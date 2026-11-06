package grpcclient

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const serviceConfig = `{"loadBalancingConfig": [{"round_robin": {}}]}`

func New(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	defaults := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(serviceConfig),
	}

	conn, err := grpc.NewClient(target, append(defaults, opts...)...)
	if err != nil {
		return nil, fmt.Errorf("create grpc client for %s: %w", target, err)
	}
	return conn, nil
}
