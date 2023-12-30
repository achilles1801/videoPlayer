package main

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func main() {
    
	router := gin.Default() // initialize gin router

	router.POST("/upload", handleFileUpload) // handle file upload
	router.GET("/list", handleListBucket)  // New route for listing contents


	router.Run(":8080") // serve on port 8080
	
}
func handleListBucket(c *gin.Context) {
    sess := createAWSSession()
    svc := s3.New(sess)

    bucketName := "majdks-video-player-bucket"  // Replace with your bucket name

    // Call S3 to list current objects
    result, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
        Bucket: aws.String(bucketName),
    })
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, result)
}


func handleFileUpload(c *gin.Context) {
    sess := createAWSSession()
    svc := s3.New(sess)

    // Define bucket and object name (key) for the upload
    bucket := "majdks-video-player-bucket"
	uniqueID := uuid.New().String()
    timestamp := time.Now().Format("20060102150405") // Format: YYYYMMDDHHMMSS
    key := "uploads/" + timestamp + "-" + uniqueID + ".mp4"

    // Generate pre-signed URL
    req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
        Bucket: &bucket,
        Key:    &key,
    })
    urlStr, err := req.Presign(15 * time.Minute) // URL expires in 15 minutes
    if err != nil {
        // Handle error
    }

    c.JSON(http.StatusOK, gin.H{"url": urlStr})
}

func createAWSSession() *session.Session {
    sess, err := session.NewSession(&aws.Config{
        Region: aws.String("us-east-2"), // Set your AWS region
    })
    if err != nil {
        // Handle session creation error
    }
    return sess
}

