package TestService

import (
	"context"

	services "github.com/TAULargeScaleWorkshop/RLAD/services/common"       // import common as services
	. "github.com/TAULargeScaleWorkshop/RLAD/services/test-service/common" // from test-service/common import *
	. "github.com/TAULargeScaleWorkshop/RLAD/utils"                        // from utils import *
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
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

func (obj *testServiceImplementation) HelloWorld(_ context.Context, _ *emptypb.Empty) (res *wrapperspb.StringValue, err error) {
	Logger.Printf("HelloWorld called")
	return wrapperspb.String("Hello World"), nil
}
