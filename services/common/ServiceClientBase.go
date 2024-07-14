package common

import (
	"fmt"

	"google.golang.org/grpc"
)

type ServiceClientBase[client_t any] struct {
	RegistryAddresses []string
	CreateClient      func(grpc.ClientConnInterface) client_t
}

// randomly picks a registry node address to connect to
func pickNode() string {
	// TODO: call Discover()
	return "a"
}

func (obj *ServiceClientBase[client_t]) Connect() (res client_t, closeFunc func(), err error) {
	node_address := pickNode()
	conn, err := grpc.Dial(node_address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		var empty client_t
		return empty, nil, fmt.Errorf("failed to connect client to %v: %v", node_address, err)
	}
	c := obj.CreateClient(conn)
	return c, func() { conn.Close() }, nil
}
