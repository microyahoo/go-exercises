package s3

import (
	"log"
	"strings"

	"aws-sdk-demo/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3 struct {
	client *s3.S3
}

func (s *S3) ListBuckets() string {
	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(
			config.AwsAccessKey,
			config.AwsAccessSecretKey,
			config.AwsAccessToken),
		// Endpoint:    aws.String(s.endpoint),
		Region: aws.String(config.AwsRegion),
		// S3ForcePathStyle: aws.Bool(true),
	}

	newSession := session.New(s3Config)
	s3Client := s3.New(newSession)

	result, err := s3Client.ListBuckets(nil)
	if err != nil {
		log.Fatalf("Unable to list buckets, %v", err)
	}

	return result.String()
}

func (s *S3) ListObjects(bucket string) (string, error) {
	out, err := s.client.ListObjects(&s3.ListObjectsInput{Bucket: aws.String(bucket)})
	return out.String(), err
}

func (*S3) CreateBucket() error {
	return nil
}

func Test() {
	bucket := aws.String("liang-bucket")
	key := aws.String("test_object")

	id := "AKIAOBL4GJUKA3PFLEAA"
	secret := "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"
	token := ""

	region := aws.String("cn-north-1")
	endpoint := aws.String("s3.cn-north-1.amazonaws.com.cn")

	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(id, secret, token),
		Endpoint:    endpoint,
		Region:      region,
		// S3ForcePathStyle: aws.Bool(true),
	}

	newSession := session.New(s3Config)
	s3Client := s3.New(newSession)
	_, err := s3Client.CreateBucket(&s3.CreateBucketInput{
		Bucket: bucket,
	})
	if err != nil {
		log.Println("Failed to create bucket", err)
		return
	}
	if err = s3Client.WaitUntilBucketExists(&s3.HeadBucketInput{Bucket: bucket}); err != nil {
		log.Printf("Failed to wait for bucket to exist %s, %s\n", bucket, err)
		return
	}
	log.Printf("Successfully created bucket %s\n", *bucket)

	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Body:   strings.NewReader("hello world storage"),
		Bucket: bucket,
		Key:    key,
	})
	if err != nil {
		log.Printf("Failed to upload object %s/%s, %s\n", *bucket, *key, err.Error())
		return
	}
	log.Printf("Successfully uploaded key %s\n", *key)

	//Get Object
	_, err = s3Client.GetObject(&s3.GetObjectInput{
		Bucket: bucket,
		Key:    key,
	})
	if err != nil {
		log.Println("Failed to download file", err)
		return
	}
	log.Printf("Successfully Downloaded key %s\n", *key)
}

func init() {
}
