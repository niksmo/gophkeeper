package grpccl

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient[T any] struct {
	addr     string
	conn     *grpc.ClientConn
	clFabric func(grpc.ClientConnInterface) T
	client   T
}

func New[T any](
	addr string, clFabric func(grpc.ClientConnInterface) T,
) *GRPCClient[T] {
	client := &GRPCClient[T]{addr: addr, clFabric: clFabric}
	return client
}

func (c *GRPCClient[T]) ConnClient() (T, error) {
	dialOpt := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.NewClient(c.addr, dialOpt)
	if err != nil {
		return c.client, fmt.Errorf("failed to connect addr: %s", c.addr)
	}
	c.conn = conn
	return c.clFabric(conn), nil
}

func (c *GRPCClient[T]) Close() error {
	return c.conn.Close()
}
