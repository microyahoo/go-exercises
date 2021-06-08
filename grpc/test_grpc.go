package main

import (
	"context"
	"log"
	"time"

	grpccodes "google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
)

func call() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.Req{
		Arg: 1,
	}
	resp, err := client.CallSomething(ctx, req)
	if err != nil {
		if rpcErr, ok := grpcstatus.FromError(err); ok {
			if rpcErr.Code() == grpccodes.Canceled ||
				rpcErr.Code() == grpccodes.DeadlineExceeded {
				log.Warnf("RPC call timeout: %s", err)
				return nil
			}
		}
		return err
	}

	return nil
}
