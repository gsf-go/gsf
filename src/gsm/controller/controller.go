package controller

import (
	"github.com/gsf/gsf/src/gsc/rpc"
	"github.com/gsf/gsf/src/gsf/peer"
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
	peerType := reflect.TypeOf((*peer.IPeer)(nil)).Elem()
	index := -1
	for i := 0; i < method.Type().NumIn(); i++ {
		fieldType := method.Type().In(i)
		if fieldType.Implements(peerType) {
			index = i
			break
		}
	}

	rpc.GetRpcRegisterInstance().Add(name,
		func(p peer.IPeer, values []reflect.Value) []reflect.Value {
			if index > 0 {
				values = append(values, make([]reflect.Value, 1)...)
				values = append(values[:index], values[index:]...)
				values[index] = reflect.ValueOf(p)
			}
			return method.Call(values)
		})
}

func (controller *Controller) Invoke(
	name string,
	p peer.IPeer,
	args func() []interface{}) []interface{} {

	if len(name) == 0 || args == nil {
		return nil
	}

	resultChan := make(chan []interface{})
	defer func() {
		close(resultChan)
	}()

	rpcRegister := rpc.GetRpcRegisterInstance()
	rpcRegister.Add("#"+name, func(_ peer.IPeer, values []reflect.Value) []reflect.Value {

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
	connection := p.GetConnection()
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
	p peer.IPeer,
	args func() []interface{},
	result func([]interface{})) {

	if len(name) == 0 || args == nil {
		return
	}

	rpcRegister := rpc.GetRpcRegisterInstance()
	rpcRegister.Add("#"+name, func(_ peer.IPeer, values []reflect.Value) []reflect.Value {

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
	connection := p.GetConnection()
	if connection != nil {
		connection.Send(bytes)
	}

}
