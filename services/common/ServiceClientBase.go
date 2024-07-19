package common

import (
	"fmt"
	"math/rand"

	RegServiceClient "github.com/TAULargeScaleWorkshop/RLAD/services/reg-service/client"
	"google.golang.org/grpc"
)

type ServiceClientBase[client_t any] struct {
	RegistryAddresses []string // to connect to the registry service
	ServiceName       string   // to disover our server's nodes
	CreateClient      func(grpc.ClientConnInterface) client_t
}

// randomly picks a service node address to connect to
func (obj *ServiceClientBase[client_t]) pickNode() string {
	regClient := RegServiceClient.NewRegServiceClient(obj.RegistryAddresses)
	nodes, err := regClient.Discover(obj.ServiceName)
	if err != nil || len(nodes) == 0 {
		return ""
	}
	randomIndex := rand.Intn(len(nodes))
	return nodes[randomIndex]
}

func (obj *ServiceClientBase[client_t]) Connect() (res client_t, closeFunc func(), err error) {
	// pick some node of the service
	node_address := obj.pickNode()
	if node_address == "" {
		var empty client_t
		return empty, nil, fmt.Errorf("no available nodes found")
	}

	// connect
	conn, err := grpc.Dial(node_address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		var empty client_t
		return empty, nil, fmt.Errorf("failed to connect client to %v: %v", node_address, err)
	}
	c := obj.CreateClient(conn)
	return c, func() { conn.Close() }, nil
}
