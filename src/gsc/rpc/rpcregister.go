package rpc

import (
	"github.com/sf-go/gsf/src/gsf/peer"
	"reflect"
	"sync"
)

//go:generate ../../../gen/Singleton.exe -struct=rpcRegister -out=./rpcregister.go

var RpcRegisterInstance = NewRpcRegister()

type RpcRegister struct {
	cache *sync.Map
}

func NewRpcRegister() *RpcRegister {
	return &RpcRegister{
		cache: new(sync.Map),
	}
}

func (rpcRegister *RpcRegister) Add(name string,
	method func(peer peer.IPeer, values []reflect.Value) []reflect.Value) {

	rpcRegister.cache.Store(name, method)
}

func (rpcRegister *RpcRegister) Remove(name string) {
	rpcRegister.cache.Delete(name)
}

func (rpcRegister *RpcRegister) GetRpcByName(name string) func(peer peer.IPeer, values []reflect.Value) []reflect.Value {
	value, ok := rpcRegister.cache.Load(name)
	if ok {
		return value.(func(peer peer.IPeer, values []reflect.Value) []reflect.Value)
	}
	return nil
}
