package rpc

import (
	"github.com/sf-go/gsf/src/gsc/bytestream"
	"github.com/sf-go/gsf/src/gsc/serialization"
	"reflect"
)

type IRpcResponse interface {
	Response(data []byte) []reflect.Value
}

type RpcResponse struct {
}

func NewRpcResponse() *RpcResponse {
	return &RpcResponse{}
}

func (rpcResponse *RpcResponse) Response(data []byte, args ...interface{}) (string, []reflect.Value) {
	byteReader := bytestream.NewByteReader2(data)
	deserializable := serialization.NewDeserializable(byteReader)
	result := deserializable.Deserialize(args...)
	methodId := result[0].Interface().(string)
	return methodId, result[1:]
}
