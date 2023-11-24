package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {
	bucket := "test"
	key := "core"
	filePath := "core"

	file, err := os.Open(filePath)
	if err != nil {
		exitErrorf("Unable to open file %q, %v", filePath, err)
	}

	defer file.Close()

	sess, _ := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials("testy", "testy", ""),
		Endpoint:         aws.String("http://10.9.8.95:80"),
		Region:           aws.String("us-east-1"),
		S3ForcePathStyle: aws.Bool(true),
	})

	svc := s3.New(sess)

	// Step 1: Start a new multipart upload.
	createResp, err := svc.CreateMultipartUpload(&s3.CreateMultipartUploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		exitErrorf("Unable to start multipart upload %q, %v", key, err)
	}

	var completedParts []*s3.CompletedPart
	partNumber := int64(1)
	buffer := make([]byte, 5*1024*1024) // 5MB part size

	for {
		n, err := file.Read(buffer)
		if n > 0 {
			// Step 2: Upload parts.
			uploadResp, err := svc.UploadPart(&s3.UploadPartInput{
				Bucket:     aws.String(bucket),
				Key:        aws.String(key),
				PartNumber: aws.Int64(partNumber),
				UploadId:   createResp.UploadId,
				Body:       bytes.NewReader(buffer[:n]),
			})
			if err != nil {
				exitErrorf("Failed to upload part %v, %v", partNumber, err)
			}
			completedParts = append(completedParts, &s3.CompletedPart{
				ETag:       uploadResp.ETag,
				PartNumber: aws.Int64(partNumber),
			})
			partNumber++
		}

		if err != nil {
			if err == io.EOF {
				break
			}
			exitErrorf("Failed to read file %q, %v", filePath, err)
		}
	}

	// Step 3: Complete multipart upload.
	_, err = svc.CompleteMultipartUpload(&s3.CompleteMultipartUploadInput{
		Bucket:   aws.String(bucket),
		Key:      aws.String(key),
		UploadId: createResp.UploadId,
		MultipartUpload: &s3.CompletedMultipartUpload{
			Parts: completedParts,
		},
	})
	if err != nil {
		exitErrorf("Failed to complete multipart upload %q, %v", key, err)
	}

	fmt.Printf("Successfully uploaded %q to %q\n", filePath, bucket)
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
