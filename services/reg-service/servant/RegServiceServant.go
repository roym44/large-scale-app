package RegServiceServant

import (
	"context"
	"fmt"
	"sync"
	"time"

	dht "github.com/TAULargeScaleWorkshop/RLAD/services/reg-service/servant/dht"
	testservicecommon "github.com/TAULargeScaleWorkshop/RLAD/services/test-service/common"
	"github.com/TAULargeScaleWorkshop/RLAD/utils"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	is_first  bool
	chordNode *dht.Chord
	cacheMap  map[string][]Node
	mutex     sync.Mutex
)

type Node struct {
	Address   string
	FailCount int
	Alive     bool
}

func IsFirst() bool {
	utils.Logger.Printf("IsFirst() called, result: %v", is_first)
	return is_first
}

func InitServant(chord_name string) {
	utils.Logger.Printf("RegServiceServant::InitServant() called with %v", chord_name)
	var err error
	// if chord_name == "8502" {
	// }
	// TODO: what happens when a second RegService needs to join? how does he know if he's first without calling NewChord first?
	chordNode, err = dht.NewChord("root", 6666)
	if err != nil {
		utils.Logger.Fatalf("could not create new chord: %v", err)
		return
	}
	utils.Logger.Printf("NewChord returned: %v", chordNode)

	is_first, err = chordNode.IsFirst()
	if err != nil {
		utils.Logger.Fatalf("could not call IsFirst: %v", err)
		return
	}
	utils.Logger.Printf("chordNode.IsFirst() result: %v", is_first)

	// join
	if !is_first {
		utils.Logger.Printf("not first")
		// join already existing "root" with a new chord_name
		chordNode, err = dht.JoinChord(chord_name, "root", 6666)
		if err != nil {
			utils.Logger.Fatalf("could not join chord: %v", err)
			return
		}
		utils.Logger.Printf("JoinChord returned: %v", chordNode)
	} else {
		// we are root, initialize the cache map :)
		utils.Logger.Printf("first!")
		cacheMap = make(map[string][]Node)
	}
}

// Registry API
func Register(service_name string, node_address string) {
	mutex.Lock()
	defer mutex.Unlock()

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
	mutex.Lock()
	defer mutex.Unlock()

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
	mutex.Lock()
	defer mutex.Unlock()

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

// TestServiceClient code (we duplicate some code for the IsAlive grpc connection)
type TestServiceClient struct {
	Address string // we have a specified address, not using the registry
	// client_t == testservicecommon.TestServiceClient
	CreateClient func(grpc.ClientConnInterface) testservicecommon.TestServiceClient
}

func NewTestServiceClient(address string) *TestServiceClient {
	return &TestServiceClient{
		Address:      address,
		CreateClient: testservicecommon.NewTestServiceClient,
	}
}

func (obj *TestServiceClient) Connect() (res testservicecommon.TestServiceClient, closeFunc func(), err error) {
	// Set a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, obj.Address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		var empty testservicecommon.TestServiceClient
		return empty, nil, fmt.Errorf("failed to connect client to %v: %v", obj.Address, err)
	}
	c := obj.CreateClient(conn)
	return c, func() { conn.Close() }, nil
}

func (obj *TestServiceClient) IsAlive() (bool, error) {
	c, closeFunc, err := obj.Connect()
	if err != nil {
		return false, fmt.Errorf("could not connect to server: %v", err)
	}
	defer closeFunc()
	// Call the IsAlive RPC function
	r, err := c.IsAlive(context.Background(), &emptypb.Empty{})
	if err != nil {
		return false, fmt.Errorf("could not call IsAlive: %v", err)
	}
	return r.Value, nil
}

// Internal logic, health checking the nodes
func IsAliveCheck() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		mutex.Lock()
		utils.Logger.Printf("IsAliveCheck: called\n")
		for serviceName, nodes := range cacheMap {
			for i := 0; i < len(nodes); i++ {
				utils.Logger.Printf("IsAliveCheck: Service = %s, Node = %v\n", serviceName, nodes[i])
				c := NewTestServiceClient(nodes[i].Address)
				alive, err := c.IsAlive()
				// we assume that (!alive) iff (err != nil) in IsAlive implementation
				if !alive {
					utils.Logger.Printf("IsAliveCheck: Node %v is not alive: error = %v\n", nodes[i], err)
					nodes[i].FailCount++
					if nodes[i].FailCount >= 2 {
						utils.Logger.Printf("IsAliveCheck: marking node as not alive! %v\n", nodes[i])
						// mark the node to be unregistered later
						nodes[i].Alive = false
					}
				} else {
					nodes[i].FailCount = 0
					nodes[i].Alive = true
				}
			}
			// unregister (manually, not calling the API function) the "dead" nodes
			for i := len(nodes) - 1; i >= 0; i-- {
				if !nodes[i].Alive {
					utils.Logger.Printf("Node %s is not alive, unregistering...\n", nodes[i].Address)
					cacheMap[serviceName] = append(cacheMap[serviceName][:i], cacheMap[serviceName][i+1:]...)
				}
			}
		}
		mutex.Unlock()
	}
}
