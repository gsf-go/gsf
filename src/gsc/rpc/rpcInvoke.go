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
	register     *RpcRegister
}

func NewRpcInvoke(register *RpcRegister) *RpcInvoke {
	return &RpcInvoke{
		serializable: serialization.NewSerializable(),
		register:     register,
	}
}

func (rpcInvoke *RpcInvoke) Invoke(methodId []byte, peer peer.IPeer, args ...interface{}) []reflect.Value {

	method := rpcInvoke.register.GetRpcByName(string(methodId))
	values := make([]reflect.Value, len(args))

	for i, item := range args {
		v := reflect.ValueOf(item)
		values[i] = v
	}

	return method(peer, values)
}

func (rpcInvoke *RpcInvoke) Request(methodId []byte, args ...interface{}) []byte {
	slice := make([]interface{}, 0, len(args)+1)
	slice = append(slice, methodId)
	slice = append(slice, args...)
	return rpcInvoke.serializable.Serialize(args...)
}
