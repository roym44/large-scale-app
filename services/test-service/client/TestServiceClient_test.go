package TestServiceClient

import (
	"testing"
)
// TODO: update tests to pass addresses []string from config
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

func TestStoreAndGet(t *testing.T) {
	c := NewTestServiceClient("localhost:50051")
	err := c.Store("key1", "value1")
	if err != nil {
		t.Fatalf("could not call Store: %v", err)
		return
	}
	r, err := c.Get("key1")
	if err != nil {
		t.Fatalf("could not call Get: %v", err)
		return
	}
	if r != "value1" {
		t.Fatalf("wrong value: received %s, expected value1", r)
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

func TestIsAlive(t *testing.T) {
	c := NewTestServiceClient("localhost:50051")
	r, err := c.IsAlive()
	if err != nil {
		t.Fatalf("could not call IsAlive: %v", err)
		return
	}
	if !r {
		t.Fatalf("IsAlive returned false")
		return
	}
	t.Logf("Response: %v", r)
}

func TestExtractLinksFromURL(t *testing.T) {
	c := NewTestServiceClient("localhost:50051")

	url := "https://www.microsoft.com"
	links, err := c.ExtractLinksFromURL(url, 1)
	if err != nil {
		t.Fatalf("ExtractLinksFromURL failed with error: %v", err)
	}

	// make sure you got some links
	if len(links) == 0 {
		t.Fatalf("ExtractLinksFromURL returned no links")
	}

	// print the links
	t.Logf("links: %v\n", links)
}
