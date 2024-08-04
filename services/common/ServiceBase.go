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

// regAddresses - the registry service addresses to connect to (for all the services who aren't the registry)
func Start(serviceName string, grpcListenPort int, regAddresses []string, bindgRPCToService func(s grpc.ServiceRegistrar)) (err error) {
	Logger.Printf("Start(%s, %d)\n", serviceName, grpcListenPort)
	// start the service
	listeningAddress, grpcServer, startListening, err := startgRPC(grpcListenPort)
	if err != nil {
		return err
	}
	bindgRPCToService(grpcServer)

	// make sure it registers to the registry service
	unregister := registerAddress(serviceName, regAddresses, listeningAddress)
	defer unregister()

	Logger.Printf("%s starts listening on %s\n", serviceName, listeningAddress)
	startListening()
	return nil
}

func registerAddress(serviceName string, regAddresses []string, listeningAddress string) (unregister func()) {
	regClient := RegServiceClient.NewRegServiceClient(regAddresses)
	err := regClient.Register(serviceName, listeningAddress)
	if err != nil {
		Logger.Fatalf("Failed to register to registry service: %v", err)
	}
	return func() {
		regClient.Unregister(serviceName, listeningAddress)
	}
}
