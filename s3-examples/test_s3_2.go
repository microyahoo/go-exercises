package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	// "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
)

func main() {
	const (
		// accessKey = "8EUE3A2L1KKQRMPTR79S"
		// secretKey = "GvMG155EfuT2E4DscaEm5bh9vbjZHLDMw5sTdQcJ"
		// bucket    = "zhengliang-test"
		// endpoint  = "http://10.3.8.32:80"

		// accessKey = "minioadmin"
		// secretKey = "minioadmin"
		// bucket    = "bucket"
		// endpoint  = "http://10.24.32.220:9000"

		// accessKey = "test"
		// secretKey = "test1"
		// bucket    = "zhengliang"
		// endpoint  = "http://10.9.8.72:80"

		accessKey = "CITSRGRDIGAZF7KO6LXA"
		secretKey = "hq5uc45HdISk6euyNwPDaXX7PpUzHzrRvVuAts3N"
		bucket    = "01160027"
		endpoint  = "http://10.9.8.113:7480"
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

	mgr := manager.NewDownloader(client, func(mgr *manager.Downloader) {
		mgr.PartSize = 30 * 1024 * 1024
	})
	// buf := make([]byte, int(headObject.ContentLength))
	// w := manager.NewWriteAtBuffer(buf)
	for _, key := range keys {
		head, err := client.HeadObject(context.TODO(), &s3.HeadObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
		if err != nil {
			log.Fatal(err)
		}
		length := head.ContentLength
		buf := make([]byte, length)
		w := manager.NewWriteAtBuffer(buf)
		log.Printf("key: %s, content type: %s, content length: %d", key, *head.ContentType, length)
		if length == 0 {
			fmt.Println("****")
			n, err := mgr.Download(context.Background(), w, &s3.GetObjectInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(key),
			})
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("##****download %d", n)
			break
		}
	}

}

// https://sdk.amazonaws.com/java/api/latest/software/amazon/awssdk/transfer/s3/internal/DefaultS3TransferManager.html
