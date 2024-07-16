package grpcclientlb

import (
	"errors"

	"google.golang.org/grpc"
)

type GrpcClientLB interface {
	Get() *grpc.ClientConn
}

type grpcClientLB struct {
	grpcClient []*grpc.ClientConn
	size       int
	offset     int
}

func NewGrpcClientLB(factory func() (*grpc.ClientConn, error), poolSize int) (GrpcClientLB, error) {
	if poolSize <= 0 {
		return nil, errors.New("poolSize must be greater than 0")
	}

	grpcClient := make([]*grpc.ClientConn, poolSize)
	for i := 0; i < poolSize; i++ {
		conn, err := factory()
		if err != nil {
			return nil, err
		}

		grpcClient[i] = conn
	}

	return &grpcClientLB{
		grpcClient: grpcClient,
		size:       len(grpcClient),
		offset:     0,
	}, nil
}

func (g *grpcClientLB) Get() *grpc.ClientConn {
	if g.offset >= g.size {
		g.offset = 0
	}

	conn := g.grpcClient[g.offset]
	g.offset = (g.offset + 1) % g.size

	return conn
}
