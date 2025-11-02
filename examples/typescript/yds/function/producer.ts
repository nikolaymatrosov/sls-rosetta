import { Http } from "@yandex-cloud/function-types/dist/src/http";
import { Driver } from "@ydbjs/core";
import { MetadataCredentialsProvider } from "@ydbjs/auth/metadata";
import { createTopicWriter } from "@ydbjs/topic/writer";
import { AccessTokenCredentialsProvider } from "@ydbjs/auth/access-token";

/**
 * Producer function that receives HTTP requests and writes messages to YDS topic.
 *
 * Expected request body:
 * {
 *   "message": "string",
 *   "user_id": "string",
 *   "action": "string"
 * }
 */
export async function handler(event: Http.Event): Promise<Http.Result> {
  try {
    // Parse request body
    let body: string = event.body || "{}";
    if (event.isBase64Encoded) {
      body = Buffer.from(body, "base64").toString("utf-8");
    }

    const data = JSON.parse(body);

    // Validate required fields
    const requiredFields = ["message", "user_id", "action"];
    for (const field of requiredFields) {
      if (!(field in data)) {
        return {
          statusCode: 400,
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ error: `Missing required field: ${field}` }),
        };
      }
    }

    // Get environment variables
    const ydbEndpoint = process.env.YDB_ENDPOINT;
    const ydsTopicPath = process.env.YDS_TOPIC_PATH;

    if (!ydbEndpoint || !ydsTopicPath) {
      console.error(
        "Missing environment variables: YDB_ENDPOINT or YDS_TOPIC_PATH"
      );
      return {
        statusCode: 500,
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ error: "Configuration error" }),
      };
    }

    // Prepare message with timestamp
    const message = {
      message: data.message,
      user_id: data.user_id,
      action: data.action,
      timestamp: new Date().toISOString(),
    };

    // Connect to YDB and write to topic
    console.log(`Connecting to YDB: ${ydbEndpoint}`);
    console.log(`Topic path: ${ydsTopicPath}`);

    const iamToken = process.env.YDB_IAM_TOKEN;
    let credentialsProvider;
    if (iamToken) {
      console.log("Using IAM token for authentication");
      credentialsProvider = new AccessTokenCredentialsProvider({
        token: iamToken,
      });
    } else {
      console.log("No IAM token provided; using default credentials provider");
      credentialsProvider = new MetadataCredentialsProvider();
    }

    const driver = new Driver(ydbEndpoint, {
      credentialsProvider,
    });

    try {
      console.log("Waiting for driver to be ready...");
      await driver.ready();
      console.log("Driver is ready");

      // Create topic writer using the documented API
      console.log("Creating topic writer...");
      await using writer = createTopicWriter(driver, {
        topic: ydsTopicPath,
        producer: "producer-ts",
      });
      console.log("Topic writer created");

      // Write message to topic
      const messageData = JSON.stringify(message);
      console.log(`Writing message: ${messageData}`);
      const seqNo = writer.write(
        new TextEncoder().encode(messageData)
      );
      console.log(`Message queued with sequence number: ${seqNo}`);

      console.log("Flushing writer...");
      await writer.flush();
      console.log("Writer flushed successfully");

      console.log(`Message sent to YDS topic: ${JSON.stringify(message)}`);

      return {
        statusCode: 200,
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          status: "success",
          message: "Message sent to YDS topic",
          data: message,
          sequenceNumber: seqNo.toString(),
        }),
      };
    } catch (writeError) {
      console.error("Error during YDB write operation:", writeError);
      console.error(
        "Error stack:",
        writeError instanceof Error ? writeError.stack : String(writeError)
      );
      throw writeError;
    } finally {
      // Close driver
      console.log("Closing driver...");
      await driver.close();
      console.log("Driver closed");
    }
  } catch (error) {
    if (error instanceof SyntaxError) {
      console.error("Invalid JSON in request body:", error);
      return {
        statusCode: 400,
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ error: "Invalid JSON format" }),
      };
    }

    console.error("Error processing request:", error);
    return {
      statusCode: 500,
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        error: error instanceof Error ? error.message : String(error),
      }),
    };
  }
}
