package main

import (
	"context"
	"encoding/json"
	"net/http" // Added this line
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"

	// "github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// ListHandler is the Lambda function handler
func ListHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    sess, err := createAWSSessionForList()
    if err != nil {
        return serverError(err)
    }

    svc := s3.New(sess)
    bucketName := os.Getenv("BUCKET_NAME") // Set your bucket name in environment variables
    cloudFrontURL := os.Getenv("CLOUDFRONT_URL") // CloudFront URL

    result, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
        Bucket: aws.String(bucketName),
    })
    if err != nil {
        return serverError(err)
    }

    var urls []string
    for _, object := range result.Contents {
        if isVideoFile(*object.Key) {
            url := cloudFrontURL + *object.Key
            urls = append(urls, url)
        }
    }

    return apiResponse(http.StatusOK, urls)
}

// serverError sends an API Gateway proxy response on error
func serverError(err error) (events.APIGatewayProxyResponse, error) {
    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusInternalServerError,
        Body:       err.Error(),
    }, nil
}

// apiResponse sends a successful API Gateway proxy response
func apiResponse(status int, body interface{}) (events.APIGatewayProxyResponse, error) {
    resp, err := json.Marshal(body)
    if err != nil {
        return serverError(err)
    }

    return events.APIGatewayProxyResponse{
        StatusCode: status,
        Body:       string(resp),
        Headers:    map[string]string{
            "Content-Type": "application/json",
            "Access-Control-Allow-Origin": "*", // Add this line
        },
    }, nil
}

// isVideoFile checks if the filename corresponds to a supported video file
func isVideoFile(filename string) bool {
    return strings.HasSuffix(filename, ".mp4") ||
        strings.HasSuffix(filename, ".mov") ||
        strings.HasSuffix(filename, ".3gp") ||
        strings.HasSuffix(filename, ".webm") ||
        strings.HasSuffix(filename, ".ogg")
}

// createAWSSessionForList creates an AWS session for accessing S3
func createAWSSessionForList() (*session.Session, error) {
    awsRegion := os.Getenv("AWS_REGION")
    awsAccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
    awsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

    return session.NewSession(&aws.Config{
        Region:      aws.String(awsRegion),
        Credentials: credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, ""),
    })
}

func main() {
    lambda.Start(ListHandler)
}