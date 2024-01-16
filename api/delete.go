package main

import (
    "context"
    "net/http"
    "os"
    "strings"

    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
)

func DeleteHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    sess, err := session.NewSession(&aws.Config{
        Region: aws.String(os.Getenv("AWS_REGION")),
    })
    if err != nil {
        return serverError(err)
    }

    svc := s3.New(sess)

    // Extract the full CloudFront URL from the request parameters
    videoUrl := request.QueryStringParameters["videoUrl"]
    if videoUrl == "" {
        return clientError(http.StatusBadRequest, "Video URL is required")
    }

    // Extract the key from the CloudFront URL
    key := strings.TrimPrefix(videoUrl, os.Getenv("CLOUDFRONT_URL"))
    if key == "" {
        return clientError(http.StatusBadRequest, "Invalid video URL")
    }

    // Delete the object from the bucket
    _, err = svc.DeleteObject(&s3.DeleteObjectInput{
        Bucket: aws.String(os.Getenv("BUCKET_NAME")),
        Key:    aws.String(key),
    })
    if err != nil {
        return serverError(err)
    }

    return successResponse("Video deleted successfully")
}

func serverError(err error) (events.APIGatewayProxyResponse, error) {
    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusInternalServerError, 
        Body:       err.Error(),
        Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
    }, nil
}

func clientError(statusCode int, message string) (events.APIGatewayProxyResponse, error) {
    return events.APIGatewayProxyResponse{
        StatusCode: statusCode, 
        Body:       message,
        Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
    }, nil
}

func successResponse(message string) (events.APIGatewayProxyResponse, error) {
    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusOK, 
        Body:       message,
        Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
    }, nil
}
    func main() {
        lambda.Start(DeleteHandler)
        }