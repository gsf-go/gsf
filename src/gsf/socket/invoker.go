package socket

import (
	"github.com/sf-go/gsf/src/gsc/logger"
	"github.com/sf-go/gsf/src/gsc/rpc"
	"github.com/sf-go/gsf/src/gsf/peer"
	"strings"
)

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
	methodId, result := response.Response(data, peer)
	if len(result) == 0 {
		return
	}

	method := rpc.GetRpcRegisterInstance().GetRpcByName(methodId)
	if method == nil {
		logger.Log.Error("没有注册ID:" + methodId + "的RPC")
		return
	}

	value := method(peer, result)
	values := make([]interface{}, len(value))
	for i, item := range value {
		values[i] = item.Interface()
	}

	if strings.Contains(methodId, "#") {
		return
	}

	invoke := rpc.NewRpcInvoke()
	bytes := invoke.Request("#"+methodId, values...)
	connection := peer.GetConnection()
	if connection != nil {
		connection.Send(bytes)
	}
}
