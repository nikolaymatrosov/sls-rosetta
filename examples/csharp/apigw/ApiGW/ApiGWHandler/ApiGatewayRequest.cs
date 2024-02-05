using System.Text.Json;
using System.Text.Json.Serialization;

namespace ApiGWHandler;

public class ApiGatewayRequest
{
    [JsonPropertyName("version")]
    public string Version { get; set; }

    [JsonPropertyName("resource")]
    public string Resource { get; set; }

    [JsonPropertyName("path")]
    public string Path { get; set; }

    [JsonPropertyName("httpMethod")]
    public string HttpMethod { get; set; }

    [JsonPropertyName("headers")]
    public Dictionary<string, string> Headers { get; set; }

    [JsonPropertyName("multiValueHeaders")]
    public Dictionary<string, List<string>> MultiValueHeaders { get; set; }

    [JsonPropertyName("queryStringParameters")]
    public Dictionary<string, string> QueryStringParameters { get; set; }

    [JsonPropertyName("multiValueQueryStringParameters")]
    public Dictionary<string, List<string>> MultiValueQueryStringParameters { get; set; }

    [JsonPropertyName("requestContext")]
    public Dictionary<string, object> RequestContext { get; set; }

    [JsonPropertyName("pathParameters")]
    public Dictionary<string, string> PathParameters { get; set; }

    [JsonPropertyName("body")]
    public string Body { get; set; }

    [JsonPropertyName("isBase64Encoded")]
    public bool IsBase64Encoded { get; set; }

    [JsonPropertyName("parameters")]
    public Dictionary<string, string> Parameters { get; set; }

    [JsonPropertyName("multiValueParameters")]
    public Dictionary<string, List<string>> MultiValueParameters { get; set; }

    [JsonPropertyName("operationId")]
    public string OperationId { get; set; }
}