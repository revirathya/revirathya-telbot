package grpc

import (
	"fmt"

	"google.golang.org/grpc"
)

func dial(name string, addr string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(addr, defaultDialOptions()...)
	if err != nil {
		return nil, fmt.Errorf("dial %s (%s): %w", name, addr, err)
	}
	return conn, nil
}
