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

func (rpcResponse *RpcResponse) HandleMessageId(data []byte, args ...interface{}) ([]byte, []byte) {
	byteReader := bytestream.NewByteReader2(data)
	deserializable := serialization.NewDeserializable(byteReader)
	bytes := deserializable.DeserializeSingle(args...).Interface()
	length := byteReader.GetUnreadLength()
	return bytes.([]byte), data[len(data)-length:]
}

func (rpcResponse *RpcResponse) HandleData(data []byte, args ...interface{}) []reflect.Value {
	byteReader := bytestream.NewByteReader2(data)
	deserializable := serialization.NewDeserializable(byteReader)
	result := deserializable.Deserialize(args...)
	return result
}
