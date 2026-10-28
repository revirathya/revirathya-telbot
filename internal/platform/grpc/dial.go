package grpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	// Interceptor
	"github.com/revirathya/revirathya-telbot/internal/auth"
	"github.com/revirathya/revirathya-telbot/internal/platform/grpc/interceptor/withtelegramid"
)

const (
	keepAliveTime    = 15 * time.Second
	keepAliveTimeout = 5 * time.Second
)

func telegramIDLookup(ctx context.Context, req any) (int64, error) {
	telegramID, err := auth.TelegramIDFrom(ctx)
	if err != nil {
		return 0, err
	}
	return telegramID, nil
}

func Dial(name string, addr string, apiKey string) (*grpc.ClientConn, error) {
	// Define dial options
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(auth.NewApiKey(apiKey)),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                keepAliveTime,
			Timeout:             keepAliveTimeout,
			PermitWithoutStream: true,
		}),
		grpc.WithChainUnaryInterceptor(
			withtelegramid.UnaryClientInterceptor(telegramIDLookup),
		),
	}

	// Create new client
	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		return nil, fmt.Errorf("dial %s (%s): %w", name, addr, err)
	}
	return conn, nil
}
