package main

import (
	"net"

	"google.golang.org/grpc"

	snapshotsapi "github.com/containerd/containerd/api/services/snapshots/v1"
	"github.com/containerd/containerd/contrib/snapshotservice"
)

func main() {
	// Create a gRPC server
	rpc := grpc.NewServer()

	// Configure your custom snapshotter
	sn, _ := NewMyCoolSnapshotter()
	// Convert the snapshotter to a gRPC service,
	service := snapshotservice.FromSnapshotter(sn)
	// Register the service with the gRPC server
	snapshotsapi.RegisterSnapshotsServer(rpc, service)

	// Listen and serve
	l, _ := net.Listen("unix", "/var/my-cool-snapshooter-socket")
	_ = rpc.Serve(l)
}
