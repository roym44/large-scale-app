package CacheService

import (
	"context"
	"fmt"
	"net"

	"github.com/TAULargeScaleWorkshop/RLAD/config"
	. "github.com/TAULargeScaleWorkshop/RLAD/services/cache-service/common"
	. "github.com/TAULargeScaleWorkshop/RLAD/utils"
	"github.com/TAULargeScaleWorkshop/RLAD/services/reg-service/servant/dht"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gopkg.in/yaml.v2"
)

type CacheService struct {
	UnimplementedCacheServiceServer
	chord dht.Chord
}

func NewCacheService(nodeName string, port int32) (*CacheService, error) {
    chordNode, err := dht.NewChord(nodeName, port)
    if err != nil {
        return nil, err
    }

    return &CacheService{chord: chordNode}, nil
}

func (cs *CacheService) Set(_ context.Context, params *SetKeyValueReq) (_ *emptypb.Empty, err error) {
	Logger.Printf("Set called with key: %s, value: %s", params.Key, params.Value)
	cs.chord.Set(params.Key, params.Value)
	return &emptypb.Empty{}, nil
}

func (cs *CacheService) Get(_ context.Context, k *GetKeyReq) (*GetValueReq, err error) {
	Logger.Printf("Get called with key: %s", k.Key)
	value, err := cs.chord.Get(k.Key)
	return &GetValueReq{Value: value}, nil
}

func (cs *CacheService) Delete(_ context.Context, k *GetKeyReq) (_ *emptypb.Empty, err error) {
	Logger.Printf("Delete called with key: %s", k.Key)
	cs.chord.Delete(k.Key)
	return &emptypb.Empty{}, nil
}

func (cs *CacheService) IsAlive(_ context.Context, _ *emptypb.Empty) (*wrapperspb.BoolValue, err error) {
	Logger.Printf("IsAlive called ")
	IsFirst, err := cs.chord.IsFirst()
	return wrapperspb.Bool(isFirst), nil
}

