package TestServiceServant

import (
	"math/rand"
	"time"
)

var cacheMap map[string]string

func init() {
	cacheMap = make(map[string]string)
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

func Get(key string) (string, bool) {
	// returns the value (or "" if not found), and a boolean indicating whether the key was found in the map
	value, ok := cacheMap[key]
	return value, ok
}

func WaitAndRand(seconds int32, sendToClient func(x int32) error) error {
	time.Sleep(time.Duration(seconds) * time.Second)
	return sendToClient(int32(rand.Intn(10)))
}

func IsAlive() bool {
	return true
}
