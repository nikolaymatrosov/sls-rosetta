namespace Postbox;

public class Response
{
    public int statusCode { get; set; }
    public String body { get; set; }

    public Response(int statusCode, String body)
    {
        this.statusCode = statusCode;
        this.body = body;
    }
}