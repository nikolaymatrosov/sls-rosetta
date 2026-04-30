import { faker } from "@faker-js/faker";
import { Driver } from "@ydbjs/core";
import { query } from "@ydbjs/query";
import { AccessTokenCredentialsProvider } from "@ydbjs/auth/access-token";
import { MetadataCredentialsProvider } from "@ydbjs/auth/metadata";
import { Uuid } from "@ydbjs/value/primitive";
import { randomUUID } from "node:crypto";
import { mkdir, writeFile, readFile } from "node:fs/promises";
import { existsSync } from "node:fs";
import { join } from "node:path";

const USER_COUNT = 500;
const PROGRESS_FILE = join("./data", "progress.json");

type Progress = Record<string, number>;

async function readProgress(): Promise<Progress> {
  if (!existsSync(PROGRESS_FILE)) return {};
  return JSON.parse(await readFile(PROGRESS_FILE, "utf-8")) as Progress;
}

async function saveProgress(progress: Progress): Promise<void> {
  await writeFile(PROGRESS_FILE, JSON.stringify(progress, null, 2));
}

const POSTS_MIN = 10;
const POSTS_MAX = 100;
const COMMENTS_MIN = 10;
const COMMENTS_MAX = 100;
const LIKES_MIN = 0;
const LIKES_MAX = 200;
const BATCH_SIZE = 500;
const DATA_DIR = "./data";

function randInt(min: number, max: number): number {
  return Math.floor(Math.random() * (max - min + 1)) + min;
}

function sample<T>(items: T[]): T {
  return items[Math.floor(Math.random() * items.length)];
}

// ── Plain data types (stored in files as JSON) ───────────────────────────────

interface UserRow {
  id: string;
  name: string;
  email: string;
}

interface PostRow {
  id: string;
  user_id: string;
  title: string;
  content: string;
}

interface CommentRow {
  id: string;
  post_id: string;
  user_id: string;
  content: string;
}

interface LikeRow {
  id: string;
  post_id: string;
  user_id: string;
}

type SqlClient = ReturnType<typeof query>;

// ── Phase 1: Generate data and dump to files ─────────────────────────────────

async function generateData(): Promise<void> {
  await mkdir(DATA_DIR, { recursive: true });

  // Generate users
  console.log(`Generating ${USER_COUNT} users …`);
  const users: UserRow[] = [];
  const userIds: string[] = [];

  for (let i = 0; i < USER_COUNT; i++) {
    const id = randomUUID();
    userIds.push(id);
    users.push({
      id,
      name: faker.person.fullName(),
      email: faker.internet.email(),
    });
  }
  await writeFile(join(DATA_DIR, "users.json"), JSON.stringify(users));
  console.log(`  Wrote ${users.length} users to ${DATA_DIR}/users.json\n`);

  // Generate posts, comments, and likes
  console.log("Generating posts, comments, and likes …");
  const posts: PostRow[] = [];
  const comments: CommentRow[] = [];
  const likes: LikeRow[] = [];

  for (let u = 0; u < userIds.length; u++) {
    const userId = userIds[u];
    const postCount = randInt(POSTS_MIN, POSTS_MAX);

    for (let p = 0; p < postCount; p++) {
      const postId = randomUUID();
      posts.push({
        id: postId,
        user_id: userId,
        title: faker.lorem.sentence({ min: 3, max: 8 }),
        content: faker.lorem.paragraphs({ min: 1, max: 3 }),
      });

      const commentCount = randInt(COMMENTS_MIN, COMMENTS_MAX);
      for (let c = 0; c < commentCount; c++) {
        comments.push({
          id: randomUUID(),
          post_id: postId,
          user_id: sample(userIds),
          content: faker.lorem.sentence({ min: 3, max: 15 }),
        });
      }

      const likeCount = randInt(LIKES_MIN, LIKES_MAX);
      for (let l = 0; l < likeCount; l++) {
        likes.push({
          id: randomUUID(),
          post_id: postId,
          user_id: sample(userIds),
        });
      }
    }

    if ((u + 1) % 50 === 0 || u === userIds.length - 1) {
      console.log(
        `  Users processed: ${u + 1}/${userIds.length} | ` +
          `posts: ${posts.length}, ` +
          `comments: ${comments.length}, ` +
          `likes: ${likes.length}`
      );
    }
  }

  await writeFile(join(DATA_DIR, "posts.json"), JSON.stringify(posts));
  await writeFile(join(DATA_DIR, "comments.json"), JSON.stringify(comments));
  await writeFile(join(DATA_DIR, "likes.json"), JSON.stringify(likes));

  console.log(`\n  Wrote ${posts.length} posts    → ${DATA_DIR}/posts.json`);
  console.log(`  Wrote ${comments.length} comments → ${DATA_DIR}/comments.json`);
  console.log(`  Wrote ${likes.length} likes    → ${DATA_DIR}/likes.json`);
}

