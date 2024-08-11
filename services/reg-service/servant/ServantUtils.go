package RegServiceServant

import (
	"fmt"
	"errors"

	"github.com/pebbe/zmq4"

	"github.com/TAULargeScaleWorkshop/RLAD/utils"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"google.golang.org/protobuf/proto"

)

// -------------------- TestServiceMQ --------------------
func (obj *TestServiceClient) ConnectMQ() (socket *zmq4.Socket, err error) {
	socket, err = zmq4.NewSocket(zmq4.REQ)
	utils.Logger.Printf("ConnectMQ(): created NewSocket %v", socket)
	if err != nil {
		utils.Logger.Fatalf("Failed to create a new zmq socket: %v", err)
	}
	utils.Logger.Printf("ConnectMQ(): calling Connect for address: %s", obj.AddressMQ)
	err = socket.Connect(obj.AddressMQ)
	utils.Logger.Printf("after Connect")
	if err != nil {
		utils.Logger.Printf("Failed to connect a zmq socket: %v", err)
	}
	return socket, err
}

func (obj *TestServiceClient) NewMarshaledCallParameter(method string, proto_data proto.Message) ([]byte, error) {
	var msg []byte

	// handle data
	data, err := proto.Marshal(proto_data)
	if err != nil {
		utils.Logger.Printf("NewMarshaledCallParameter(): Marshal(proto_data) failed: %v\n", err)
	}

	// handle call params
	callParams := &CallParameters{Method: method, Data: data}
	msg, err = proto.Marshal(callParams)
	if err != nil {
		utils.Logger.Printf("NewMarshaledCallParameter(): Marshal(callParams) failed: %v\n", err)
	}
	return msg, err
}

func (obj *TestServiceClient) IsAliveAsync() (func() (bool, error), error) {
	utils.Logger.Printf("before ConnectMQ")
	mqsocket, err := obj.ConnectMQ()
	utils.Logger.Printf("after ConnectMQ")
	if err != nil {
		return nil, fmt.Errorf("IsAliveAsync(): ConnectMQ failed: %v\n", err)
	}

	// packing
	utils.Logger.Printf("before obj.NewMarshaledCallParameter")
	msg, err := obj.NewMarshaledCallParameter("IsAlive", &emptypb.Empty{})
	utils.Logger.Printf("after obj.NewMarshaledCallParameter")
	if err != nil {
		return nil, fmt.Errorf("IsAliveAsync(): NewMarshaledCallParameter failed: %v\n", err)
	}

	_, err = mqsocket.SendBytes(msg, 0)
	if err != nil {
		return nil, fmt.Errorf("IsAliveAsync(): SendBytes failed: %v\n", err)
	}

	// return function (future pattern)
	ret := func() (bool, error) {
		defer mqsocket.Close()
		rv, err := mqsocket.RecvBytes(0)
		if err != nil {
			return false, fmt.Errorf("IsAliveAsync(): RecvBytes failed: %v\n", err)
		}

		returnValue := &ReturnValue{}
		// handle return value
		err = proto.Unmarshal(rv, returnValue)
		if err != nil {
			return false, fmt.Errorf("IsAliveAsync(): Unmarshal(rv) failed: %v\n", err)
		}

		// handle data
		r := &wrapperspb.BoolValue{}
		err = proto.Unmarshal(returnValue.Data, r)
		if err != nil {
			return false, fmt.Errorf("IsAliveAsync(): Unmarshal(returnValue.Data) failed: %v\n", err)
		}

		// error
		if returnValue.Error != "" {
			err = errors.New(returnValue.Error)
		}
		return r.Value, err
	}

	return ret, nil
}
