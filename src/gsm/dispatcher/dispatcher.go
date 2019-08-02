package dispatcher

import (
	"github.com/sf-go/gsf/src/gsc/logger"
	"github.com/sf-go/gsf/src/gsc/rpc"
	"github.com/sf-go/gsf/src/gsm/peer"
	"reflect"
)

type Dispatcher struct {
	PreInvoke func() bool
	EndInvoke func()
	register  *rpc.RpcRegister
}

func NewDispatcher(register *rpc.RpcRegister) *Dispatcher {
	return &Dispatcher{
		register: register,
	}
}

func (dispatcher *Dispatcher) Dispatch(peer peer.IPeer, data []byte) {
	defer func() {
		if r := recover(); r != nil {
			logger.Log.Error("Recovered in %s", r)
		}
	}()

	response := rpc.NewRpcResponse()
	messageBytes, dataBytes := response.HandleMessageId(data, peer)
	messageId := string(messageBytes)
	method := dispatcher.register.GetRpcByName(messageId)
	if method == nil {
		logger.Log.Error("没有注册ID:", messageId, "的RPC")
		return
	}

	method(peer, messageBytes, dataBytes)
}

func (dispatcher *Dispatcher) Register(
	id []byte,
	before func() interface{},
	handle func() interface{},
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

	rpcRegister.AddRequest(id, func(p peer.IPeer, methodId []byte, data []byte) {
		response := rpc.NewRpcResponse()
		result := response.HandleData(data, p)
		if len(result) == 0 {
			return
		}

		if index > 0 {
			result = append(result, make([]reflect.Value, 1)...)
			result = append(result[:index], result[index:]...)
			result[index] = reflect.ValueOf(p)
		}

		args := result[:argsLength]
		if beforeValue != reflect.ValueOf(nil) {
			beforeRet := beforeValue.Call(args)
			if !beforeRet[0].Bool() {
				return
			}
		}

		ret := method.Call(args)
		if afterValue != reflect.ValueOf(nil) {
			afterValue.Call(args)
		}

		values := make([]interface{}, len(ret))
		for i, item := range ret {
			values[i] = item.Interface()
		}

		invoke := rpc.NewRpcInvoke()
		methodId[len(methodId)-1] = 1
		bytes := invoke.Request(methodId, values...)
		connection := p.GetConnection()
		if connection != nil {
			connection.Send(bytes)
		}
	})
}

func (dispatcher *Dispatcher) RawRegister(
	id []byte,
	before func(peer peer.IPeer, data []byte) bool,
	handle func(peer peer.IPeer, data []byte) []interface{},
	after func(peer peer.IPeer, data []byte)) {

	if len(id) == 0 || handle == nil {
		return
	}

	dispatcher.register.AddRequest(id, func(p peer.IPeer, methodId []byte, data []byte) {
		if before != nil {
			if !before(p, data) {
				return
			}
		}

		if handle == nil {
			return
		}

		values := handle(p, data)
		if after != nil {
			after(p, data)
		}

		invoke := rpc.NewRpcInvoke()
		methodId[len(methodId)-1] = 1
		bytes := invoke.Request(methodId, values...)
		connection := p.GetConnection()
		if connection != nil {
			connection.Send(bytes)
		}
	})
}

func (dispatcher *Dispatcher) Invoke(
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
	rpcRegister.AddResponse(id, func(_ peer.IPeer, method []byte, data []byte) {

		response := rpc.NewRpcResponse()
		values := response.HandleData(data, p)
		if len(values) == 0 {
			return
		}

		res := make([]interface{}, len(values))
		for i, item := range values {
			res[i] = item.Interface()
		}

		rpcRegister.RemoveResponse(id)
		resultChan <- res
	})

	rpcInvoke := rpc.NewRpcInvoke()
	if dispatcher.PreInvoke != nil && !dispatcher.PreInvoke() {
		return nil
	}

	components := make([]interface{}, 0)
	//p.Range(func(key string, component component.IComponent) bool {
	//	components = append(components, component)
	//	return true
	//})

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

func (dispatcher *Dispatcher) AsyncInvoke(
	id []byte,
	p peer.IPeer,
	args func() []interface{},
	result func([]interface{})) {

	if len(id) == 0 || args == nil {
		return
	}

	rpcRegister := dispatcher.register
	rpcRegister.AddResponse(id, func(_ peer.IPeer, methodId []byte, data []byte) {

		response := rpc.NewRpcResponse()
		values := response.HandleData(data, p)
		if len(values) == 0 {
			return
		}

		res := make([]interface{}, len(values))
		for i, item := range values {
			res[i] = item.Interface()
		}

		rpcRegister.RemoveResponse(id)
		if result != nil {
			result(res)
			if dispatcher.EndInvoke != nil {
				dispatcher.EndInvoke()
			}
		}
	})

	rpcInvoke := rpc.NewRpcInvoke()
	if dispatcher.PreInvoke != nil && !dispatcher.PreInvoke() {
		rpcRegister.RemoveResponse(id)
		return
	}

	components := make([]interface{}, 0)
	//p.Range(func(key string, component component.IComponent) bool {
	//	components = append(components, component)
	//	return true
	//})

	bytes := rpcInvoke.Request(id, append(args(), components...)...)
	connection := p.GetConnection()
	if connection != nil {
		connection.Send(bytes)
	}
}
