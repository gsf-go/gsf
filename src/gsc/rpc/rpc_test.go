package rpc

import (
	"gsf/peer"
	"reflect"
	"strconv"
	"testing"
)

func Func(num int, str string) string {
	return strconv.Itoa(num) + " " + str
}

func TestRpc(t *testing.T) {
	method := reflect.ValueOf(Func)
	GetRpcRegisterInstance().Add("Func",
		func(peer peer.IPeer, args []reflect.Value) []reflect.Value {
			return method.Call(args)
		})
	req := NewRpcInvoke()
	ret := req.Invoke("Func", nil, 100, "xxxxx")
	t.Log(ret)
}

func TestRpc2(t *testing.T) {
	method := reflect.ValueOf(Func)
	GetRpcRegisterInstance().Add("Func",
		func(peer peer.IPeer, args []reflect.Value) []reflect.Value {
			return method.Call(args)
		})
	request := NewRpcInvoke()
	ret := request.Request("Func", 100, "xxxxx")

	response := NewRpcResponse()
	_, res := response.Response(nil, ret)
	t.Log(res)
}
