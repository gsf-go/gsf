package dispatcher

import (
	"github.com/sf-go/gsf/src/gsc/logger"
	"github.com/sf-go/gsf/src/gsc/rpc"
	"github.com/sf-go/gsf/src/gsm/component"
	"github.com/sf-go/gsf/src/gsm/peer"
	"reflect"
)

type RpcDispatcher struct {
	PreInvoke func() bool
	EndInvoke func()
	register  *rpc.RpcRegister
}

func NewRpcDispatcher(register *rpc.RpcRegister) *RpcDispatcher {
	return &RpcDispatcher{
		register: register,
	}
}

func (dispatcher *RpcDispatcher) Dispatch(peer peer.IPeer, data []byte) {
	defer func() {
		if r := recover(); r != nil {
			logger.Log.Error("Recovered in %s", r)
		}
	}()

	response := rpc.NewRpcResponse()
	methodId, result := response.Response(data, peer)
	if len(result) == 0 {
		return
	}

	method := dispatcher.register.GetRpcByName(string(methodId))
	if method == nil {
		logger.Log.Error("没有注册ID:", string(methodId), "的RPC")
		return
	}

	value := method(peer, result)
	values := make([]interface{}, len(value))
	for i, item := range value {
		values[i] = item.Interface()
	}

	invoke := rpc.NewRpcInvoke(dispatcher.register)
	methodId[len(methodId)-1] = 1
	bytes := invoke.Request(methodId, values...)
	connection := peer.GetConnection()
	if connection != nil {
		connection.Send(bytes)
	}
}

func (dispatcher *RpcDispatcher) Register(
	id []byte,
	handle func() interface{},
	before func() interface{},
	after func() interface{}) {

	if len(id) == 0 || handle == nil {
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

	rpcRegister := dispatcher.register
	beforeValue := reflect.ValueOf(nil)
	if before != nil {
		beforeValue = reflect.ValueOf(before())
	}

	afterValue := reflect.ValueOf(nil)
	if after != nil {
		afterValue = reflect.ValueOf(after())
	}

	rpcRegister.AddRequest(id, func(p peer.IPeer, values []reflect.Value) []reflect.Value {

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

func (dispatcher *RpcDispatcher) Invoke(
	id []byte,
	p peer.IPeer,
	args func() []interface{}) []interface{} {

	if len(id) == 0 || args == nil {
		return nil
	}

	resultChan := make(chan []interface{})
	defer func() {
		close(resultChan)
	}()

	rpcRegister := dispatcher.register
	rpcRegister.AddResponse(id, func(_ peer.IPeer, values []reflect.Value) []reflect.Value {

		response := make([]interface{}, len(values))
		for i, item := range values {
			response[i] = item.Interface()
		}

		defer func() {
			rpcRegister.RemoveResponse(id)
			resultChan <- response
		}()
		return values
	})

	rpcInvoke := rpc.NewRpcInvoke(dispatcher.register)
	if dispatcher.PreInvoke != nil && !dispatcher.PreInvoke() {
		return nil
	}

	components := make([]interface{}, 0)
	p.Range(func(key string, component component.IComponent) bool {
		components = append(components, component)
		return true
	})

	bytes := rpcInvoke.Request(id, append(args(), components...)...)
	connection := p.GetConnection()
	if connection != nil {
		connection.Send(bytes)
	}

	defer func() {
		if dispatcher.EndInvoke != nil {
			dispatcher.EndInvoke()
		}
	}()

	select {
	case result := <-resultChan:
		return result
	}
}

func (dispatcher *RpcDispatcher) AsyncInvoke(
	id []byte,
	p peer.IPeer,
	args func() []interface{},
	result func([]interface{})) {

	if len(id) == 0 || args == nil {
		return
	}

	rpcRegister := dispatcher.register
	rpcRegister.AddResponse(id, func(_ peer.IPeer, values []reflect.Value) []reflect.Value {

		response := make([]interface{}, len(values))
		for i, item := range values {
			response[i] = item.Interface()
		}

		defer func() {
			rpcRegister.RemoveResponse(id)
			if result != nil {
				result(response)
				if dispatcher.EndInvoke != nil {
					dispatcher.EndInvoke()
				}
			}
		}()
		return values
	})

	rpcInvoke := rpc.NewRpcInvoke(dispatcher.register)
	if dispatcher.PreInvoke != nil && !dispatcher.PreInvoke() {
		rpcRegister.RemoveResponse(id)
		return
	}

	components := make([]interface{}, 0)
	p.Range(func(key string, component component.IComponent) bool {
		components = append(components, component)
		return true
	})

	bytes := rpcInvoke.Request(id, append(args(), components...)...)
	connection := p.GetConnection()
	if connection != nil {
		connection.Send(bytes)
	}
}
