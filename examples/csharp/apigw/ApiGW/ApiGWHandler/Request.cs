using System.Text.Json.Serialization;

namespace ApiGWHandler;

public record Request
{
    public String Name { get; init; }
}