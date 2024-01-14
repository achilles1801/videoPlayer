package main

import ( // import the needed packages
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
    router := gin.Default() // entry point for our application
    err := godotenv.Load() // load the .env file
    if err != nil {
        log.Fatal("Error loading .env file")
    }
    config := cors.DefaultConfig() // configure CORS
    config.AllowOrigins = []string{"*"} // allow requests from localhost:3000
    router.Use(cors.New(config)) // use the CORS configuration

    router.GET("/presign", handlePresign) // when a GET request is made to /presign, call handlePresign
    router.GET("/list", handleListBucket)  // when a GET request is made to /list, call handleListBucket
    router.DELETE("/delete", handleDelete) // when a DELETE request is made to /delete, call handleDelete


    router.Run(":8080") // serve on port 8080
}

func handleListBucket(c *gin.Context) { //responsible for listing all the videos in the bucket
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
        // Check if the object is a video file
        if strings.HasSuffix(*object.Key, ".mp4") || strings.HasSuffix(*object.Key, ".mov") || strings.HasSuffix(*object.Key, ".3gp") || strings.HasSuffix(*object.Key, ".webm") || strings.HasSuffix(*object.Key, ".ogg") {
            // Generate the URL for each object using CloudFront URL
            url := cloudFrontURL + *object.Key
            urls = append(urls, url)
        }
    }

    c.JSON(http.StatusOK, urls)
}

func handlePresign(c *gin.Context) { //generate the presigned url
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
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
            return
        }
    
    // Generate pre-signed URL
    req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
        Bucket: aws.String("majdks-video-player-bucket"),
        Key:    &key,
        // ACL:   aws.String("public-read"),
        ContentType: aws.String(contentType),
    })
    urlStr, err := req.Presign(15 * time.Minute) // URL expires in 15 minutes
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate pre-signed URL"})
        return
    }
    log.Println("pre-signed url:",urlStr)

    c.JSON(http.StatusOK, gin.H{"url": urlStr})
}
func handleDelete(c *gin.Context) {
    sess := createAWSSession()
    svc := s3.New(sess)

    // Extract videoUrl from query parameter
    videoUrl := c.Query("videoUrl")
    if videoUrl == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Video URL is required"})
        return
    }

    // Extract the key from the videoUrl
    key := strings.TrimPrefix(videoUrl, "https://d2ufs6yhgycudn.cloudfront.net/")

    // Delete the object from the bucket
    _, err := svc.DeleteObject(&s3.DeleteObjectInput{
        Bucket: aws.String("majdks-video-player-bucket"),
        Key:    aws.String(key),
    })
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete video"})
        return
    }

    // Wait to see if the object was deleted
    err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
        Bucket: aws.String("majdks-video-player-bucket"),
        Key:    aws.String(key),
    })
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete video"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Video deleted successfully"})
}


func createAWSSession() *session.Session { // similar to logging in to AWS
    awsAccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
    awsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

    sess, err := session.NewSession(&aws.Config{
        Region: aws.String("us-east-2"), // Set your AWS region
        Credentials: credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, ""),
    })
    if err != nil {
        log.Fatalf("failed to create session: %s", err.Error())
    }
    return sess
}


