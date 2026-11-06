package grpcserver

import (
	"context"
	"fmt"
	"log"
	"net"
	"runtime/debug"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

const shutdownTimeout = 10 * time.Second

func Run(ctx context.Context, addr string, register func(grpc.ServiceRegistrar)) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", addr, err)
	}

	srv := grpc.NewServer(grpc.ChainUnaryInterceptor(recoverPanics))

	healthSrv := health.NewServer()
	healthpb.RegisterHealthServer(srv, healthSrv)
	reflection.Register(srv)
	register(srv)

	serveErr := make(chan error, 1)
	go func() {
		log.Printf("gRPC server listening on %s", addr)
		serveErr <- srv.Serve(lis)
	}()

	select {
	case err := <-serveErr:
		return fmt.Errorf("grpc server: %w", err)
	case <-ctx.Done():
		log.Println("Shutdown signal received, draining connections...")
	}

	healthSrv.Shutdown()

	stopped := make(chan struct{})
	go func() {
		srv.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
	case <-time.After(shutdownTimeout):
		log.Println("Graceful stop timed out, forcing stop")
		srv.Stop()
	}

	log.Println("gRPC server stopped")
	return nil
}

func recoverPanics(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic in %s: %v\n%s", info.FullMethod, r, debug.Stack())
			err = status.Error(codes.Internal, "internal error")
		}
	}()
	return handler(ctx, req)
}
