package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"

	// "github.com/aws/aws-sdk-go-v2/service/internal/s3shared"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws/awsutil"
)

func main() {
	const (
		key       = "install_client.sh"
		accessKey = "simulation-developer"
		secretKey = "Q6nuvt7yA70cRZpn"
		bucket    = "prophet-client"
		endpoint  = "http://10.3.8.31:80"
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

	output, err := client.HeadObject(context.Background(), &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Printf("%T\n", err) // *smithy.OperationError
	}
	fmt.Printf("head object output: %s\n", awsutil.Prettify(output))
}
