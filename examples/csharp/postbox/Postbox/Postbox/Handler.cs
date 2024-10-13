using Amazon.Runtime;
using Amazon.SimpleEmailV2;
using Amazon.SimpleEmailV2.Model;
using Newtonsoft.Json;


namespace Postbox;

public class Handler
{
    // ReSharper disable once UnusedMember.Global
    public async Task<Response> FunctionHandler(Request request)
    {
        var client = new AmazonSimpleEmailServiceV2Client(
            new AmazonSimpleEmailServiceV2Config
            {
                ServiceURL = "https://postbox.cloud.yandex.net",
                SignatureMethod = SigningAlgorithm.HmacSHA256,
                SignatureVersion = "4",
                AuthenticationRegion = "ru-central1"
            });
        string messageId = null;
        try
        {
            var response = await client.SendEmailAsync(
                new SendEmailRequest
                {
                    Destination = new Destination
                    {
                        ToAddresses = ["nikolay.matrosov@gmail.com"]
                    },
                    Content = new EmailContent
                    {
                        Simple = new Message
                        {
                            Body = new Body
                            {
                                Text = new Content
                                {
                                    Charset = "UTF-8",
                                    Data = "Hello, world!"
                                }
                            },
                            Subject = new Content
                            {
                                Charset = "UTF-8",
                                Data = "Test email"
                            }
                        }
                    },
                    FromEmailAddress = "noreply@ycld.ru"
                });
            messageId = response.MessageId;
        }
        catch (Exception ex)
        {
            // Log the exception as JSON
            Console.WriteLine(JsonConvert.SerializeObject(ex));
        }

        return new Response(
            statusCode: 200,
            body: messageId ?? "Failed to send email"
        );
    }
}