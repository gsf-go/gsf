package rpc

import (
	"gsc/logger"
	"gsc/serialization"
	"gsf/peer"
	"reflect"
)

type IRpcResponse interface {
	Response(peer peer.IPeer, data []byte) (string, []reflect.Value)
}

type RpcResponse struct {
	deserializable *serialization.Deserializable
}

func NewRpcResponse() *RpcResponse {
	return &RpcResponse{
		deserializable: serialization.NewDeserializable(),
	}
}

func (rpcResponse *RpcResponse) Response(peer peer.IPeer, data []byte) (string, []reflect.Value) {
	result := rpcResponse.deserializable.Deserialize(data)
	messageId := result[0].String()
	method := GetRpcRegisterInstance().GetRpcByName(messageId)
	if method == nil {
		logger.Log.Error("没有注册ID:" + messageId + "的RPC")
		return messageId, nil
	}
	return messageId, method(peer, result[1:])
}
