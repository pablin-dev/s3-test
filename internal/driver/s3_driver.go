package driver

// Removed unused imports: "github.com/aws/aws-sdk-go-v2/aws", "github.com/aws/aws-sdk-go-v2/service/s3"

// S3Driver implements storage operations using AWS S3.
type S3Driver struct {
	client S3ClientAPI // Uses the interface
	bucket string
}

// NewS3Driver creates a new instance of S3Driver with an injected client.
func NewS3Driver(client S3ClientAPI, bucket string) *S3Driver {
	return &S3Driver{
		client: client,
		bucket: bucket,
	}
}
