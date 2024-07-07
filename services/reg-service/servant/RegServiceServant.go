package RegServiceServant

import (
	"fmt"
	"time"

	. "github.com/TAULargeScaleWorkshop/RLAD/services/test-service/client"
	"github.com/TAULargeScaleWorkshop/RLAD/utils"
)

var (
	cacheMap map[string][]Node
)

type Node struct {
	Address   string
	FailCount int
	Alive     bool
}

func init() {
	cacheMap = make(map[string][]Node)

}

func Register(service_name string, node_address string) {
	if cacheMap[service_name] != nil {
		for _, node := range cacheMap[service_name] {
			if node.Address == node_address {
				utils.Logger.Printf("Address %s already exists for service %s\n", node_address, service_name)
				return
			}
		}
	}
	cacheMap[service_name] = append(cacheMap[service_name], Node{Address: node_address, FailCount: 0, Alive: true})
	utils.Logger.Printf("Address %s added for service %s\n", node_address, service_name)
}

func Unregister(service_name string, node_address string) {
	if cacheMap[service_name] != nil {
		for i, node := range cacheMap[service_name] {
			if node.Address == node_address {
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
	nodes, ok := cacheMap[service_name]
	addresses := []string{}
	if !ok {
		return addresses, fmt.Errorf("service not found: %v", service_name)
	}
	for _, node := range nodes {
		addresses = append(addresses, node.Address)
	}
	return addresses, nil
}

func IsAliveCheck() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		for serviceName, nodes := range cacheMap {
			for i, node := range nodes {
				c := NewTestServiceClient(node.Address)
				alive, err := c.IsAlive()
				// we assume that (!alive) iff (err != nil) in IsAlive implementation
				if !alive {
					utils.Logger.Printf("Service is not alive, address %s: %v\n", node.Address, err)
					nodes[i].FailCount++
					if node.FailCount >= 2 {
						// mark the node to be unregistered later
						nodes[i].Alive = false
					}
				} else {
					nodes[i].FailCount = 0
					nodes[i].Alive = true
				}
			}
			// unregister the "dead" nodes
			for i := len(nodes); i >= 0; i-- {
				if !nodes[i].Alive {
					Unregister(serviceName, nodes[i].Address)
				}
			}
		}
	}
}
