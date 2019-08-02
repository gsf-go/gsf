package rpc

import (
	"github.com/sf-go/gsf/src/gsc/serialization"
	"github.com/sf-go/gsf/src/gsm/peer"
	"reflect"
)

type IRpcInvoke interface {
	Invoke(methodId []byte, peer peer.IPeer, args ...interface{}) []reflect.Value
	Request(methodId []byte, args ...interface{}) []byte
}

type RpcInvoke struct {
	serializable serialization.ISerializable
}

func NewRpcInvoke() *RpcInvoke {
	return &RpcInvoke{
		serializable: serialization.NewSerializable(),
	}
}

//func (rpcInvoke *RpcInvoke) Invoke(methodId []byte, peer peer.IPeer, args ...interface{}) []reflect.Value {
//
//	method := rpcInvoke.register.GetRpcByName(string(methodId))
//	values := make([]reflect.Value, len(args))
//
//	for i, item := range args {
//		v := reflect.ValueOf(item)
//		values[i] = v
//	}
//
//	return method(peer,  ,values)
//}

func (rpcInvoke *RpcInvoke) Request(methodId []byte, args ...interface{}) []byte {
	dataBytes := rpcInvoke.serializable.Serialize(args...)
	methodIdBytes := rpcInvoke.serializable.SerializeSingle(methodId)
	buffer := make([]byte, 0, len(methodIdBytes)+len(dataBytes))
	buffer = append(buffer, methodIdBytes...)
	buffer = append(buffer, dataBytes...)
	return buffer
}
