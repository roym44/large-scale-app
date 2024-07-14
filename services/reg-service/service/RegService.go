package RegService

import (
	"context"
	"fmt"

	services "github.com/TAULargeScaleWorkshop/RLAD/services/common"
	. "github.com/TAULargeScaleWorkshop/RLAD/services/reg-service/common"
	RegServiceServant "github.com/TAULargeScaleWorkshop/RLAD/services/reg-service/servant"
	. "github.com/TAULargeScaleWorkshop/RLAD/utils"
	"github.com/TAULargeScaleWorkshop/RLAD/config"

	"gopkg.in/yaml.v2"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type regServiceImplementation struct {
	UnimplementedRegServiceServer
}

func Start(configData []byte) error {
	var config config.RegConfig
	err := yaml.Unmarshal(configData, &config) // parses YAML
	if err != nil {
		Logger.Fatalf("error unmarshaling data: %v", err)
	}

	bindgRPCToService := func(s grpc.ServiceRegistrar) {
		RegisterRegServiceServer(s, &regServiceImplementation{})
	}
	var listening_address string
	for port := config.ListenPort; port < config.ListenPort + 10; port++ {
		services.Start("RegistryService", port, &listening_address, bindgRPCToService)
		// will return only if failed to connect
	}
	Logger.Fatalf("Failed to connect to all the registry servers")
	return fmt.Errorf("Failed to connect to all the registry servers")
}

func (obj *regServiceImplementation) Register(_ context.Context, params *UpdateRegistryParameters) (_ *emptypb.Empty, err error) {
	Logger.Printf("Register called with service: %s, address: %s", params.ServiceName, params.NodeAddr)
	RegServiceServant.Register(params.ServiceName, params.NodeAddr)
	return &emptypb.Empty{}, nil
}

func (obj *regServiceImplementation) Unregister(_ context.Context, params *UpdateRegistryParameters) (_ *emptypb.Empty, err error) {
	Logger.Printf("Unregister called with service: %s, address: %s", params.ServiceName, params.NodeAddr)
	RegServiceServant.Unregister(params.ServiceName, params.NodeAddr)
	return &emptypb.Empty{}, nil
}

func (obj *regServiceImplementation) Discover(_ context.Context, service_name *wrapperspb.StringValue) (*DiscoveredServers, error) {
	Logger.Printf("Discover called with service: %s", service_name.Value)
	value, err := RegServiceServant.Discover(service_name.Value)
	return &DiscoveredServers{Addresses: value}, err
}
