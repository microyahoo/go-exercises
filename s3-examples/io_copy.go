package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func createS3Client(ak, sk, endpoint string) *s3.Client {
	dialer := &net.Dialer{
		Control: func(network, address string, conn syscall.RawConn) error {
			var operr error
			if err := conn.Control(func(fd uintptr) {
				operr = syscall.SetsockoptInt(int(fd), syscall.IPPROTO_TCP, syscall.TCP_QUICKACK, 1)
				// operr = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.TCP_QUICKACK, 1)
				// operr = syscall.SetsockoptInt(int(fd), syscall.SOL_TCP, syscall.TCP_QUICKACK, 1)
			}); err != nil {
				return err
			}
			return operr
		},
	}
	tr := &http.Transport{
		DialContext: dialer.DialContext,
	}
	hc := &http.Client{Transport: tr}
	appCreds := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
		ak, sk, ""))
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(appCreds),
		config.WithRegion("us-east-1"),
		config.WithHTTPClient(hc),
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
	return client
}

func main() {
	var (
		accessKey = "testy"
		secretKey = "testy"
		bucket    = "test"
		endpoint  = "http://10.3.9.221:80"
	)

	keys := []string{"1k", "4k", "8k", "16k", "32k", "64k", "128k", "256k", "512k"}
	size := int64(1024 * 1)
	for i, key := range keys {
		client := createS3Client(accessKey, secretKey, endpoint)
		ctx := context.Background()
		result, err := client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: &bucket,
			Key:    &key,
		})
		if err != nil {
			log.Fatal(err)
		}
		start := time.Now()
		numBytes, err := io.Copy(io.Discard, result.Body)
		if err != nil {
			log.Fatal(err)
		}
		result.Body.Close()
		duration := time.Since(start)
		fmt.Println(key, duration)
		if numBytes != size {
			log.Fatalf("size %d not match", numBytes)
		}
		if i == 0 {
			size *= 4
		} else {
			size *= 2
		}
	}
}
