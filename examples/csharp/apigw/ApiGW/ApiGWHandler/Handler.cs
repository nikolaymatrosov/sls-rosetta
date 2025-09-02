using System.Text.Json;
using Yandex.Cloud.Functions;

namespace ApiGWHandler;

public class Handler : YcFunction<ApiGatewayRequest, Response>
{
    public Response FunctionHandler(ApiGatewayRequest request, Context c)
    {
        Console.WriteLine($"Function name: {c.FunctionId}");
        Console.WriteLine($"Function version: {c.FunctionVersion}");
        Console.WriteLine($"Service account token: {c.TokenJson}");
        Console.WriteLine($"Request body: {request.Body}");
        var data = JsonSerializer.Deserialize<Request>(request.Body, new JsonSerializerOptions()
            {
                PropertyNamingPolicy = JsonNamingPolicy.CamelCase
            }
        );
    
        string name = data?.Name ?? "World";

        string body = JsonSerializer.Serialize(
            new Dictionary<string, string>
            {
                { "message", $"Hello, {name}!" },
            });

        Console.WriteLine($"Response body: {body}");

        return new Response(200, body, new Dictionary<string, string>
        {
            { "Content-Type", "application/json" },
        });
    }
}