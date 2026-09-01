package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

func timeOut(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return invoker(timeoutCtx, method, req, reply, cc, opts...)
}
