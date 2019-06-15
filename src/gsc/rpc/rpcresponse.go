package rpc

import (
	"github.com/gsf/gsf/src/gsc/serialization"
	"reflect"
)

type IRpcResponse interface {
	Response(data []byte) []reflect.Value
}

type RpcResponse struct {
	deserializable *serialization.Deserializable
}

func NewRpcResponse() *RpcResponse {
	return &RpcResponse{
		deserializable: serialization.NewDeserializable(),
	}
}

func (rpcResponse *RpcResponse) Response(data []byte, args ...interface{}) []reflect.Value {
	result := rpcResponse.deserializable.Deserialize(data, args...)
	return result
}
