package TestServiceClient

import (
	context "context"
	"fmt"

	services "github.com/TAULargeScaleWorkshop/RLAD/services/common"
	service "github.com/TAULargeScaleWorkshop/RLAD/services/test-service/common"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type TestServiceClient struct {
	services.ServiceClientBase[service.TestServiceClient]
}

func NewTestServiceClient(address string) *TestServiceClient {
	return &TestServiceClient{
		ServiceClientBase: services.ServiceClientBase[service.TestServiceClient]{Address: address,
			CreateClient: service.NewTestServiceClient},
	}
}
func (obj *TestServiceClient) HelloWorld() (string, error) {
	c, closeFunc, err := obj.Connect()
	if err != nil {
		return "", fmt.Errorf("could not connect to server: %v", err)
	}
	defer closeFunc()
	// Call the HelloWorld RPC function
	r, err := c.HelloWorld(context.Background(), &emptypb.Empty{})
	if err != nil {
		return "", fmt.Errorf("could not call HelloWorld: %v", err)
	}
	return r.Value, nil
}

func (obj *TestServiceClient) HelloToUser(userName string) (string, error) {
	c, closeFunc, err := obj.Connect()
	if err != nil {
		return "", fmt.Errorf("could not connect to server: %v", err)
	}
	defer closeFunc()
	// Call the HelloWorld RPC function
	r, err := c.HelloToUser(context.Background(), wrapperspb.String(userName))
	if err != nil {
		return "", fmt.Errorf("could not call HelloToUser: %v", err)
	}
	return r.Value, nil
}
func (obj *TestServiceClient) Store(key,value string) error {
	c, closeFunc, err := obj.Connect()
	if err != nil {
		return fmt.Errorf("could not connect to server: %v", err)
	}
	defer closeFunc()
	// Call the Store RPC function
	_, err := c.Store(context.Background(), &service.StoreKeyValue{Key:key, Value:value})
	if err != nil {
		return fmt.Errorf("could not call Store: %v", err)
	}
	return nil
}

func (obj *TestServiceClient) Get(key string) (string, error) {
	c, closeFunc, err := obj.Connect()
	if err != nil {
		return "",fmt.Errorf("could not connect to server: %v", err)
	}
	defer closeFunc()
	// Call the Get RPC function
	r, err := c.Get(context.Background(), wrapperspb.String(key))
	if err != nil {
		return "",fmt.Errorf("could not call Store: %v", err)
	}
	return r.Value, nil
}
