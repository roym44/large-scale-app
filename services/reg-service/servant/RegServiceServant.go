package RegServiceServant

import (
	"context"
	"fmt"
	"strings"
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
	// address : NodeStatus
	cacheMap map[string]*NodeStatus
	mutex    sync.Mutex
)

type NodeStatus struct {
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
		cacheMap = make(map[string]*NodeStatus)
	}
}

func encodeStrings(lst []string) string {
	return strings.Join(lst, ",")
}

func decodeStrings(enc string) []string {
	return strings.Split(enc, ",")
}

// Registry API
func Register(service_name string, node_address string) {
	mutex.Lock()
	defer mutex.Unlock()

	// get the current list
	utils.Logger.Printf("chordNode.Get before")
	enc, err := chordNode.Get(service_name)
	utils.Logger.Printf("chordNode.Get after")
	if err != nil {
		utils.Logger.Fatalf("chordNode.Get failed with error: %v", err)
	}
	utils.Logger.Printf("decodeStrings before")
	lst := decodeStrings(enc)
	if len(lst) > 0 {
		for _, address := range lst {
			if address == node_address {
				utils.Logger.Printf("Address %s already exists for service %s\n", node_address, service_name)
				return
			}
		}
	}

	// add to list and set back to chord
	lst = append(lst, node_address)
	updated_enc := encodeStrings(lst)
	utils.Logger.Printf("chordNode.Set before")
	err = chordNode.Set(service_name, updated_enc)
	utils.Logger.Printf("chordNode.Set after")
	if err != nil {
		utils.Logger.Fatalf("chordNode.Set failed with error: %v", err)
	}
	utils.Logger.Printf("Address %s added for service %s\n", node_address, service_name)
}

func unregisterFromChord(service_name string, node_address string) {
	// get the current list
	enc, err := chordNode.Get(service_name)
	if err != nil {
		utils.Logger.Fatalf("chordNode.Get failed with error: %v", err)
	}
	lst := decodeStrings(enc)
	if len(lst) > 0 {
		for i, address := range lst {
			if address == node_address {
				// remove from list and set back to chord
				lst = append(lst[:i], lst[i+1:]...)
				utils.Logger.Printf("Address %s removed for service %s\n", node_address, service_name)
				updated_enc := encodeStrings(lst)
				err = chordNode.Set(service_name, updated_enc)
				if err != nil {
					utils.Logger.Fatalf("chordNode.Set failed with error: %v", err)
				}
				return
			}
		}
		utils.Logger.Printf("Address %s not found for service %s\n", node_address, service_name)
	} else {
		utils.Logger.Printf("Service %s not found\n", service_name)
	}
}

func Unregister(service_name string, node_address string) {
	mutex.Lock()
	defer mutex.Unlock()
	unregisterFromChord(service_name, node_address)
}

func Discover(service_name string) ([]string, error) {
	mutex.Lock()
	defer mutex.Unlock()

	// get the current list
	enc, err := chordNode.Get(service_name)
	if err != nil {
		utils.Logger.Fatalf("chordNode.Get failed with error: %v", err)
	}
	lst := decodeStrings(enc)
	if len(lst) == 0 {
		return lst, fmt.Errorf("service not found: %v", service_name)
	}
	return lst, nil
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

		// get all the services
		services, err := chordNode.GetAllKeys()
		if err != nil {
			utils.Logger.Fatalf("chordNode.GetAllKeys failed with error: %v", err)
		}

		for _, serviceName := range services {
			// get the current list
			enc, err := chordNode.Get(serviceName)
			if err != nil {
				utils.Logger.Fatalf("chordNode.Get failed with error: %v", err)
			}
			addresses := decodeStrings(enc)

			for _, address := range addresses {
				utils.Logger.Printf("IsAliveCheck: Service = %s, Node = %v\n", serviceName, address)
				c := NewTestServiceClient(address)
				alive, err := c.IsAlive()

				// create node status if doesn't exist
				_, ok := cacheMap[address]
				if !ok {
					cacheMap[address] = &NodeStatus{FailCount: 0, Alive: true}
				}

				// we assume that (!alive) iff (err != nil) in IsAlive implementation
				if !alive {
					utils.Logger.Printf("IsAliveCheck: Node %v is not alive: error = %v\n", address, err)
					cacheMap[address].FailCount++
					if cacheMap[address].FailCount >= 2 {
						utils.Logger.Printf("IsAliveCheck: marking node as not alive! %v\n", address)
						// mark the node to be unregistered later
						cacheMap[address].Alive = false
					}
				} else {
					cacheMap[address].FailCount = 0
					cacheMap[address].Alive = true
				}
			}

			// unregister (manually, not calling the API function) the "dead" nodes
			for i := len(addresses) - 1; i >= 0; i-- {
				if !cacheMap[addresses[i]].Alive {
					utils.Logger.Printf("Node %s is not alive, unregistering...\n", addresses[i])
					// remove from cache
					delete(cacheMap, addresses[i])

					// unregister
					unregisterFromChord(serviceName, addresses[i])
				}
			}
		}
		mutex.Unlock()
	}
}
