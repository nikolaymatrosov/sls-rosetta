// See https://aka.ms/new-console-template for more information

using Amazon;
using Amazon.Runtime;
using Amazon.S3;
using Amazon.S3.Transfer;

AWSConfigs.LoggingConfig.LogTo = LoggingOptions.Console;
var s3Config = new AmazonS3Config
{
    ServiceURL = "https://s3.yandexcloud.net",
    AuthenticationRegion = "ru-central1",
    ForcePathStyle = true,
    // AuthenticationRegion = "us-east-1",
    // SignatureVersion ,
    // SignatureMethod = SigningAlgorithm.HmacSHA1
};
var s3Cred = new EnvironmentVariablesAWSCredentials();
using var _s3Client = new AmazonS3Client(s3Cred, s3Config);

var s3Transfer = new TransferUtility(_s3Client);

string filePath = "/Users/nikthespirit/Downloads/star.png";
string bucketName = "test-api-gw-bucket";
string key = "star.png";
s3Transfer.Upload(filePath, bucketName, key);

Console.WriteLine("File uploaded successfully");