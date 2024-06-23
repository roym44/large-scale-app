package TestServiceServant

import (
	"fmt"
	"math/rand"
	"time"

	metaffi "github.com/MetaFFI/lang-plugin-go/api"
	"github.com/MetaFFI/plugin-sdk/compiler/go/IDL"
	"github.com/TAULargeScaleWorkshop/RLAD/utils"
)

// globals
var pythonRuntime *metaffi.MetaFFIRuntime
var crawlerModule *metaffi.MetaFFIModule

// extract_links_from_url(url: str, depth: int) -> list:
var extract_links_from_url func(...interface{}) ([]interface{}, error)
var cacheMap map[string]string

func init() {
	cacheMap = make(map[string]string)

	// Load the Python3.11 runtime
	pythonRuntime = metaffi.NewMetaFFIRuntime("python311")
	err := pythonRuntime.LoadRuntimePlugin()
	if err != nil {
		msg := fmt.Sprintf("Failed to load runtime plugin: %v", err)
		utils.Logger.Fatalf(msg)
		panic(msg)
	}
	// Load the Crawler module
	crawlerModule, err = pythonRuntime.LoadModule("./services/test-service/servant/crawler.py")
	if err != nil {
		msg := fmt.Sprintf("Failed to load ./crawler/crawler.py module: %v", err)
		utils.Logger.Fatalf(msg)
		panic(msg)
	}
	// Load the crawler function
	extract_links_from_url, err = crawlerModule.Load("callable=extract_links_from_url",
		[]IDL.MetaFFIType{IDL.STRING8, IDL.INT64}, // parameters types
		[]IDL.MetaFFIType{IDL.STRING8_ARRAY})      // return type
	if err != nil {
		msg := fmt.Sprintf("Failed to load extract_links_from_url function: %v", err)
		utils.Logger.Fatalf(msg)
		panic(msg)
	}
}

func HelloWorld() string {
	return "Hello World"
}

func HelloToUser(userName string) string {
	return "Hello " + userName
}

func Store(key string, value string) {
	cacheMap[key] = value
}

func Get(key string) (string, error) {
	// returns the value (or "" if not found), and a boolean indicating whether the key was found in the map
	value, ok := cacheMap[key]
	if !ok {
		return value, fmt.Errorf("key not found: %v", key)
	}
	return value, nil
}

func WaitAndRand(seconds int32, sendToClient func(x int32) error) error {
	time.Sleep(time.Duration(seconds) * time.Second)
	return sendToClient(int32(rand.Intn(10)))
}

func IsAlive() bool {
	return true
}

func ExtractLinksFromURL(url string, depth int32) ([]string, error) {
	// Call Python's extract_links_from_url
	res, err := extract_links_from_url(url, int64(depth))
	if err != nil {
		return nil, err
	}
	return res[0].([]string), nil
}
