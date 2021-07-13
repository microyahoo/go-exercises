// https://github.com/grpc/grpc-go/issues/2767
// 1. run case 1
// go test -run "TestKeepaliveMaxAge/case1-Fatal-Unavailable-transport is closing" -count 100
// result: test fail
//
package testcase1

import (
	"context"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/examples/route_guide/routeguide"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/keepalive"
)

type routeGUIDService struct {
	routeguide.UnimplementedRouteGuideServer
}

func (*routeGUIDService) GetFeature(ctx context.Context, p *routeguide.Point) (*routeguide.Feature, error) {
	return &routeguide.Feature{Name: "hello world"}, nil
}

func (*routeGUIDService) ListFeatures(p *routeguide.Rectangle, s routeguide.RouteGuide_ListFeaturesServer) error {
	return nil
}

func (*routeGUIDService) RecordRoute(s routeguide.RouteGuide_RecordRouteServer) error {
	return nil
}

func (*routeGUIDService) RouteChat(s routeguide.RouteGuide_RouteChatServer) error {
	go func() {
		err := s.Send(&routeguide.RouteNote{Message: "server send"})
		if err != nil {
			log.Printf("failed to send:%s", err)
		}
	}()
	for {
		_, err := s.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			log.Printf("failed to recv:%s", err)
			return err
		}
	}
}

// go test -run "TestKeepaliveMaxAge/case1-Fatal-Unavailable-transport is closing" -count 100
func TestKeepaliveMaxAge(t *testing.T) {
	// export GRPC_GO_LOG_VERBOSITY_LEVEL=99 GRPC_GO_LOG_SEVERITY_LEVEL=info
	grpclog.SetLoggerV2(grpclog.NewLoggerV2WithVerbosity(os.Stderr, os.Stderr, os.Stderr, 100))
	keepAliveMaxConnectionAge := time.Second * 5

	// 1. start server
	s := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionAge: keepAliveMaxConnectionAge,
		}),
	)
	routeguide.RegisterRouteGuideServer(s, &routeGUIDService{})
	const addr = "127.0.0.1:8000"
	const pprofAddr = "127.0.0.1:8001"
	go func() {
		_ = pprof.Index
		err := http.ListenAndServe(pprofAddr, nil)
		if err != nil && err != http.ErrServerClosed {
			return
		}
	}()
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		t.Fatal(err)
		return
	}
	go func() {
		err := s.Serve(ln)
		if err != nil {
			t.Error(err)
			return
		}
	}()
	defer func() {
		s.GracefulStop()
	}()
	// wait server start up
	time.Sleep(time.Millisecond * 100)
	// 2. start client
	cc, err := grpc.DialContext(context.Background(), addr, grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
		return
	}
	// go test -run "TestKeepaliveMaxAge/case1-Fatal-Unavailable-transport is closing" -count 100
	t.Run("case1-Fatal-Unavailable-transport is closing", func(t *testing.T) {
		timer := time.NewTimer(keepAliveMaxConnectionAge + time.Second)
		for {
			_, err := routeguide.NewRouteGuideClient(cc).GetFeature(context.Background(), &routeguide.Point{})
			if err != nil {
				// unexpected: its Fatal
				// must use grpc.FailFast(false), it's ok
				t.Fatal(err)
				return
			}
			// log.Println(feature)
			select {
			case <-timer.C:
				return
			default:
			}
		}
	})
}
