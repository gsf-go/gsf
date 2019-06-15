package rpc

import (
	"github.com/gsf/gsf/src/gsf/peer"
	"reflect"
	"sync"
)

//go:generate ../../../gen/Singleton.exe -struct=rpcRegister -out=./rpcregister.go

var rpcRegisterInstance *rpcRegister
var rpcRegisterOnce sync.Once

func GetRpcRegisterInstance() *rpcRegister {
	rpcRegisterOnce.Do(func() {
		rpcRegisterInstance = NewRpcRegister()
	})
	return rpcRegisterInstance
}

type rpcRegister struct {
	cache *sync.Map
}

func NewRpcRegister() *rpcRegister {
	return &rpcRegister{
		cache: new(sync.Map),
	}
}

func (rpcRegister *rpcRegister) Add(name string,
	method func(peer peer.IPeer, values []reflect.Value) []reflect.Value) {

	rpcRegister.cache.Store(name, method)
}

func (rpcRegister *rpcRegister) Remove(name string) {
	rpcRegister.cache.Delete(name)
}

func (rpcRegister *rpcRegister) GetRpcByName(name string) func(peer peer.IPeer, values []reflect.Value) []reflect.Value {
	value, ok := rpcRegister.cache.Load(name)
	if ok {
		return value.(func(peer peer.IPeer, values []reflect.Value) []reflect.Value)
	}
	return nil
}
