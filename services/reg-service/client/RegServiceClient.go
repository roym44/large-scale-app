package RegServiceClient

import (
	context "context"
	"fmt"

	services "github.com/TAULargeScaleWorkshop/RLAD/services/common"
	service "github.com/TAULargeScaleWorkshop/RLAD/services/reg-service/common"
)

type RegServiceClient struct {
	services.ServiceClientBase[service.RegServiceClient]
}

func NewRegServiceClient(address string) *RegServiceClient {
	return &RegServiceClient{
		ServiceClientBase: services.ServiceClientBase[service.RegServiceClient]{
			Address: address, CreateClient: service.NewRegServiceClient},
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
