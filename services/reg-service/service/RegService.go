package TestService

import (
	"context"

	services "github.com/TAULargeScaleWorkshop/RLAD/services/common"       // import common as services
	// . "github.com/TAULargeScaleWorkshop/RLAD/services/test-service/common" // from test-service/common import *
	// RegistryServiceServant "github.com/TAULargeScaleWorkshop/RLAD/services/test-service/servant"
	// . "github.com/TAULargeScaleWorkshop/RLAD/utils" // from utils import *

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type testServiceImplementation struct {
	UnimplementedTestServiceServer
}

func Start(configData []byte) error {
	bindgRPCToService := func(s grpc.ServiceRegistrar) {
		RegTestServiceServer(s, &regServiceImplementation{})
	}
	services.Start("RegistryService", 8502, bindgRPCToService)
	return nil
}
