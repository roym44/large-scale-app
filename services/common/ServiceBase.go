package common

import (
	"fmt"
	"net"

	RegServiceClient "github.com/TAULargeScaleWorkshop/RLAD/services/reg-service/client"
	. "github.com/TAULargeScaleWorkshop/RLAD/utils"
	"google.golang.org/grpc"
)

func startgRPC(listenPort int) (listeningAddress string, grpcServer *grpc.Server, startListening func(), err error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", listenPort))
	if err != nil {
		Logger.Fatalf("failed to listen: %v", err)
		return "", nil, nil, err
	}
	listeningAddress = lis.Addr().String()
	grpcServer = grpc.NewServer()
	startListening = func() {
		if err := grpcServer.Serve(lis); err != nil {
			Logger.Fatalf("failed to serve: %v", err)
		}
	}
	return
}

func Start(serviceName string, grpcListenPort int, addressPtr *string, bindgRPCToService func(s grpc.ServiceRegistrar)) (err error) {
	listeningAddress, grpcServer, startListening, err := startgRPC(grpcListenPort)
	if err != nil {
		return err
	}
	bindgRPCToService(grpcServer)
	*addressPtr = listeningAddress
	startListening()
	return nil
}

func RegisterAddress(serviceName string, regAddresses []string, listeningAddress string) (unregister func()) {
	regClient := RegServiceClient.NewRegServiceClient(regAddresses)
	err := regClient.Register(serviceName, listeningAddress)
	if err != nil {
		Logger.Fatalf("Failed to register to registry service: %v", err)
	}
	return func() {
		regClient.Unregister(serviceName, listeningAddress)
	}
}
