package main

import (
	"context"
	"fmt"

	"github.com/0x5010/grpcp"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

func main() {
	var addr, name string

	pool := grpcp.New(func(addr string) (*grpc.ClientConn, error) {
		return grpc.Dial(
			addr,
			grpc.WithInsecure(),
		)
	})
	conn, _ := pool.GetConn(addr)

	client := pb.NewGreeterClient(conn)
	r, _ := client.SayHello(context.Background(), &pb.HelloRequest{Name: name})
	fmt.Println(r.GetMessage())
}
