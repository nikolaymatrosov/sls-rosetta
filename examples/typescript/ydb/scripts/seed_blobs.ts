import { faker } from "@faker-js/faker";
import { Driver } from "@ydbjs/core";
import { query } from "@ydbjs/query";
import { AccessTokenCredentialsProvider } from "@ydbjs/auth/access-token";
import { MetadataCredentialsProvider } from "@ydbjs/auth/metadata";
import { Uint8, Bytes } from "@ydbjs/value/primitive";
import { randomBytes } from "node:crypto";

const KEY_COUNT = 100;
const PARTS_MIN = 1;
const PARTS_MAX = 10;
const PART_SIZE = 1024 * 1024; // 1 MB

function generateKey(): string {
  return `${faker.word.adjective()}-${faker.word.noun()}-${faker.number.int({ min: 1, max: 9999 })}`;
}

function randInt(min: number, max: number): number {
  return Math.floor(Math.random() * (max - min + 1)) + min;
}

type SqlClient = ReturnType<typeof query>;

async function seedBlobs(sql: SqlClient): Promise<void> {
  const keys = Array.from({ length: KEY_COUNT }, generateKey);

  let totalParts = 0;

  for (let k = 0; k < 1; k++) {
    const key = keys[k];
    const partCount = randInt(PARTS_MIN, PARTS_MAX);

    for (let p = 0; p < partCount; p++) {
      const data = randomBytes(PART_SIZE);
      const row = [{ key, partNumber: new Uint8(p), data: new Bytes(data) }];
      let stmt = sql`UPSERT INTO blobs SELECT key, partNumber, data FROM AS_TABLE(${row})`;
      console.log("text:", stmt.text, "parameters:", stmt.parameters);
      await stmt;
      totalParts++;
    }

    console.log(
      `  [${k + 1}/${keys.length}] key="${key}" parts=${partCount} (total parts uploaded: ${totalParts})`
    );
  }

  console.log(`\nDone — ${KEY_COUNT} keys, ${totalParts} parts, ~${totalParts} MB uploaded.`);
}

async function main() {
  const ydbEndpoint = process.env.YDB_ENDPOINT;
  const ydbDatabase = process.env.YDB_DATABASE;

  if (!ydbEndpoint || !ydbDatabase) {
    console.error("Set YDB_ENDPOINT and YDB_DATABASE environment variables.");
    console.error("  export YDB_ENDPOINT=grpcs://ydb.serverless.yandexcloud.net:2135");
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
    await seedBlobs(sql);
    console.log("Blob seeding complete.");
  } catch (err) {
    console.error("Error during seeding:", JSON.stringify(err, null, 2));
  } finally {
    await driver.close();
    console.log("Connection closed.");
  }
}

main().catch((err) => {
  console.error("Seed failed:", err);
  process.exit(1);
});
