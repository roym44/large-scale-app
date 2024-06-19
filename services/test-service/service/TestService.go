package TestService

import (
	"context"

	services "github.com/TAULargeScaleWorkshop/RLAD/services/common"       // import common as services
	. "github.com/TAULargeScaleWorkshop/RLAD/services/test-service/common" // from test-service/common import *
	TestServiceServant "github.com/TAULargeScaleWorkshop/RLAD/services/test-service/servant"
	. "github.com/TAULargeScaleWorkshop/RLAD/utils" // from utils import *

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
	return nil // TODO: Ask Zvi, should we return nil, or remove the error in signiture?
}

func (obj *testServiceImplementation) HelloWorld(ctxt context.Context, _ *emptypb.Empty) (res *wrapperspb.StringValue, err error) {
	Logger.Printf("HelloWorld called")
	return wrapperspb.String(TestServiceServant.HelloWorld()), nil
}

func (obj *testServiceImplementation) HelloToUser(_ context.Context, userName *wrapperspb.StringValue) (res *wrapperspb.StringValue, err error) {
	Logger.Printf("HelloToUser called")
	return wrapperspb.String(TestServiceServant.HelloToUser(userName.Value)), nil
}

func (obj *testServiceImplementation) Store(ctx context.Context, req *StoreKeyValue) (_ *emptypb.Empty, err error) {
	Logger.Printf("Store called with key: %s, value: %s",req.Key,req.Value)
	TestServiceServant.Store(req.Key,req.Value)
	return &emptypb.Empty{}, nil
}

func (obj *testServiceImplementation) Get(ctx context.Context, key *wrapperspb.StringValue) (res *wrapperspb.StringValue, err error) {
	Logger.Printf("Get called with key: %s",key.Value)
	value,ok:=TestServiceServant.Get(key.Value)
	if(ok!=nil){
		return nil, fmt.Errorf("key not found: %v", key.Value)
	}
	return wrapperspb.String(value), nil
}

func (obj *testServiceImplementation) WaitAndRand(seconds
	*wrapperspb.Int32Value, streamRet TestService_WaitAndRandServer) error {
	Logger.Printf("WaitAndRand called")
	streamClient := func(x int32) error {
		return streamRet.Send(wrapperspb.Int32(x))
	}
	return TestServiceServant.WaitAndRand(seconds.Value, streamClient)
	}