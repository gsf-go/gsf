package socket

import (
	"gsc/rpc"
	"gsf/peer"
	"strings"
	"sync"
)

//go:generate ../../gen/Singleton.exe -struct=invoker -out=./invoker.go

var invokerInstance *invoker
var invokerOnce sync.Once

func GetInvokerInstance() *invoker {
	invokerOnce.Do(func() {
		invokerInstance = NewInvoker()
	})
	return invokerInstance
}

type invoker struct {
}

func NewInvoker() *invoker {
	return &invoker{}
}

func (invoker *invoker) Invoke(peer peer.IPeer, data []byte) {
	response := rpc.NewRpcResponse()
	messageId, ret := response.Response(data)

	values := make([]interface{}, len(ret))
	for i, item := range ret {
		values[i] = item.Interface()
	}

	if strings.Contains(messageId, "#") {
		return
	}

	invoke := rpc.NewRpcInvoke()
	bytes := invoke.Request("#"+messageId, values...)
	connection := peer.GetConnection()
	if connection != nil {
		connection.Send(bytes)
	}
}
