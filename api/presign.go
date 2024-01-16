package main

import (
    "context"
    "encoding/json"
    "log"
    "os"
    "strings"
    "time"

    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/google/uuid"
    "github.com/aws/aws-sdk-go/service/sts"
    
)

func PresignHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    log.Println("Starting PresignHandler")
    sess, err := createAWSSessionForPresign()
    if err != nil {
        log.Println("Error creating AWS session:", err)
        return events.APIGatewayProxyResponse{
            StatusCode: 500,
            Body: err.Error(),
            Headers: map[string]string{
                "Access-Control-Allow-Origin": "*",
            },
        }, nil
    }
        // Add this code to create an STS client and call GetCallerIdentity
        stsSvc := sts.New(sess)
        identity, err := stsSvc.GetCallerIdentity(&sts.GetCallerIdentityInput{})
        if err != nil {
            log.Println("Error getting caller identity:", err)
        } else {
            log.Println("Caller identity:", identity)
        }

    svc := s3.New(sess)
    filename := request.QueryStringParameters["filename"]
    if filename == "" {
        log.Println("Filename is empty")
        return events.APIGatewayProxyResponse{
            StatusCode: 400,
            Body: "Filename is required",
            Headers: map[string]string{
                "Access-Control-Allow-Origin": "*",
            },
        }, nil
    }

    // Determine the content type based on the file extension
    var contentType string
    if strings.HasSuffix(filename, ".mp4") {
        contentType = "video/mp4"
    } else if strings.HasSuffix(filename, ".mov") {
        contentType = "video/quicktime"
    } else if strings.HasSuffix(filename, ".3gp") {
        contentType = "video/3gpp"
    } else if strings.HasSuffix(filename, ".webm") {
        contentType = "video/webm"
    } else if strings.HasSuffix(filename, ".ogg") {
        contentType = "video/ogg"
    } else {
        log.Println("Invalid file type:", filename)
        return events.APIGatewayProxyResponse{
            StatusCode: 400,
            Body: "Invalid file type",
            Headers: map[string]string{
                "Access-Control-Allow-Origin": "*",
            },
        }, nil
    }

    uniqueID := uuid.New().String()
    timestamp := time.Now().Format("20060102150405")
    key := "uploads/" + timestamp + "-" + uniqueID + "-" + filename

    req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
        Bucket:      aws.String("majdks-video-player-bucket"),
        Key:         &key,
        ContentType: aws.String(contentType),
    })
    log.Println("content type:", contentType)

    urlStr, err := req.Presign(15 * time.Minute)
    if err != nil {
        log.Println("Failed to generate pre-signed URL:", err)
        return events.APIGatewayProxyResponse{
            StatusCode: 500,
            Body:       "Failed to generate pre-signed URL: " + err.Error(),
            Headers: map[string]string{
                "Access-Control-Allow-Origin": "*",
            },
        }, nil
    }
    
    response := struct {
        URL string `json:"url"`
    }{
        URL: urlStr,
    }
    
    body, err := json.Marshal(response)
    if err != nil {
        log.Println("Failed to marshal response to JSON:", err)
        return events.APIGatewayProxyResponse{
            StatusCode: 500,
            Body:       "Failed to marshal response to JSON: " + err.Error(),
            Headers: map[string]string{
                "Access-Control-Allow-Origin": "*",
            },
        }, nil
    }
    
    log.Println("Successfully generated pre-signed URL")
    return events.APIGatewayProxyResponse{
        StatusCode: 200,
        Body:       string(body),
        Headers: map[string]string{
            "Access-Control-Allow-Origin": "*",
        },
    }, nil
}

func createAWSSessionForPresign() (*session.Session, error) {
    log.Println("Creating AWS session for presign")
    awsRegion := os.Getenv("AWS_REGION")
    log.Println("AWS_REGION:", awsRegion)

    return session.NewSession(&aws.Config{
        Region: aws.String(awsRegion),
    })
}

func main() {
    lambda.Start(PresignHandler)
}
