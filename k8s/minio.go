package main

import (
	"context"
	"log"

	argos3 "github.com/argoproj/pkg/s3"
	"github.com/minio/minio-go/v7"
)

func main() {
	opts := argos3.S3ClientOpts{
		Endpoint:        "10.9.8.72:80",
		Region:          "us-east-1",
		Secure:          false,
		AccessKey:       "testy",
		SecretKey:       "testy",
		AddressingStyle: argos3.PathStyle,
	}

	if tr, err := minio.DefaultTransport(false); err == nil {
		opts.Transport = tr
	}
	ctx := context.Background()
	client, err := argos3.NewS3Client(ctx, opts)
	if err != nil {
		panic(err)
	}
	bucket := "test"
	log.Println("start to upload empty file")
	// create an empty file first
	err1 := client.PutFile(bucket, "empty-file", "/root/go/src/github.com/microyahoo/go-exercises/k8s/empty-file")
	if err1 != nil {
		log.Printf("Failed to upload empty file: %s\n", err1)
	} else {
		log.Println("Successfully uploaded an empty file")
	}
	log.Println("start to upload empty directory")
	// create an empty directory first
	err2 := client.PutDirectory(bucket, "empty-directory", "/root/go/src/github.com/microyahoo/go-exercises/k8s/empty-dir")
	if err2 != nil {
		log.Printf("Failed to upload empty directory: %s\n", err2)
	} else {
		log.Println("Successfully uploaded an empty directory")
	}
}
