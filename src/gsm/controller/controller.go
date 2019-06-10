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

func (controller *Controller) Invoke(
	name string,
	peer peer.IPeer,
	args func() []interface{}) []interface{} {

	if len(name) == 0 || args == nil {
		return nil
	}

	resultChan := make(chan []interface{})
	defer func() {
		close(resultChan)
	}()

	rpcRegister := rpc.GetRpcRegisterInstance()
	rpcRegister.Add("#"+name,
		func(values []reflect.Value) []reflect.Value {

			response := make([]interface{}, len(values))
			for i, item := range values {
				response[i] = item.Interface()
			}

			defer func() {
				rpcRegister.Remove("#" + name)
				resultChan <- response
			}()
			return values
		})

	rpcInvoke := rpc.NewRpcInvoke()
	bytes := rpcInvoke.Request(name, args()...)
	connection := peer.GetConnection()
	if connection != nil {
		connection.Send(bytes)
	}

	select {
	case result := <-resultChan:
		return result
	}

}

func (controller *Controller) AsyncInvoke(
	name string,
	peer peer.IPeer,
	args func() []interface{},
	result func([]interface{})) {

	if len(name) == 0 || args == nil {
		return
	}

	rpcRegister := rpc.GetRpcRegisterInstance()
	rpcRegister.Add("#"+name,
		func(values []reflect.Value) []reflect.Value {

			response := make([]interface{}, len(values))
			for i, item := range values {
				response[i] = item.Interface()
			}

			defer func() {
				rpcRegister.Remove("#" + name)
				if result != nil {
					result(response)
				}
			}()
			return values
		})

	rpcInvoke := rpc.NewRpcInvoke()
	bytes := rpcInvoke.Request(name, args()...)
	connection := peer.GetConnection()
	if connection != nil {
		connection.Send(bytes)
	}

}
