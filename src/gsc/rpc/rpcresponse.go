package rpc

import (
	"gsc/serialization"
	"reflect"
)

type IRpcResponse interface {
	Response(data []byte) (string, []reflect.Value)
}

type RpcResponse struct {
	deserializable serialization.IDeserializable
}

func NewRpcResponse() *RpcResponse {
	return &RpcResponse{
		deserializable: new(serialization.Deserializable),
	}
}

func (rpcResponse *RpcResponse) Response(data []byte) (string, []reflect.Value) {
	result := rpcResponse.deserializable.Deserialize(data)
	messageId := result[0].String()
	method := GetRpcRegisterInstance().GetRpcByName(messageId)
	return messageId, method(result[1:])
}
