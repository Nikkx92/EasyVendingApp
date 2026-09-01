package grpc

import (
	"context"
	"sync/atomic"

	"google.golang.org/grpc"
)

var connectStatus atomic.Bool

func IsConnected() bool {
	return connectStatus.Load()
}

func statusInterceptor(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {

	err := invoker(ctx, method, req, reply, cc, opts...)

	connectStatus.Store(err == nil)

	return err
}
