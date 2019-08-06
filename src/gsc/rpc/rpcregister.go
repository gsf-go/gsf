package rpc

import (
	"github.com/sf-go/gsf/src/gsm/peer"
)

//go:generate ../../../gen/Singleton.exe -struct=rpcRegister -out=./rpcregister.go

type RpcRegister struct {
	cache map[string]func(peer.IPeer, string, []byte)
}

func NewRpcRegister() *RpcRegister {
	return &RpcRegister{
		cache: make(map[string]func(peer.IPeer, string, []byte)),
	}
}

func (rpcRegister *RpcRegister) AddRequest(id string,
	method func(peer peer.IPeer, messageId string, data []byte)) {
	rpcRegister.cache[id] = method
}

func (rpcRegister *RpcRegister) AddResponse(id string,
	method func(peer peer.IPeer, messageId string, data []byte)) {
	rpcRegister.cache["#"+id] = method
}

func (rpcRegister *RpcRegister) GetResponseId(id string) string {
	return "#" + id
}

func (rpcRegister *RpcRegister) RemoveResponse(id string) {
	delete(rpcRegister.cache, "#"+id)
}

func (rpcRegister *RpcRegister) GetRpcByName(id string) func(peer peer.IPeer, messageId string, data []byte) {
	value, ok := rpcRegister.cache[id]
	if ok {
		return value
	}
	return nil
}
