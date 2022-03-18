package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/gofrs/uuid"
)

func main() {
	const (
		accessKey      = "AKG9173KPSW1MEU7FJYT"
		secretKey      = "3CuKGRgNxg3p1F6xmTdk4s94a4wmCXjEyL9mVtLc"
		endpoint       = "http://10.9.8.123:32615"
		prefix1        = "go"
		prefix2        = "s3"
		prefix3        = "he"
		prefix4        = "hi"
		prefix5        = "event"
		coldStorage    = "COLD"
		archiveStorage = "ARCHIVE"
		tagKey         = "testKey"
		tagVal         = "testVal"
	)
	bucket := "mybucket-2"
	appCreds := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
		accessKey, secretKey, ""))
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(appCreds),
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
	// TODO(zhengliang): whether need to close?

	// Get the first page of results for ListObjectsV2 for a bucket
	output, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("first page results:")
	for _, object := range output.Contents {
		log.Printf("key=%s size=%d, storageClass=%s", aws.ToString(object.Key), object.Size, object.StorageClass)
	}

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

	now := time.Now().UTC()
	date1 := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	date2 := time.Date(now.Year()+1, now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	// date3 := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, time.UTC)

	input := &s3.PutBucketLifecycleConfigurationInput{
		Bucket: &bucket,
		LifecycleConfiguration: &types.BucketLifecycleConfiguration{
			Rules: []types.LifecycleRule{
				types.LifecycleRule{
					Status: types.ExpirationStatusEnabled,
					ID:     aws.String(UUID4()),
					Filter: &types.LifecycleRuleFilterMemberPrefix{Value: prefix1},
					Expiration: &types.LifecycleExpiration{
						Date: aws.Time(date1),
					},
				},
				types.LifecycleRule{
					Status: types.ExpirationStatusEnabled,
					ID:     aws.String(UUID4()),
					Filter: &types.LifecycleRuleFilterMemberPrefix{Value: prefix2},
					Expiration: &types.LifecycleExpiration{
						Date: aws.Time(date2),
					},
				},
				types.LifecycleRule{
					Status: types.ExpirationStatusEnabled,
					ID:     aws.String(UUID4()),
					Filter: &types.LifecycleRuleFilterMemberPrefix{Value: prefix3},
					Expiration: &types.LifecycleExpiration{
						Days: 3,
					},
				},
				types.LifecycleRule{
					Status: types.ExpirationStatusEnabled,
					ID:     aws.String(UUID4()),
					Filter: &types.LifecycleRuleFilterMemberAnd{
						Value: types.LifecycleRuleAndOperator{
							Prefix:             aws.String(prefix4),
							ObjectSizeLessThan: 120000000,
						},
					},
					Transitions: []types.Transition{
						types.Transition{
							// Date:         aws.Time(date3),
							Days:         1,
							StorageClass: types.TransitionStorageClass(coldStorage),
						},
						types.Transition{
							Days:         60,
							StorageClass: types.TransitionStorageClass(archiveStorage),
						},
					},
					Expiration: &types.LifecycleExpiration{
						Days: 90,
					},
				},
				types.LifecycleRule{
					Status: types.ExpirationStatusEnabled,
					ID:     aws.String(UUID4()),
					Filter: &types.LifecycleRuleFilterMemberAnd{
						Value: types.LifecycleRuleAndOperator{
							Prefix: aws.String(prefix5),
						},
					},
					Transitions: []types.Transition{
						types.Transition{
							Days:         1,
							StorageClass: types.TransitionStorageClass(coldStorage),
						},
					},
				},
				types.LifecycleRule{
					Status: types.ExpirationStatusEnabled,
					ID:     aws.String(UUID4()),
					Filter: &types.LifecycleRuleFilterMemberAnd{
						Value: types.LifecycleRuleAndOperator{
							Tags: []types.Tag{
								types.Tag{
									Key:   aws.String(tagKey),
									Value: aws.String(tagVal),
								},
							},
						},
					},
					Transitions: []types.Transition{
						types.Transition{
							Days:         1,
							StorageClass: types.TransitionStorageClass(coldStorage),
						},
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

	// put object
	poo, err := client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String("test"),
		Body:          bytes.NewReader([]byte("PAYLOAD")),
		StorageClass:  types.StorageClass(coldStorage),
		ContentLength: 7,
		Tagging:       aws.String(fmt.Sprintf("%s=%s", tagKey, tagVal)),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("put object output: %s\n", awsutil.Prettify(poo))
}

func UUID4() string {
	return uuid.Must(uuid.NewV4()).String()
}
