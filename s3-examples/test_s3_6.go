package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"

	// "github.com/aws/aws-sdk-go-v2/service/internal/s3shared"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/smithy-go"
)

func main() {
	const (
		accessKey = "testy"
		secretKey = "testy"
		bucket    = "test"
		endpoint  = "http://10.9.8.72:80"
	)
	appCreds := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
		accessKey, secretKey, ""))
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(appCreds),
		// config.WithRegion("us-east-1"),
		config.WithClientLogMode(aws.LogResponseWithBody|aws.LogRequestWithBody),
		config.WithEndpointResolver(
			aws.EndpointResolverFunc(
				func(service, region string) (aws.Endpoint, error) {
					return aws.Endpoint{
						Source: aws.EndpointSourceCustom,
						URL:    endpoint,
					}, nil
				})),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Create an Amazon S3 service client
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	// file, err := os.Create("hhhhhhh")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// // put object
	// poo, err := client.PutObject(context.Background(), &s3.PutObjectInput{
	// 	Bucket: aws.String(bucket),
	// 	Key:    aws.String("hhhhhhh"),
	// 	Body:   file,
	// })

	output, err := client.HeadObject(context.Background(), &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String("hhhhhhh"),
	})
	if err != nil {
		log.Printf("%T\n", err) // *smithy.OperationError

		// https://github.com/aws/aws-sdk-go-v2/issues/2084
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			log.Printf("S3 error: %T, %s, %s\n", apiErr, apiErr.ErrorCode(), apiErr.ErrorMessage()) //  *smithy.GenericAPIError
			switch apiErr.(type) {
			case *types.NoSuchKey, *types.NotFound, *types.NoSuchBucket:
				log.Printf("no such key or not found: %T\n", err)
			}
			if apiErr.ErrorCode() == "NotFound" {
				log.Printf("no such key or not found: %T\n", err)
			}
		}
		// var notFoundErr *types.NotFound
		// if errors.As(err, &notFoundErr) {
		// 	log.Printf("*type of error: %T\n", notFoundErr)
		// } else {
		// 	log.Printf("**type of error: %T\n", notFoundErr)
		// }
		var oe *smithy.OperationError
		if errors.As(err, &oe) {
			log.Printf("type of error: %T\n", errors.Unwrap(err))
			log.Printf("type of error: %T\n", errors.Unwrap(oe))
			// var x *s3shared.ResponseError
			// if errors.As(errors.Unwrap(err), &x) {
			// 	log.Printf("type of error: %T, %s\n", x, x)
			// }

			log.Printf("type of error: %T\n", oe.Unwrap())
			log.Fatal(oe.Unwrap())
		}
	}
	fmt.Printf("head object output: %s\n", awsutil.Prettify(output))
}
