using System.Text.Json.Serialization;

namespace ApiGWHandler;

// We need to provide parameters in lower camel case
// because Yandex API Gateway expects them in this format
// Otherwise, you'll this 502 error:
// {"message":"no statusCode provided by function"}
public class Response(
    int statusCode,
    string body,
    Dictionary<string, string> headers = null,
    bool isBase64Encoded = false)
{
    [JsonPropertyName("statusCode")]
    public int StatusCode { get; set; } = statusCode;

    [JsonPropertyName("body")]
    public string Body { get; set; } = body;

    [JsonPropertyName("headers")]
    public Dictionary<string, string> Headers { get; set; } = headers;

    [JsonPropertyName("isBase64Encoded")]
    public bool IsBase64Encoded { get; set; } = isBase64Encoded;
}