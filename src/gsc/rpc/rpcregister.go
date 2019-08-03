package rpc

import (
	"github.com/sf-go/gsf/src/gsm/peer"
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
	method func(peer peer.IPeer, messageId []byte, data []byte)) {
	rpcRegister.cache.Store(string(id), method)
}

func (rpcRegister *RpcRegister) AddResponse(id []byte,
	method func(peer peer.IPeer, messageId []byte, data []byte)) {

	id = append(id, 1)
	rpcRegister.cache.Store(string(id), method)
}

func (rpcRegister *RpcRegister) GetResponseId(id []byte) []byte {
	id = append(id, 1)
	return id
}

func (rpcRegister *RpcRegister) RemoveResponse(id []byte) {
	id = append(id, 1)
	rpcRegister.cache.Delete(string(id))
}

func (rpcRegister *RpcRegister) GetRpcByName(id []byte) func(peer peer.IPeer, messageId []byte, data []byte) {
	value, ok := rpcRegister.cache.Load(string(id))
	if ok {
		return value.(func(peer peer.IPeer, messageId []byte, data []byte))
	}
	return nil
}
