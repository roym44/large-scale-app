package common

import (
	"fmt"

	"google.golang.org/grpc"
)

type ServiceClientBase[client_t any] struct {
	Address      string
	CreateClient func(grpc.ClientConnInterface) client_t
}

func (obj *ServiceClientBase[client_t]) Connect() (res client_t, closeFunc func(), err error) {
	conn, err := grpc.Dial(obj.Address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		var empty client_t
		return empty, nil, fmt.Errorf("failed to connect client to %v: %v", obj.Address, err)
	}
	c := obj.CreateClient(conn)
	return c, func() { conn.Close() }, nil
}
