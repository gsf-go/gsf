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

func (rpcResponse *RpcResponse) Split(data []byte, args ...interface{}) (string, []byte) {
	byteReader := bytestream.NewByteReader2(data)
	messageId := ""
	byteReader.Read(&messageId)
	tmp := []byte(nil)
	byteReader.Read(&tmp)
	return messageId, tmp
}

func (rpcResponse *RpcResponse) HandleData(data []byte, args ...interface{}) []reflect.Value {
	byteReader := bytestream.NewByteReader2(data)
	deserializable := serialization.NewDeserializable(byteReader)
	result := deserializable.Deserialize(args...)
	return result
}
