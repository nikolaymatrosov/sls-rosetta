import { query } from "@ydbjs/query";

// Connection record type
export interface ConnectionRecord {
  user_id: string;
  connection_id: string;
}

// Store connection in YDB
export async function storeConnection(
  sql: ReturnType<typeof query>,
  userId: string,
  connectionId: string,
  connectedAt: Date
): Promise<void> {
  try {
    // Delete any existing connection for this user (enforce one-per-user)
    console.log(`Deleting existing connections for user: ${userId}`);
    await sql`DELETE FROM connections WHERE user_id = ${userId};`;

    // Insert new connection
    console.log(`Inserting connection: ${connectionId} for user: ${userId}`);
    await sql`
      UPSERT INTO connections (connection_id, user_id, connected_at)
      VALUES (${connectionId}, ${userId}, ${connectedAt});
    `;
    console.log(`Connection stored successfully`);
  } catch (error) {
    console.error(`Failed to store connection:`, error);
    throw error;
  }
}

// Remove connection from YDB by user_id
export async function removeConnection(
  sql: ReturnType<typeof query>,
  userId: string
): Promise<void> {
  await sql`
		DELETE FROM connections WHERE user_id = ${userId};
	`;
}

// Remove connection from YDB by connection_id
export async function removeConnectionById(
  sql: ReturnType<typeof query>,
  connectionId: string
): Promise<void> {
  await sql`
		DELETE FROM connections WHERE connection_id = ${connectionId};
	`;
}

// Get all connections from YDB
export async function getAllConnections(
  sql: ReturnType<typeof query>
): Promise<ConnectionRecord[]> {
  const result = await sql`
    SELECT user_id, connection_id FROM connections;
  `;

  // Result is nested: [[{...}, {...}]]
  return (result as unknown as Array<ConnectionRecord[]>)[0] ?? [];
}

// Get user_id by connection_id
export async function getUserIdByConnectionId(
  sql: ReturnType<typeof query>,
  connectionId: string
): Promise<string | null> {
  console.log(`Looking up user for connection: ${connectionId}`);
  const result = await sql`
    SELECT user_id FROM connections WHERE connection_id = ${connectionId};
  `;

  // Result is nested: [[{user_id: "..."}]]
  const rows =
    (result as unknown as Array<Array<{ user_id: string }>>)[0] ?? [];
  console.log(
    `Found ${rows.length} rows for connection ${connectionId}:`,
    JSON.stringify(rows)
  );
  return rows.length > 0 ? rows[0].user_id : null;
}
