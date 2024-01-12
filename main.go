package main

import (
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
    "github.com/gin-contrib/cors"
)

func main() {
    router := gin.Default() // entry point for our application

    router.GET("/presign", handlePresign) // when a GET request is made to /presign, call handlePresign
    router.GET("/list", handleListBucket)  // when a GET request is made to /list, call handleListBucket

    config := cors.DefaultConfig()
    config.AllowOrigins = []string{"http://localhost:8080"} // Replace with your frontend URL
    router.Use(cors.New(config))
    
    // Serve the main app when visiting the root path "/"
    router.GET("/", func(c *gin.Context) {
        c.File("./frontend/index.html")
    })
    router.Static("/js", "./frontend/js")

    router.Run(":8080") // serve on port 8080
}

func handleListBucket(c *gin.Context) {
    sess := createAWSSession() // create a new AWS session
    svc := s3.New(sess) // create a new S3 client

    bucketName := "majdks-video-player-bucket"  // name of the bucket

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

func handlePresign(c *gin.Context) {
    sess := createAWSSession()
    svc := s3.New(sess)

    // Extract filename from query parameter
    filename := c.Query("filename")
    if filename == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Filename is required"})
        return
    }

    // Generate a unique key based on the filename
    uniqueID := uuid.New().String()
    timestamp := time.Now().Format("20060102150405") // Format: YYYYMMDDHHMMSS
    key := "uploads/" + timestamp + "-" + uniqueID + "-" + filename

    // Generate pre-signed URL
    req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
        Bucket: aws.String("majdks-video-player-bucket"),
        Key:    &key,
        ACL:   aws.String("public-read"),
    })
    urlStr, err := req.Presign(15 * time.Minute) // URL expires in 15 minutes
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate pre-signed URL"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"url": urlStr})
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
        log.Fatalf("failed to create session: %s", err.Error())
    }
    return sess
}


