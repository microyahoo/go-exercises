package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Type int

const (
	Read Type = 1 << iota
	Write
	Delete
	List
)

func (t Type) String() string {
	var types []string
	if t&Read != 0 {
		types = append(types, "read")
	}
	if t&Write != 0 {
		types = append(types, "write")
	}
	if t&Delete != 0 {
		types = append(types, "delete")
	}
	if t&List != 0 {
		types = append(types, "list")
	}
	return strings.Join(types, ",")
}

func createClient(ak, sk, endpoint string) *s3.Client {
	dialer := &net.Dialer{
		Control: func(network, address string, conn syscall.RawConn) error {
			var operr error
			if err := conn.Control(func(fd uintptr) {
				operr = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.TCP_QUICKACK, 1)
			}); err != nil {
				return err
			}
			return operr
		},
	}
	tr := &http.Transport{
		// dial tcp 10.3.9.232:80: connect: cannot assign requested address
		// https://github.com/golang/go/issues/16012
		// https://github.com/valyala/fasthttp/blob/d795f13985f16622a949ea9fc3459cf54dc78b3e/client_timing_test.go#L292
		// https://github.com/etcd-io/etcd/blob/55de68d18c63fde8747f2cd9f7c2ff242346f756/client/pkg/transport/timeout_transport.go#L39
		// overrides the DefaultMaxIdleConnsPerHost = 2
		MaxIdleConnsPerHost: 1024, // allow more idle connections between peers to avoid unnecessary port allocation

		// https://gitlab.com/gitlab-org/gitlab-pages/-/merge_requests/274#note_335755221
		//
		// Not setting this results in connections in the pool staying open forever. When the pool size is 2,
		// this is unlikely to matter much, as we're highly unlikely to leave those connections idle.
		// But with a higher number of connections in the pool, we'll probably want to set this so that after a burst of activity,
		// we can close idle connections. Perhaps a value like 15 minutes for benchmark test?
		//
		// https://github.com/yarpc/yarpc-go/issues/1906
		// https://blog.cloudflare.com/when-tcp-sockets-refuse-to-die/
		IdleConnTimeout: 15 * time.Minute,
		DialContext:     dialer.DialContext,
		// Set TCP_NODELAY
		// DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
		// 	conn, err := (&net.Dialer{}).DialContext(ctx, network, addr)
		// 	if tc, ok := conn.(*net.TCPConn); ok {
		// 		tc.SetNoDelay(true)
		// 		rc, e := tc.SyscallConn()
		// 		if e == nil {
		// 			rc.Control(func(fd uintptr) {
		// 				syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.TCP_QUICKACK, 1)
		// 			})
		// 		}
		// 	}
		// 	return conn, err
		// },
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

	keys := []string{"4k", "8k", "16k", "32k", "64k", "128k", "256k", "512k"}
	size := int64(1024 * 4)
	for _, key := range keys {
		client := createClient(accessKey, secretKey, endpoint)
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
		size *= 2
	}
	// fmt.Println("new discarder")
	// size = int64(1024 * 4)
	// for _, key := range keys {
	// 	ctx := context.Background()
	// 	result, err := client.GetObject(ctx, &s3.GetObjectInput{
	// 		Bucket: &bucket,
	// 		Key:    &key,
	// 	})
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	start := time.Now()
	// 	// d := &discarder{}
	// 	// numBytes, err := d.ReadFrom(result.Body)
	// 	// numBytes, err := io.Copy(&discarder{}, result.Body)
	// 	numBytes, err := io.CopyN(&discarder{}, result.Body, size)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	duration := time.Since(start)
	// 	fmt.Println(key, duration)
	// 	if numBytes != size {
	// 		log.Fatalf("size %d not match", numBytes)
	// 	}
	// 	size *= 2
	// }

	// fmt.Println(Read)
	// fmt.Println(Write)
	// fmt.Println(Delete)
	// fmt.Println(List)
	// t := Read | Write | List
	// fmt.Println(t)
	// t = Delete | Write | List
	// fmt.Println(t)

	// d := (111111 * time.Second) + (222222 * time.Millisecond) + (3333 * time.Microsecond) + (333 * time.Nanosecond)
	// durStr, _ := durationfmt.Format(d, "%0hh:%0mm:%0ss")
	// fmt.Println(durStr)
	// durStr, _ = durationfmt.Format(d, "%0s")
	// fmt.Println(durStr)
	// fmt.Println(d.Round(time.Second))
}

// type discarder struct {
// }

// // func (d *discarder) ReadFrom(r io.Reader) (int64, error) {
// // 	p, _ := io.ReadAll(r)
// // 	return int64(len(p)), nil
// // }

// func (d *discarder) Write(p []byte) (n int, err error) {
// 	return len(p), nil
// }

// func (d *discarder) WriteAt(p []byte, off int64) (n int, err error) {
// 	return len(p), nil
// }
