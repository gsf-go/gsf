package rpc

import (
	"github.com/sf-go/gsf/src/gsc/bytestream"
	"github.com/sf-go/gsf/src/gsc/serialization"
	"github.com/sf-go/gsf/src/gsm/peer"
	"reflect"
)

type IRpcInvoke interface {
	Invoke(methodId []byte, peer peer.IPeer, args ...interface{}) []reflect.Value
	Request(methodId []byte, args ...interface{}) []byte
}

type RpcInvoke struct {
	serializable serialization.ISerializable
}

func NewRpcInvoke() *RpcInvoke {
	return &RpcInvoke{
		serializable: serialization.NewSerializable(),
	}
}

func (rpcInvoke *RpcInvoke) Request(methodId string, args ...interface{}) []byte {
	dataBytes := rpcInvoke.serializable.Serialize(args...)
	writer := bytestream.NewByteWriter3()
	writer.Write(methodId)
	writer.Write(dataBytes)
	return writer.ToBytes()
}
