package main

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	awscred "github.com/aws/aws-sdk-go/aws/credentials"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gofrs/uuid"
)

const (
	accessKey = "AKG9173KPSW1MEU7FJYT"
	secretKey = "3CuKGRgNxg3p1F6xmTdk4s94a4wmCXjEyL9mVtLc"
	endpoint  = "http://10.9.8.123:32615"
)

func main() {
	var (
		bucket   = "mybucket-2"
		region   = "us-east-1"
		prefix1  = "go"
		prefix2  = "s3"
		enabledS = "Enabled"
	)

	creds := awscred.NewStaticCredentials(accessKey, secretKey, "")
	// cfg := aws.NewConfig().WithEndpoint(endpoint).
	// 	WithS3ForcePathStyle(true).WithDisableSSL(true).
	// 	WithCredentials(creds).WithRegion(region)

	cfg := &aws.Config{
		Region:           aws.String(region),
		Endpoint:         aws.String(endpoint),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      creds,
		DisableSSL:       aws.Bool(true),
	}
	session, err := awssession.NewSession(cfg)
	if err != nil {
		log.Fatal(err)
	}
	client := s3.New(session)

	in := &s3.ListObjectsV2Input{}
	in.SetBucket(bucket).SetDelimiter("/").SetPrefix("")
	out, err := client.ListObjectsV2(in)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(out.CommonPrefixes)
	for _, content := range out.Contents {
		fmt.Printf("content: %s\n", content)
	}

	// os.Exit(1)

	// lifecycle
	lcc, err := client.GetBucketLifecycleConfiguration(&s3.GetBucketLifecycleConfigurationInput{Bucket: &bucket})
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("********lcc: %#v\n", lcc)

	now := time.Now()
	fmt.Println(now.Format("2006-01-02"))
	date1 := time.Date(now.Year(), now.Month(), 0, 0, 0, 0, 0, time.UTC)
	date2 := time.Date(now.Year()+1, now.Month(), 0, 0, 0, 0, 0, time.UTC)
	// time.Now().UTC().Format("2006-01-02T15:04:05Z07:00")
	// var days int64 = 4
	expiration1 := &s3.LifecycleExpiration{}
	expiration1.SetDate(date1)
	rule1 := &s3.LifecycleRule{
		Status: aws.String(enabledS),
		ID:     aws.String(UUID4()),
		Filter: &s3.LifecycleRuleFilter{
			Prefix: aws.String(prefix1),
		},
		Expiration: expiration1,
	}
	input := &s3.PutBucketLifecycleConfigurationInput{
		Bucket: &bucket,
		LifecycleConfiguration: &s3.BucketLifecycleConfiguration{
			Rules: []*s3.LifecycleRule{
				rule1,
				&s3.LifecycleRule{
					Status: aws.String(enabledS),
					ID:     aws.String(UUID4()),
					Filter: &s3.LifecycleRuleFilter{
						Prefix: aws.String(prefix2),
					},
					Expiration: &s3.LifecycleExpiration{
						Date: aws.Time(date2),
					},
				},
			},
		},
	}
	fmt.Printf("lifecycle input: %#v\n", input)
	// lcOutput, err := client.PutBucketLifecycleConfigurationWithContext(context.TODO(), input)
	lcOutput, err := client.PutBucketLifecycleConfiguration(input)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("lcc output: %#v\n", lcOutput)

	lcc, err = client.GetBucketLifecycleConfiguration(&s3.GetBucketLifecycleConfigurationInput{Bucket: &bucket})
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("lcc: %#v\n", lcc)
}

func UUID4() string {
	return uuid.Must(uuid.NewV4()).String()
}
