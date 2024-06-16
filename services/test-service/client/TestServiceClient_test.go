package TestServiceClient // TODO: Ask Zvi regarding package naming, and import cycle

import (
	"testing"

	// client "github.com/TAULargeScaleWorkshop/RLAD/services/test-service/client/TestServiceClient"
)

func TestHelloWorld(t *testing.T) {
	c := NewTestServiceClient("localhost:50051")
	r, err := c.HelloWorld()
	if err != nil {
		t.Fatalf("could not call HelloWorld: %v", err)
		return
	}
	t.Logf("Response: %v", r)
}

func TestHelloToUser(t *testing.T) {
	c := NewTestServiceClient("localhost:50051")
	r, err := c.HelloToUser("Zvi")
	if err != nil {
		t.Fatalf("could not call HelloToUser: %v", err)
		return
	}
	t.Logf("Response: %v", r)
}