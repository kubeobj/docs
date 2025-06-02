package main

import (
	"context"
	"flag"
	"log"
	"os"
	"path/filepath" // To get the base filename

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types" // For ACL if needed
)

func main() {
	// Define command-line flags
	bucketName := flag.String("bucket", "", "The S3 bucket name.")
	filePath := flag.String("file", "", "The path to the file to upload.")
	s3Key := flag.String("key", "", "The S3 object key (path within the bucket). Defaults to the filename.")
	region := flag.String("region", "", "The AWS region for the S3 bucket. (e.g., us-east-1). If not set, SDK will try to determine it.")
	acl := flag.String("acl", "", "Optional: Canned ACL for the object (e.g., private, public-read).")

	flag.Parse()

	// Validate required flags
	if *bucketName == "" || *filePath == "" {
		log.Println("Bucket name and file path are required.")
		flag.Usage()
		os.Exit(1)
	}

	// Open the file to upload
	file, err := os.Open(*filePath)
	if err != nil {
		log.Fatalf("Failed to open file %q, %v", *filePath, err)
	}
	defer file.Close()

	// Determine the S3 object key
	objectKey := *s3Key
	if objectKey == "" {
		objectKey = filepath.Base(*filePath) // Use the filename as the key if not specified
	}

	// Load AWS configuration. The SDK will use the default credential chain.
	// You can specify a region explicitly or let the SDK try to determine it.
	var cfg aws.Config
	if *region != "" {
		cfg, err = config.LoadDefaultConfig(context.TODO(), config.WithRegion(*region))
	} else {
		cfg, err = config.LoadDefaultConfig(context.TODO())
	}
	if err != nil {
		log.Fatalf("Unable to load SDK config, %v", err)
	}

	// Create an S3 client
	s3Client := s3.NewFromConfig(cfg)

	// Create an S3 Uploader. This is preferred for larger files as it handles multipart uploads.
	uploader := manager.NewUploader(s3Client)

	// Prepare the PutObjectInput
	putInput := &s3.PutObjectInput{
		Bucket: aws.String(*bucketName),
		Key:    aws.String(objectKey),
		Body:   file,
	}

	// Set ACL if provided
	if *acl != "" {
		// Validate ACL value or convert to types.ObjectCannedACL
		// For simplicity, directly casting. In a real app, validate.
		putInput.ACL = types.ObjectCannedACL(*acl)
	}

	// Perform the upload
	log.Printf("Uploading %s to s3://%s/%s in region %s...\n", *filePath, *bucketName, objectKey, cfg.Region)
	result, err := uploader.Upload(context.TODO(), putInput)
	if err != nil {
		log.Fatalf("Failed to upload file to S3, %v", err)
	}

	log.Printf("Successfully uploaded %q to S3 location: %s\n", objectKey, result.Location)
}
