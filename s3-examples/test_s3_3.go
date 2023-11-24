package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	// "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
)

func main() {
	const (
		accessKey = "testy"
		secretKey = "testy"
		bucket    = "test"
		endpoint  = "http://10.9.8.95:80"
	)
	appCreds := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
		accessKey, secretKey, ""))
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(appCreds),
		config.WithRegion("us-east-1"),
		// config.WithClientLogMode(aws.LogResponseWithBody|aws.LogRequestWithBody),
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

	params := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		// MaxKeys: 5000,
	}

	// Get the first page of results for ListObjectsV2 for a bucket
	output, err := client.ListObjectsV2(context.TODO(), params)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("max keys: %d, key count: %d", output.MaxKeys, output.KeyCount)

	var (
		totalObjects = 0
		keys         []string
	)
	paginator := s3.NewListObjectsV2Paginator(client, params)
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		totalObjects += len(output.Contents)
		for _, content := range output.Contents {
			keys = append(keys, *content.Key)
		}
	}
	fmt.Println("total objects:", totalObjects, ", key size: ", len(keys))
}

// https://sdk.amazonaws.com/java/api/latest/software/amazon/awssdk/transfer/s3/internal/DefaultS3TransferManager.html
