package rpc

import (
	"github.com/sf-go/gsf/src/gsc/serialization"
	"github.com/sf-go/gsf/src/gsf/peer"
	"reflect"
)

type IRpcInvoke interface {
	Invoke(methodId string, peer peer.IPeer, args ...interface{}) []reflect.Value
	Request(methodId string, args ...interface{}) []byte
}

type RpcInvoke struct {
	serializable serialization.ISerializable
}

func NewRpcInvoke() *RpcInvoke {
	return &RpcInvoke{
		serializable: serialization.NewSerializable(),
	}
}

func (rpcInvoke *RpcInvoke) Invoke(methodId string, peer peer.IPeer, args ...interface{}) []reflect.Value {

	method := RpcRegisterInstance.GetRpcByName(methodId)
	values := make([]reflect.Value, len(args))

	for i, item := range args {
		v := reflect.ValueOf(item)
		values[i] = v
	}

	return method(peer, values)
}

func (rpcInvoke *RpcInvoke) Request(methodId string, args ...interface{}) []byte {
	slice := make([]interface{}, 0, len(args)+1)
	slice = append(slice, methodId)
	slice = append(slice, args...)
	return rpcInvoke.serializable.Serialize(slice...)
}
