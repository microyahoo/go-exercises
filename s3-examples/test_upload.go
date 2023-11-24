package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func RetryUpload(svc *s3.S3, bucket, key string) error {
	var (
		byteSize int = 101 * 1024 * 1024 // 101M
		// byteSize       int = 26 * 1024 * 1024 // 26M
		completedParts []*s3.CompletedPart
	)
	buffer := generateBytes(byteSize)

	fmt.Println("start CreateMultipartUpload")
	resp, err := svc.CreateMultipartUpload(&s3.CreateMultipartUploadInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		fmt.Println("failed to CreateMultipartUpload ")
		return err
	}
	log.Println("success to CreateMultipartUpload")

	log.Println("start UploadPart PartNumber 1")
	// first update 8 MB
	fileBytes := buffer[0 : 8*1024*1024]
	uploadResult, err := svc.UploadPart(&s3.UploadPartInput{
		Body:          bytes.NewReader(fileBytes),
		Bucket:        &bucket,
		Key:           &key,
		PartNumber:    aws.Int64(int64(1)),
		UploadId:      resp.UploadId,
		ContentLength: aws.Int64(int64(len(fileBytes))),
	})

	if err != nil {
		log.Println("failed to UploadPart PartNumber 1")
		return err
	}

	etag := uploadResult.ETag
	log.Printf("etag: %s\n", *etag)
	completedParts = append(completedParts, &s3.CompletedPart{
		ETag:       etag,
		PartNumber: aws.Int64(int64(1)),
	})
	log.Println("success to UploadPart PartNumber 1")

	var wg sync.WaitGroup
	wg.Add(3)
	f := func(bytes io.ReadSeeker) {
		defer wg.Done()

		var buf = make([]byte, 64)
		var stk = buf[:runtime.Stack(buf, false)]
		log.Printf("start UploadPart PartNumber 2, goroutine id: %s\n", string(stk))
		// second part
		uploadResult2, err := svc.UploadPart(&s3.UploadPartInput{
			Body:          bytes,
			Bucket:        &bucket,
			Key:           &key,
			PartNumber:    aws.Int64(int64(2)),
			UploadId:      resp.UploadId,
			ContentLength: aws.Int64(int64(100 * 1024 * 1024)), // 100M
			// ContentLength: aws.Int64(int64(16 * 1024 * 1024)), // 16M
		})
		if err != nil {
			log.Printf("Failed to UploadPart PartNumber 2, goroutine id: %s, error: %s\n", string(stk), err.Error())
			return
		}
		log.Printf("success to UploadPart PartNumber 2, goroutine id: %s\n", string(stk))
		completedParts = append(completedParts, &s3.CompletedPart{
			ETag:       uploadResult2.ETag,
			PartNumber: aws.Int64(int64(2)),
		})

	}
	go f(bytes.NewReader(buffer[1*1024*1024:]))
	// go f(bytes.NewReader(buffer[:16*1024*1024]))
	// go f(bytes.NewReader(buffer[2*1024*1024 : 18*1024*1024]))

	go f(bytes.NewReader(buffer[2*1024*1024 : 80*1024*1024]))
	go f(bytes.NewReader(buffer[19*1024*1024 : 100*1024*1024]))

	// go f(bytes.NewReader(buffer[8*1024*1024 : 24*1024*1024]))
	// // go f(bytes.NewReader(buffer[:16*1024*1024]))
	// // go f(bytes.NewReader(buffer[2*1024*1024 : 18*1024*1024]))

	// go f(bytes.NewReader(buffer[2*1024*1024 : 8*1024*1024]))
	// go f(bytes.NewReader(buffer[19*1024*1024 : 26*1024*1024]))

	wg.Wait()

	log.Println("CompleteMultipartUpload start")
	_, err = svc.CompleteMultipartUpload(&s3.CompleteMultipartUploadInput{
		Bucket:   &bucket,
		Key:      &key,
		UploadId: resp.UploadId,
		MultipartUpload: &s3.CompletedMultipartUpload{
			Parts: completedParts,
		},
	})
	if err != nil {
		log.Println("failed to CompleteMultipartUpload")
		return err
	}
	log.Println("CompleteMultipartUpload success ")

	return nil
}

func generateBytes(size int) []byte {
	str := []byte("0123456789abcdefghijklmnopqrstuvwxyz")
	ceil := math.Ceil(float64(size) / (1024 * 1024))
	var generateBytes []byte
	for i := 0; i < int(ceil); i++ {
		j := i
		if j >= len(str) {
			j = int(math.Mod(float64(i), float64(len(str))))
		}
		generateBytes = append(generateBytes, bytes.Repeat([]byte{str[j]}, 1024*1024)...)
	}
	return generateBytes
}

// Read(p []byte) (n int, err error)
func generateRandomBytes(size int) ([]byte, error) {
	rndSrc := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(rndSrc)
	random := make([]byte, size)
	var (
		err   error
		n     = 0
		retry = 3
	)
	for i := 0; i < retry && n != size; i++ {
		n, err = io.ReadFull(rng, random)
	}
	if err != nil {
		return nil, err
	}
	return random, nil
}

func main() {
	// str1 := string(generateBytes(10))
	// str2 := string(generateBytes(1024*1024 + 1))
	// fmt.Println(string(str2))
	// fmt.Println(len(str1))
	// fmt.Println(len(str2))

	svc := s3.New(session.Must(session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials("testy", "testy", ""),
		Endpoint:         aws.String("http://10.9.8.95:80"),
		Region:           aws.String("us-east-1"),
		S3ForcePathStyle: aws.Bool(true),
	})))
	err := RetryUpload(svc, "test", "test-upload")
	if err != nil {
		log.Panic(err.Error())
	}
}
