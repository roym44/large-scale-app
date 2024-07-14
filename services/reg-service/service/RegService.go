package RegService

import (
	"context"

	services "github.com/TAULargeScaleWorkshop/RLAD/services/common"
	. "github.com/TAULargeScaleWorkshop/RLAD/services/reg-service/common"
	RegServiceServant "github.com/TAULargeScaleWorkshop/RLAD/services/reg-service/servant"
	. "github.com/TAULargeScaleWorkshop/RLAD/utils"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type regServiceImplementation struct {
	UnimplementedRegServiceServer
}

func Start(configData []byte) error {
	// TODO: use listen_port from config
	// TODO: try the first one (8502) then 8503, 8504...
	bindgRPCToService := func(s grpc.ServiceRegistrar) {
		RegisterRegServiceServer(s, &regServiceImplementation{})
	}
	services.Start("RegistryService", 8502, bindgRPCToService)
	return nil
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