// ── Phase 2: Upload data from files in batches ───────────────────────────────

async function uploadTable<T>(
  filename: string,
  label: string,
  upload: (batch: T[]) => Promise<void>,
  progress: Progress,
  batchSize: number = BATCH_SIZE
): Promise<void> {
  const raw = await readFile(join(DATA_DIR, filename), "utf-8");
  const rows: T[] = JSON.parse(raw) as T[];
  const total = rows.length;
  const startAt = progress[label] ?? 0;

  if (startAt >= total) {
    console.log(`Skipping ${label} — already complete (${total}/${total}).\n`);
    return;
  }

  if (startAt > 0) {
    console.log(`Resuming ${label} from row ${startAt}/${total} …`);
  } else {
    console.log(`Uploading ${total} ${label} …`);
  }

  for (let i = startAt; i < total; i += batchSize) {
    const batch = rows.slice(i, i + batchSize);
    await upload(batch);
    progress[label] = i + batch.length;
    await saveProgress(progress);
    console.log(`  ${label}: ${progress[label]}/${total} (batch ${Math.ceil(progress[label] / batchSize)})`);
  }
  console.log(`  Done — ${progress[label]} ${label} uploaded.\n`);
}

async function uploadData(sql: SqlClient): Promise<void> {
  const progress = await readProgress();

  await uploadTable<UserRow>("users.json", "users", async (batch) => {
    const rows = batch.map((r) => ({ id: new Uuid(r.id), name: r.name, email: r.email }));
    await sql`UPSERT INTO users SELECT Unwrap(id) AS id, name, email FROM AS_TABLE(${rows})`;
  }, progress);

  await uploadTable<PostRow>("posts.json", "posts", async (batch) => {
    const rows = batch.map((r) => ({
      id: new Uuid(r.id),
      user_id: new Uuid(r.user_id),
      title: r.title,
      content: r.content,
    }));
    await sql`UPSERT INTO posts SELECT Unwrap(id) AS id, user_id, title, content FROM AS_TABLE(${rows})`;
  }, progress);

  await uploadTable<CommentRow>("comments.json", "comments", async (batch) => {
    const rows = batch.map((r) => ({
      id: new Uuid(r.id),
      post_id: new Uuid(r.post_id),
      user_id: new Uuid(r.user_id),
      content: r.content,
    }));
    await sql`UPSERT INTO comments SELECT Unwrap(id) AS id, post_id, user_id, content FROM AS_TABLE(${rows})`;
  }, progress, 2000);

  await uploadTable<LikeRow>("likes.json", "likes", async (batch) => {
    const rows = batch.map((r) => ({
      id: new Uuid(r.id),
      post_id: new Uuid(r.post_id),
      user_id: new Uuid(r.user_id),
    }));
    await sql`UPSERT INTO likes SELECT Unwrap(id) AS id, post_id, user_id FROM AS_TABLE(${rows})`;
  }, progress, 10000);
}

// ── Main ─────────────────────────────────────────────────────────────────────

async function main() {
  // ── Phase 1: Generate ────────────────────────────────────────────────────
  const dataExists =
    existsSync(join(DATA_DIR, "users.json")) &&
    existsSync(join(DATA_DIR, "posts.json")) &&
    existsSync(join(DATA_DIR, "comments.json")) &&
    existsSync(join(DATA_DIR, "likes.json"));

  if (dataExists) {
    console.log(`Data files already exist in ${DATA_DIR}/, skipping generation.\n`);
  } else {
    await generateData();
  }

  // ── Phase 2: Upload ──────────────────────────────────────────────────────
  const ydbEndpoint = process.env.YDB_ENDPOINT;
  const ydbDatabase = process.env.YDB_DATABASE;

  if (!ydbEndpoint || !ydbDatabase) {
    console.error("Set YDB_ENDPOINT and YDB_DATABASE environment variables.");
    console.error("Example:");
    console.error(
      "  export YDB_ENDPOINT=grpcs://ydb.serverless.yandexcloud.net:2135"
    );
    console.error("  export YDB_DATABASE=/ru-central1/b1g.../etn...");
    process.exit(1);
  }

  const connectionString = `${ydbEndpoint}/?database=${ydbDatabase}`;
  const iamToken = process.env.YDB_IAM_TOKEN;
  const credentialsProvider = iamToken
    ? new AccessTokenCredentialsProvider({ token: iamToken })
    : new MetadataCredentialsProvider();

  console.log(`Connecting to ${connectionString} …`);
  const driver = new Driver(connectionString, { credentialsProvider });
  await driver.ready();
  console.log("Connected.\n");

  const sql = query(driver);

  try {
    await uploadData(sql);
    console.log("Seeding complete.");
  } catch (err) {
    console.error("Error during upload:", JSON.stringify(err, null, 2));
  } finally {
    await driver.close();
    console.log("Connection closed.");
  }
}

main().catch((err) => {
  console.error("Seed failed:", err);
  process.exit(1);
});
