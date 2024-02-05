namespace ApiGWHandler;

public class Response
{
    public int statusCode { get; set; }
    public String body { get; set; }
    public Dictionary<string, string> headers { get; set; }

    public bool isBase64Encoded { get; set; }

    public Response(int statusCode, String body, Dictionary<string, string> headers = null, bool isBase64Encoded = false)
    {
        this.statusCode = statusCode;
        this.body = body;
        this.headers = headers;
        this.isBase64Encoded = isBase64Encoded;
    }
}