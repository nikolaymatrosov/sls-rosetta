import { Http } from "@yandex-cloud/function-types/dist/src/http";
import { Driver } from "@ydbjs/core";
import { query } from "@ydbjs/query";
import { Uuid } from '@ydbjs/value/primitive'
import { StatsMode } from "@ydbjs/api/query";
import { MetadataCredentialsProvider } from "@ydbjs/auth/metadata";
import { AccessTokenCredentialsProvider } from "@ydbjs/auth/access-token";
import { randomUUID } from "crypto";

/**
 * HTTP handler that queries YDB across multiple tables and returns results
 * with execution statistics.
 *
 * Demonstrates how to use @ydbjs/query with StatsMode.FULL to retrieve
 * transaction statistics including per-table read rows and bytes —
 * useful for estimating read unit consumption.
 */
export async function handler(event: Http.Event): Promise<Http.Result> {
  const ydbEndpoint = process.env.YDB_ENDPOINT;
  const ydbDatabase = process.env.YDB_DATABASE;

  if (!ydbEndpoint || !ydbDatabase) {
    console.error(
      "Missing environment variables: YDB_ENDPOINT or YDB_DATABASE",
    );
    return {
      statusCode: 500,
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ error: "Configuration error" }),
    };
  }

  const userId = event.queryStringParameters?.userId;
  if (!userId) {
    return {
      statusCode: 400,
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        error: "Missing required query parameter: userId",
      }),
    };
  }

  const connectionString = `${ydbEndpoint}/?database=${ydbDatabase}`;
  console.log(`Connecting to YDB: ${connectionString}`);

  const iamToken = process.env.YDB_IAM_TOKEN;
  const credentialsProvider = iamToken
    ? new AccessTokenCredentialsProvider({ token: iamToken })
    : new MetadataCredentialsProvider();

  const driver = new Driver(connectionString, { credentialsProvider });

  try {
    await driver.ready();
    console.log("Driver is ready");

    const sql = query(driver);

    // Query a single user across multiple tables with StatsMode.FULL to observe
    // per-table read unit consumption (rows and bytes read from each table).
    // const q = sql`
    //   SELECT
    //     u.id        AS user_id,
    //     u.name      AS user_name,
    //     u.email     AS user_email,
    //     p.id        AS post_id,
    //     p.title     AS post_title,
    //     COUNT(c.id) AS comment_count,
    //     COUNT(l.id) AS like_count
    //   FROM users AS u
    //   LEFT JOIN posts AS p ON p.user_id = u.id
    //   LEFT JOIN comments AS c ON c.post_id = p.id
    //   LEFT JOIN likes AS l ON l.post_id = p.id
    //   WHERE u.id = CAST(${userId} AS UUID)
    //   GROUP BY u.id, u.name, u.email, p.id, p.title
    //   ORDER BY post_title
    // `.withStats(StatsMode.FULL);
    const users = [
      {
        id: "e061e6a1-65eb-47f2-b71e-0839ec3d3c72",
        name: "Dr. Dave Collier",
        email: "Dino_Windler99@yahoo.com",
      },
      {
        id: "13b09b5d-0740-4459-8a2e-d8eda58663a1",
        name: "Ms. Mildred Orn",
        email: "Anabelle80@hotmail.com",
      },
      {
        id: "ecff8f66-70a6-4b3e-9266-d9667e899a6a",
        name: "Tricia Wuckert II",
        email: "Daniel.Walsh82@yahoo.com",
      },
      {
        id: "cf26f862-6d4d-4d2a-a53b-c442aeb8e3d4",
        name: "Vince Kuhn",
        email: "Megan94@yahoo.com",
      },
      {
        id: "a708e1ca-6608-4287-8e60-62797eee71b5",
        name: "Ernie Halvorson",
        email: "Gordon_Langosh@yahoo.com",
      },
      {
        id: "153b6b40-6a49-4553-8d22-2fd90a6931bf",
        name: "Ginger Feeney",
        email: "Frankie_Borer2@yahoo.com",
      },
      {
        id: "575c8332-cfe7-42dc-b975-8451df000e8a",
        name: "Lucas Gutmann",
        email: "Dwayne38@hotmail.com",
      },
      {
        id: "572ed53b-a79d-4585-b192-0bbdcdef19c0",
        name: "Jenna Beer",
        email: "Tommie55@yahoo.com",
      },
      {
        id: "8bbe93d7-40c5-4ade-8471-20514ebd177f",
        name: "Taylor Harvey",
        email: "Margaretta_Thompson@hotmail.com",
      },
      {
        id: "4c90f201-8ee3-4a8d-97da-b72f6b1dfc86",
        name: "Sammy Kautzer",
        email: "Faye_Adams@hotmail.com",
      },
      {
        id: "bd009301-08ca-442c-baaf-a2c83fc49579",
        name: "Christine Bailey",
        email: "Delpha.OKon53@hotmail.com",
      },
      {
        id: "66513dd9-052c-41b2-ba57-743310575873",
        name: "Griffin Kiehn",
        email: "Antoinette17@yahoo.com",
      },
      {
        id: "6720df03-f761-44c5-a5fd-c648cd71a019",
        name: "Robyn Reynolds",
        email: "Elisa6@gmail.com",
      },
      {
        id: "dc1df0c4-129a-4ec7-9536-95f8273b8629",
        name: "Mark Kuvalis",
        email: "Art19@hotmail.com",
      },
      {
        id: "e685833c-04de-407b-afe5-dd3bc59f7419",
        name: "Fern Yost",
        email: "Jazlyn46@gmail.com",
      },
    ];
    const units: number[] = [];
    for (let i = 0; i < 100; i++) {
      const uuid = new Uuid(randomUUID());
      const q = sql`
      SELECT key, partNumber
      FROM blobs
      WHERE key = 'sinful-tomatillo-9383'
      `.withStats(StatsMode.FULL);
      q.on("metadata", (s) => {
        const consumedUnits = parseConsumedUnits(s.get("x-ydb-consumed-units"));
        units.push(consumedUnits);
        console.log("x-ydb-consumed-units", consumedUnits);
      });

      const resultSets = await q;
      const stats = q.stats();

      const rows = resultSets[0] ?? [];

      console.log(`Query returned ${rows.length} rows`);
      const replacer = (_key: string, value: unknown) => {
        if (typeof value === "bigint") return String(value);
        if (Buffer.isBuffer(value)) return value.toString("utf-8");
        if (
          value !== null &&
          typeof value === "object" &&
          (value as { type?: string }).type === "Buffer" &&
          Array.isArray((value as { data?: unknown[] }).data)
        ) {
          return Buffer.from((value as { data: number[] }).data).toString("utf-8");
        }
        return value;
      };
      // console.log("Query stats:", JSON.stringify(stats, replacer));
    }
    console.log(units);
    return {
      statusCode: 200,
      headers: { "Content-Type": "application/json" },
      // body: JSON.stringify({ rows, stats }, bigintReplacer),
    };
  } catch (error) {
    console.error("Error executing query:", JSON.stringify(error));
    return {
      statusCode: 500,
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        error: error instanceof Error ? error.message : String(error),
      }),
    };
  } finally {
    await driver.close();
    console.log("Driver closed");
    
  }
}


function parseConsumedUnits(stats: string | undefined): number {
  if (!stats) return 0;
  const match = stats.match(/^(\d+),/);
  return match ? parseInt(match[1], 10) : 0;
}