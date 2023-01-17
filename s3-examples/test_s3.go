package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go/aws/awsutil"
)

func main() {
	const (
		accessKey = "WF48U4VBV3D0KOIWT9JV"
		secretKey = "eJCzfGTO4dETlbQBciEEZmPB0U3adK8V22LLdr8Q"
		bucket    = "zhengliang"
		endpoint  = "http://127.0.0.1:80"
	)
	appCreds := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
		accessKey, secretKey, ""))
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(appCreds),
		config.WithRegion("us-east-1"),
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

	// lifecycle
	ctx := context.Background()
	lcc, err := client.GetBucketLifecycleConfiguration(ctx,
		&s3.GetBucketLifecycleConfigurationInput{
			Bucket: aws.String(bucket),
		})
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("get bucket lifecycle configuration: %s\n", awsutil.Prettify(lcc))

	input := &s3.PutBucketLifecycleConfigurationInput{
		Bucket: aws.String(bucket),
		LifecycleConfiguration: &types.BucketLifecycleConfiguration{
			Rules: []types.LifecycleRule{
				types.LifecycleRule{
					Status: types.ExpirationStatusEnabled,
					ID:     aws.String("rule-uuid1"),
					Filter: &types.LifecycleRuleFilterMemberAnd{
						Value: types.LifecycleRuleAndOperator{
							Prefix: aws.String("prefix"),
							Tags: []types.Tag{
								{
									Key:   aws.String("color"),
									Value: aws.String("red"),
								},
								{
									Key:   aws.String("color"),
									Value: aws.String("blue"),
								},
							},
						},
					},
					Expiration: &types.LifecycleExpiration{
						Days: 2,
					},
				},
			},
		},
	}
	fmt.Printf("lifecycle input: %s\n", awsutil.Prettify(input))
	lccOutput, err := client.PutBucketLifecycleConfiguration(ctx, input)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("lifecycle configuration output: %s\n", awsutil.Prettify(lccOutput))

	lcc, err = client.GetBucketLifecycleConfiguration(ctx,
		&s3.GetBucketLifecycleConfigurationInput{
			Bucket: aws.String(bucket),
		})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("get lifecycle configuration: %s\n", awsutil.Prettify(lcc))
}
