package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {
    router := gin.Default() // entry point for our application

    config := cors.DefaultConfig()
    config.AllowOrigins = []string{"http://localhost:3000"} // Replace with your frontend URL
    router.Use(cors.New(config))

    router.GET("/presign", handlePresign) // when a GET request is made to /presign, call handlePresign
    router.GET("/list", handleListBucket)  // when a GET request is made to /list, call handleListBucket

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

    // Create a slice to hold the URLs
    urls := []string{}

    // CloudFront distribution URL
    cloudFrontURL := "https://d2ufs6yhgycudn.cloudfront.net/"

    // Iterate over the objects in the bucket
    for _, object := range result.Contents {
        // Check if the object is a .mp4 video
        if strings.HasSuffix(*object.Key, ".mp4") {
            // Generate the URL for each object using CloudFront URL
            url := cloudFrontURL + *object.Key
            urls = append(urls, url)
        }
    }

    c.JSON(http.StatusOK, urls)
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
    log.Println("pre-signed url:",urlStr)

    c.JSON(http.StatusOK, gin.H{"url": urlStr})
}


func createAWSSession() *session.Session {
    sess, err := session.NewSession(&aws.Config{
        Region: aws.String("us-east-2"), // Set your AWS region
        Credentials: credentials.NewStaticCredentials("AKIAYP5CN6RKM73QN2QR", "pIP23PYfDOhx6qgZ8muTi53QDBuv6DKsKfOZd8/W", ""),
    })
    if err != nil {
        log.Fatalf("failed to create session: %s", err.Error())
    }
    return sess
}


