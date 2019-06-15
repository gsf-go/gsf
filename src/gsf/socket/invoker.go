package socket

import (
	"github.com/gsf/gsf/src/gsc/logger"
	"github.com/gsf/gsf/src/gsc/rpc"
	"github.com/gsf/gsf/src/gsf/peer"
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
	defer func() {
		if r := recover(); r != nil {
			logger.Log.Error("Recovered in %s", r)
		}
	}()

	response := rpc.NewRpcResponse()
	result := response.Response(data, peer)

	if len(result) == 0 {
		return
	}
	messageId := result[0].String()
	method := rpc.GetRpcRegisterInstance().GetRpcByName(messageId)
	if method == nil {
		logger.Log.Error("没有注册ID:" + messageId + "的RPC")
		return
	}

	value := method(peer, result[1:])
	values := make([]interface{}, len(value))
	for i, item := range value {
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
