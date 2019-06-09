package controller

import (
	"gsc/rpc"
	"gsf/peer"
	"reflect"
)

type Controller struct {
}

func (controller *Controller) Initialize() {

}

func (controller *Controller) Register(name string, function func() interface{}) {
	if len(name) == 0 || function == nil {
		return
	}

	method := reflect.ValueOf(function())
	rpc.GetRpcRegisterInstance().Add(name,
		func(values []reflect.Value) []reflect.Value {
			return method.Call(values)
		})
}

func (controller *Controller) Invoke(name string, peer peer.IPeer, function func() []interface{}) {
	if len(name) == 0 || function == nil {
		return
	}

	rpcInvoke := rpc.NewRpcInvoke()
	bytes := rpcInvoke.Request(name, function()...)
	connection := peer.GetConnection()
	if connection != nil {
		connection.Send(bytes)
	}
}
