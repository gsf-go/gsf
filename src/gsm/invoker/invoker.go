package invoker

import (
	"github.com/sf-go/gsf/src/gsc/logger"
	"github.com/sf-go/gsf/src/gsc/rpc"
	"github.com/sf-go/gsf/src/gsm/peer"
	"reflect"
)

type Invoker struct {
	PreInvoke func() bool
	EndInvoke func()
	register  *rpc.RpcRegister
}

func NewInvoker(register *rpc.RpcRegister) *Invoker {
	return &Invoker{
		register: register,
	}
}

func (dispatcher *Invoker) Dispatch(peer peer.IPeer, data []byte) {
	defer func() {
		if r := recover(); r != nil {
			logger.Log.Error("Recovered in %s", r)
		}
	}()

	response := rpc.NewRpcResponse()
	messageId, dataBytes := response.Split(data, peer)
	method := dispatcher.register.GetRpcByName(messageId)
	if method == nil {
		logger.Log.Error("没有注册ID:", messageId, "的RPC")
		return
	}

	method(peer, messageId, dataBytes)
}

func (dispatcher *Invoker) Register(
	id string,
	before func() interface{},
	handle func() interface{},
	after func() interface{}) {

	if len(id) == 0 || handle == nil {
		return
	}

	// 获取IPeer特殊字段
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

	beforeValue := reflect.ValueOf(nil)
	if before != nil {
		beforeValue = reflect.ValueOf(before())
	}

	afterValue := reflect.ValueOf(nil)
	if after != nil {
		afterValue = reflect.ValueOf(after())
	}

	dispatcher.register.AddRequest(id, func(p peer.IPeer, methodId string, data []byte) {
		response := rpc.NewRpcResponse()
		result := response.HandleData(data, p)
		if len(result) == 0 {
			return
		}

		if index >= 0 {
			result = append(result, reflect.ValueOf(nil))
			copy(result[index+1:], result[index:len(result)-1])
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
		dataBytes := invoke.Request(dispatcher.register.GetResponseId(methodId), values...)
		connection := p.GetConnection()
		if connection != nil {
			connection.Send(dataBytes)
		}
	})
}

func (dispatcher *Invoker) FixRegister(
	id string,
	handle func(peer peer.IPeer,
		args ...interface{}) []interface{}) {

	if len(id) == 0 || handle == nil {
		return
	}

	dispatcher.register.AddRequest(id, func(p peer.IPeer, id string, data []byte) {
		response := rpc.NewRpcResponse()
		args := response.HandleData(data, p)
		if len(args) == 0 {
			return
		}

		tmp := make([]interface{}, 0)
		for _, v := range args {
			tmp = append(tmp, v.Interface())
		}
		result := handle(p, tmp...)
		invoke := rpc.NewRpcInvoke()
		dataBytes := invoke.Request(dispatcher.register.GetResponseId(id), result...)
		connection := p.GetConnection()
		if connection != nil {
			connection.Send(dataBytes)
		}
	})
}

func (dispatcher *Invoker) RawRegister(
	id string,
	before func(peer peer.IPeer, data []byte) bool,
	handle func(peer peer.IPeer, data []byte) []byte,
	after func(peer peer.IPeer, data []byte)) {

	if len(id) == 0 || handle == nil {
		return
	}

	dispatcher.register.AddRequest(id, func(p peer.IPeer, methodId string, data []byte) {
		if handle == nil || len(data) == 0 {
			return
		}

		if before != nil {
			if !before(p, data) {
				return
			}
		}
		values := handle(p, data)
		if after != nil {
			after(p, data)
		}

		rpcInvoke := rpc.NewRpcInvoke()
		dataBytes := rpcInvoke.Request(dispatcher.register.GetResponseId(methodId))
		connection := p.GetConnection()
		if connection != nil {
			connection.Send(append(dataBytes, values...))
		}
	})
}

func (dispatcher *Invoker) Invoke(
	id string,
	p peer.IPeer,
	args func() []interface{}) []interface{} {

	if len(id) == 0 {
		return nil
	}

	resultChan := make(chan []interface{})
	defer func() {
		close(resultChan)
	}()

	dispatcher.register.AddResponse(id, func(p peer.IPeer, methodId string, data []byte) {

		response := rpc.NewRpcResponse()
		values := response.HandleData(data, p)
		if len(values) == 0 {
			return
		}

		res := make([]interface{}, len(values))
		for i, item := range values {
			res[i] = item.Interface()
		}

		dispatcher.register.RemoveResponse(methodId)
		resultChan <- res
	})

	rpcInvoke := rpc.NewRpcInvoke()
	if dispatcher.PreInvoke != nil && !dispatcher.PreInvoke() {
		return nil
	}

	dataBytes := rpcInvoke.Request(id, args()...)
	connection := p.GetConnection()
	if connection != nil {
		connection.Send(dataBytes)
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

func (dispatcher *Invoker) AsyncInvoke(
	id string,
	p peer.IPeer,
	args func() []interface{},
	result func([]interface{})) {

	if len(id) == 0 {
		return
	}

	dispatcher.register.AddResponse(id, func(_ peer.IPeer, methodId string, data []byte) {

		response := rpc.NewRpcResponse()
		values := response.HandleData(data, p)
		if len(values) == 0 {
			return
		}

		res := make([]interface{}, len(values))
		for i, item := range values {
			res[i] = item.Interface()
		}

		dispatcher.register.RemoveResponse(methodId)
		if result != nil {
			result(res)
			if dispatcher.EndInvoke != nil {
				dispatcher.EndInvoke()
			}
		}
	})

	if dispatcher.PreInvoke != nil && !dispatcher.PreInvoke() {
		dispatcher.register.RemoveResponse(id)
		return
	}

	rpcInvoke := rpc.NewRpcInvoke()
	dataBytes := rpcInvoke.Request(id, args()...)
	connection := p.GetConnection()
	if connection != nil {
		connection.Send(dataBytes)
	}
}

func (dispatcher *Invoker) RawInvoke(
	id string,
	p peer.IPeer,
	data []byte) []byte {

	if len(id) == 0 {
		return nil
	}

	resultChan := make(chan []byte)
	defer func() {
		close(resultChan)
	}()

	rpcRegister := dispatcher.register
	rpcRegister.AddResponse(id, func(p peer.IPeer, methodId string, data []byte) {
		if len(data) == 0 {
			return
		}

		rpcRegister.RemoveResponse(methodId)
		resultChan <- data
	})

	if dispatcher.PreInvoke != nil && !dispatcher.PreInvoke() {
		rpcRegister.RemoveResponse(id)
		return nil
	}

	rpcInvoke := rpc.NewRpcInvoke()
	dataBytes := rpcInvoke.Request(id)
	connection := p.GetConnection()
	if connection != nil {
		connection.Send(append(dataBytes, data...))
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

func (dispatcher *Invoker) AsyncRawInvoke(
	id string,
	p peer.IPeer,
	data []byte,
	result func(data []byte)) {

	if len(id) == 0 {
		return
	}

	dispatcher.register.AddResponse(id, func(p peer.IPeer, methodId string, data []byte) {
		if len(data) == 0 {
			return
		}

		dispatcher.register.RemoveResponse(methodId)
		if result != nil {
			result(data)
			if dispatcher.EndInvoke != nil {
				dispatcher.EndInvoke()
			}
		}
	})

	if dispatcher.PreInvoke != nil && !dispatcher.PreInvoke() {
		dispatcher.register.RemoveResponse(id)
		return
	}

	rpcInvoke := rpc.NewRpcInvoke()
	dataBytes := rpcInvoke.Request(id)
	connection := p.GetConnection()
	if connection != nil {
		connection.Send(append(dataBytes, data...))
	}
}
