# Video Player with AWS S3 and Go

This project is a simple video player that uses AWS S3 for video storage and Go for the backend. The frontend is served by the Go server and makes requests to the server to get presigned URLs for uploading and accessing videos on S3.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

- Go (version 1.16 or later)
- AWS account with an S3 bucket

### Installing

1. Clone the repository to your local machine.
2. Navigate to the project directory.
3. Run `go build` to compile the Go server.
4. Run `./main` (or `main.exe` on Windows) to start the server.

The server will start on `localhost:8080`. You can access the video player by navigating to `http://localhost:8080` in your browser.

## Usage

- To get a presigned URL for uploading a video to S3, make a GET request to `/presign` with the filename as a query parameter.
- To list the objects in the S3 bucket, make a GET request to `/list`.

## Built With

- [Go](https://golang.org/)
- [Gin](https://github.com/gin-gonic/gin) - HTTP web framework used
- [AWS SDK for Go](https://aws.amazon.com/sdk-for-go/) - AWS SDK for the Go programming language

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
