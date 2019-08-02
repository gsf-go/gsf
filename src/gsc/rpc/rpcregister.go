package rpc

import (
	"github.com/sf-go/gsf/src/gsm/peer"
	"reflect"
	"sync"
)

//go:generate ../../../gen/Singleton.exe -struct=rpcRegister -out=./rpcregister.go

type RpcRegister struct {
	cache *sync.Map
}

func NewRpcRegister() *RpcRegister {
	return &RpcRegister{
		cache: new(sync.Map),
	}
}

func (rpcRegister *RpcRegister) AddRequest(id []byte,
	method func(peer peer.IPeer, values []reflect.Value) []reflect.Value) {
	rpcRegister.cache.Store(string(id), method)
}

func (rpcRegister *RpcRegister) AddResponse(id []byte,
	method func(peer peer.IPeer, values []reflect.Value) []reflect.Value) {

	tmp := make([]byte, len(id)+1)
	tmp = append(tmp, 1)
	tmp = append(tmp, id...)

	rpcRegister.cache.Store(string(tmp), method)
}

func (rpcRegister *RpcRegister) RemoveResponse(id []byte) {
	tmp := make([]byte, len(id)+1)
	tmp = append(tmp, 1)
	tmp = append(tmp, id...)
	rpcRegister.cache.Delete(string(tmp))
}

func (rpcRegister *RpcRegister) GetRpcByName(name string) func(peer peer.IPeer, values []reflect.Value) []reflect.Value {
	value, ok := rpcRegister.cache.Load(name)
	if ok {
		return value.(func(peer peer.IPeer, values []reflect.Value) []reflect.Value)
	}
	return nil
}
