package TestService

import (
	services "github.com/TAULargeScaleWorkshop/RLAD/services/common"       // import common as services
	. "github.com/TAULargeScaleWorkshop/RLAD/services/test-service/common" // from test-service/common import *
	"google.golang.org/grpc"
)

type testServiceImplementation struct {
	UnimplementedTestServiceServer
}

func Start(configData []byte) error {
	bindgRPCToService := func(s grpc.ServiceRegistrar) {
		RegisterTestServiceServer(s, &testServiceImplementation{})
	}
	services.Start("TestService", 50051, bindgRPCToService)
}
