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

func TestStoreGet(t *testing.T) {
	c := NewTestServiceClient("localhost:50051")
	err := c.Store("key1","value1")
	if err != nil {
		t.Fatalf("could not call Store: %v", err)
		return
	}
	r, err := c.Get("key1")
	if err != nil {
		t.Fatalf("could not call Get: %v", err)
		return
	}
	if(r!="value1"){
		t.Fatalf("wrong value: received %s, expected value1",r)
		return
	}
	t.Logf("Response: %v", r)
}

func TestWaitAndRand(t *testing.T) {
	c := NewTestServiceClient("localhost:50051")
	resPromise, err := c.WaitAndRand(3)
	if err != nil {
		t.Fatalf("Calling WaitAndRand failed: %v", err)
		return
	}
	res, err := resPromise()
	if err != nil {
		t.Fatalf("WaitAndRand failed: %v", err)
		return
	}
	t.Logf("Returned random number: %v\n", res)
}