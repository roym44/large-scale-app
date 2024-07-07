package RegServiceServant

import (
	"fmt"
	"sync"
	"math/rand"
	"time"

	metaffi "github.com/MetaFFI/lang-plugin-go/api"
	"github.com/MetaFFI/plugin-sdk/compiler/go/IDL"
	"github.com/TAULargeScaleWorkshop/RLAD/utils"
)

var (
    cacheMap map[string][]string
    fails    map[string]int
    mu       sync.Mutex
)

func init() {
	cacheMap = make(map[string][]string)
	fails = make(map[string]int)

}

func Register(service_name string, node_address string) {
	if cacheMap[service_name] != nil {
		for _, addr := range cacheMap[service_name] {
			if addr == node_address {
				utils.Logger.Printf("Address %s already exists for service %s\n", node_address, service_name)
				return
			}
		}
	}
	cacheMap[service_name] = append(cacheMap[service_name], node_address)
	utils.Logger.Printf("Address %s added for service %s\n", node_address, service_name)
}

func Unregister(service_name string, node_address string) {
	if cacheMap[service_name] != nil {
		for i, addr := range cacheMap[service_name] {
			if addr == node_address {
				cacheMap[service_name] = append(cacheMap[service_name][:i], cacheMap[service_name][i+1:]...)
				utils.Logger.Printf("Address %s removed for service %s\n", node_address, service_name)
				return
			}
		}
		utils.Logger.Printf("Address %s not found for service %s\n", node_address, service_name)
	} else {
		utils.Logger.Printf("Service %s not found\n", service_name)
	}
}

func Discover(service_name string) ([]string, error) {
	value, ok := cacheMap[service_name]
	if !ok {
		return value, fmt.Errorf("service not found: %v", service_name)
	}
	return value, nil
}
//TODO:implement IsAlive
func IsAliveCheck() {
    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()
    for range ticker.C {
        for serviceName, addresses := range cacheMap {
            for _, address := range addresses {
				mu.Lock()
                alive, err := IsAlive(address)
                if err != nil {
                    utils.Logger.Printf("IsAlive check failed for address %s: %v\n", address, err)
					mu.Unlock()
					continue
                } else if !alive {
                    count := fails[address]
                    count++
                    if count >= 2 {
                        Unregister(serviceName, address)
                        delete(fails, address)
                        utils.Logger.Printf("Node %s unregistered from service\n", address, serviceName)
                    } else {
                        fails[address] = count
                    }
                } else {
                    delete(fails, address)
                }
				mu.Unlock()
            }
        }
    }
}