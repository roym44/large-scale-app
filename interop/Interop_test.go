package Interop

import (
	"os"
	"testing"
)

// entry point for the test suite
func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestExtractLinksFromURL(t *testing.T) {
	// Load Python
	err := LoadPython()
	if err != nil {
		t.Fatalf("LoadPython failed with error: %v", err)
	}

	// Extract links from www.microsoft.com
	url := "https://www.microsoft.com"

	// call the function with "depth" 1.
	links, err := ExtractLinksFromURL(url, 1)
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
