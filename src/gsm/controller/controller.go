package controller

import (
	"github.com/gsf/gsf/src/gsc/rpc"
	"github.com/gsf/gsf/src/gsf/peer"
	"reflect"
)

type Controller struct {
	PreInvoke func() bool
	EndInvoke func()
}

func NewController() *Controller {
	return &Controller{}
}

func (controller *Controller) Initialize() {
}

func (controller *Controller) Register(
	name string,
	handle func() interface{},
	before func() interface{},
	after func() interface{}) {

	if len(name) == 0 || handle == nil {
		return
	}

	method := reflect.ValueOf(handle())
	peerType := reflect.TypeOf((*peer.IPeer)(nil)).Elem()
	index := -1
	for i := 0; i < method.Type().NumIn(); i++ {
		fieldType := method.Type().In(i)
		if fieldType.Implements(peerType) {
			index = i
			break
		}
	}

	rpcRegister := rpc.GetRpcRegisterInstance()
	beforeValue := reflect.ValueOf(nil)
	if before != nil {
		beforeValue = reflect.ValueOf(before())
	}

	afterValue := reflect.ValueOf(nil)
	if after != nil {
		afterValue = reflect.ValueOf(after())
	}

	rpcRegister.Add(name,
		func(p peer.IPeer, values []reflect.Value) []reflect.Value {
			if index > 0 {
				values = append(values, make([]reflect.Value, 1)...)
				values = append(values[:index], values[index:]...)
				values[index] = reflect.ValueOf(p)
			}

			if beforeValue != reflect.ValueOf(nil) {
				beforeRet := beforeValue.Call(values)
				if beforeRet[0].Bool() {
					ret := method.Call(values)
					if afterValue != reflect.ValueOf(nil) {
						afterValue.Call(values)
					}
					return ret
				}
				return nil
			}
			ret := method.Call(values)
			return ret
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
	if controller.PreInvoke != nil && !controller.PreInvoke() {
		return nil
	}
	bytes := rpcInvoke.Request(name, args()...)
	connection := p.GetConnection()
	if connection != nil {
		connection.Send(bytes)
	}

	defer func() {
		if controller.EndInvoke != nil {
			controller.EndInvoke()
		}
	}()

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
				if controller.EndInvoke != nil {
					controller.EndInvoke()
				}
			}
		}()
		return values
	})

	rpcInvoke := rpc.NewRpcInvoke()
	if controller.PreInvoke != nil && !controller.PreInvoke() {
		rpcRegister.Remove("#" + name)
		return
	}

	bytes := rpcInvoke.Request(name, args()...)
	connection := p.GetConnection()
	if connection != nil {
		connection.Send(bytes)
	}
}
