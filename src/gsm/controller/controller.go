package controller

import (
	"github.com/sf-go/gsf/src/gsc/rpc"
	"github.com/sf-go/gsf/src/gsf/peer"
	"github.com/sf-go/gsf/src/gsm/component"
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
	argsLength := method.Type().NumIn()
	for i := 0; i < argsLength; i++ {
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

			args := values[:argsLength]

			if beforeValue != reflect.ValueOf(nil) {
				beforeRet := beforeValue.Call(args)
				if beforeRet[0].Bool() {
					ret := method.Call(args)
					if afterValue != reflect.ValueOf(nil) {
						afterValue.Call(args)
					}
					return ret
				}
				return nil
			}
			ret := method.Call(args)
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

	components := make([]interface{}, 0)
	p.Range(func(key string, component component.IComponent) bool {
		components = append(components, component)
		return true
	})

	bytes := rpcInvoke.Request(name, append(args(), components...)...)
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

	components := make([]interface{}, 0)
	p.Range(func(key string, component component.IComponent) bool {
		components = append(components, component)
		return true
	})

	bytes := rpcInvoke.Request(name, append(args(), components...)...)
	connection := p.GetConnection()
	if connection != nil {
		connection.Send(bytes)
	}
}
