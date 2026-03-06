package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func main() {
	const (
		accessKey = "CHHIEDO3WL59C6JPDA05"
		secretKey = "kpLoIuou06K7OlBgmuUVPibJsC97lVCRub8a2AQI"
		bucket    = "argo-workflow-artifacts"
		endpoint  = "http://10.3.8.34:80"
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

	key := "localization-benchmark-v2-12752-1760420800/localization-benchmark-v2-12752-1760420800-single-2536172718/main.log"
	output, err := client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("head object: %+v\n, %s", output, *output.ETag)

	params := &s3.GetObjectAttributesInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		ObjectAttributes: []types.ObjectAttributes{
			types.ObjectAttributesEtag,
			types.ObjectAttributesChecksum,
			types.ObjectAttributesObjectParts,
		},
	}

	output2, err := client.GetObjectAttributes(context.TODO(), params)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("etag: %+v, checksum: %+v, attributes: %+v", output2.ETag, output2.Checksum, output2)
}
