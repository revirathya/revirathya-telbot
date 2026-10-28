package withtelegramid

import (
	"context"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const telegramIDKey = "x-telegram-user-id"

type TelegramIDLookup func(context.Context, any) (int64, error)

func UnaryClientInterceptor(lookupFunc TelegramIDLookup) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		telegramID, err := lookupFunc(ctx, req)
		if err != nil {
			return err
		}

		ctx = metadata.AppendToOutgoingContext(ctx,
			telegramIDKey, strconv.FormatInt(telegramID, 10),
		)

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
