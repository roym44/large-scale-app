package TestService

import (
	context "context"
	"testing"

	common "github.com/TAULargeScaleWorkshop/RLAD/services/test-service/common"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestHelloWorld(t *testing.T) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		t.Fatalf("did not connect: %v", err)
		return
	}
	defer conn.Close()
	c := common.NewTestServiceClient(conn)

	// Call the HelloWorld RPC function
	r, err := c.HelloWorld(context.Background(), &emptypb.Empty{})
	if err != nil {
		t.Fatalf("could not call HelloWorld: %v", err)
		return
	}
	t.Logf("Response: %v", r.GetValue())
}
