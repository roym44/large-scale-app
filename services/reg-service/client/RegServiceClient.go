package RegServiceClient

import (
	context "context"
	"fmt"

	service "github.com/TAULargeScaleWorkshop/RLAD/services/reg-service/common"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// Note: Code Duplication to prevent importing ServiceClientBase creating an import cycle
type RegServiceClient struct {
	RegistryAddresses []string
	CreateClient      func(grpc.ClientConnInterface) service.RegServiceClient
}

func (obj *RegServiceClient) Connect() (res service.RegServiceClient, closeFunc func(), err error) {
	// TODO: we currently take the first reg address
	conn, err := grpc.Dial(obj.RegistryAddresses[0], grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		var empty service.RegServiceClient
		return empty, nil, fmt.Errorf("failed to connect client to %v: %v", obj.RegistryAddresses[0], err)
	}
	c := obj.CreateClient(conn)
	return c, func() { conn.Close() }, nil
}

func NewRegServiceClient(addresses []string) *RegServiceClient {
	return &RegServiceClient{
		RegistryAddresses: addresses, CreateClient: service.NewRegServiceClient,
	}
}

func (obj *RegServiceClient) Register(service_name string, node_addr string) error {
	c, closeFunc, err := obj.Connect()
	if err != nil {
		return fmt.Errorf("could not connect to registery server: %v", err)
	}
	defer closeFunc()
	// Call the Store RPC function
	_, err = c.Register(context.Background(), &service.UpdateRegistryParameters{ServiceName: service_name, NodeAddr: node_addr})
	if err != nil {
		return fmt.Errorf("could not call Register: %v", err)
	}
	return nil
}

func (obj *RegServiceClient) Unregister(service_name string, node_addr string) error {
	c, closeFunc, err := obj.Connect()
	if err != nil {
		return fmt.Errorf("could not connect to registery server: %v", err)
	}
	defer closeFunc()
	// Call the Store RPC function
	_, err = c.Unregister(context.Background(), &service.UpdateRegistryParameters{ServiceName: service_name, NodeAddr: node_addr})
	if err != nil {
		return fmt.Errorf("could not call Unregister: %v", err)
	}
	return nil
}

func (obj *RegServiceClient) Discover(service_name string) ([]string, error) {
	c, closeFunc, err := obj.Connect()
	if err != nil {
		return nil, fmt.Errorf("could not connect to registery server: %v", err)
	}
	defer closeFunc()
	// Call the Store RPC function
	discovered, err := c.Discover(context.Background(), wrapperspb.String(service_name))
	if err != nil {
		return nil, fmt.Errorf("could not call Unregister: %v", err)
	}
	return discovered.Addresses, nil
}
