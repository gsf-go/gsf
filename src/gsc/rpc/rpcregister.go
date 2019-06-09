package rpc

import (
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
	rwMutex sync.RWMutex
	cache   map[string]func([]reflect.Value) []reflect.Value
}

func NewRpcRegister() *rpcRegister {
	return &rpcRegister{
		cache: make(map[string]func([]reflect.Value) []reflect.Value),
	}
}

func (rpcRegister *rpcRegister) Add(name string, method func([]reflect.Value) []reflect.Value) {
	rpcRegister.rwMutex.Lock()
	defer func() {
		rpcRegister.rwMutex.Unlock()
	}()

	rpcRegister.cache[name] = method
}

func (rpcRegister *rpcRegister) Remove(name string) {
	rpcRegister.rwMutex.Lock()
	defer func() {
		rpcRegister.rwMutex.Unlock()
	}()

	delete(rpcRegister.cache, name)
}

func (rpcRegister *rpcRegister) GetRpcByName(name string) func([]reflect.Value) []reflect.Value {
	rpcRegister.rwMutex.RLock()
	defer func() {
		rpcRegister.rwMutex.RUnlock()
	}()

	return rpcRegister.cache[name]
}
