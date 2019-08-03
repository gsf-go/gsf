package dispatcher

import (
	"bytes"
	"encoding/gob"
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
	method := dispatcher.register.GetRpcByName(messageBytes)
	if method == nil {
		logger.Log.Error("没有注册ID:", string(messageBytes), "的RPC")
		return
	}

	method(peer, messageBytes, dataBytes)
}

func (dispatcher *Dispatcher) Register(
	id interface{},
	before func() interface{},
	handle func() interface{},
	after func() interface{}) {

	if id == nil || handle == nil {
		return
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(id)
	if err != nil {
		return
	}
	idBytes := buf.Bytes()

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

	dispatcher.register.AddRequest(idBytes, func(p peer.IPeer, methodId []byte, data []byte) {
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
		dataBytes := invoke.Request(dispatcher.register.GetResponseId(methodId), values...)
		connection := p.GetConnection()
		if connection != nil {
			connection.Send(dataBytes)
		}
	})
}

func (dispatcher *Dispatcher) RawRegister(
	id interface{},
	before func(peer peer.IPeer, data []byte) bool,
	handle func(peer peer.IPeer, data []byte) []byte,
	after func(peer peer.IPeer, data []byte)) {

	if id == nil || handle == nil {
		return
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(id)
	if err != nil {
		return
	}
	idBytes := buf.Bytes()

	dispatcher.register.AddRequest(idBytes, func(p peer.IPeer, methodId []byte, data []byte) {
		if before != nil {
			if !before(p, data) {
				return
			}
		}

		if handle == nil {
			return
		}

		if len(data) == 0 {
			return
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

func (dispatcher *Dispatcher) Invoke(
	id interface{},
	p peer.IPeer,
	args func() []interface{}) []interface{} {

	if id == nil {
		return nil
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(id)
	if err != nil {
		return nil
	}
	idBytes := buf.Bytes()

	resultChan := make(chan []interface{})
	defer func() {
		close(resultChan)
	}()

	dispatcher.register.AddResponse(idBytes, func(p peer.IPeer, method []byte, data []byte) {

		response := rpc.NewRpcResponse()
		values := response.HandleData(data, p)
		if len(values) == 0 {
			return
		}

		res := make([]interface{}, len(values))
		for i, item := range values {
			res[i] = item.Interface()
		}

		dispatcher.register.RemoveResponse(idBytes)
		resultChan <- res
	})

	rpcInvoke := rpc.NewRpcInvoke()
	if dispatcher.PreInvoke != nil && !dispatcher.PreInvoke() {
		return nil
	}

	dataBytes := rpcInvoke.Request(idBytes)
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

func (dispatcher *Dispatcher) AsyncInvoke(
	id interface{},
	p peer.IPeer,
	args func() []interface{},
	result func([]interface{})) {

	if id == nil {
		return
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(id)
	if err != nil {
		return
	}

	idBytes := buf.Bytes()
	dispatcher.register.AddResponse(idBytes, func(_ peer.IPeer, methodId []byte, data []byte) {

		response := rpc.NewRpcResponse()
		values := response.HandleData(data, p)
		if len(values) == 0 {
			return
		}

		res := make([]interface{}, len(values))
		for i, item := range values {
			res[i] = item.Interface()
		}

		dispatcher.register.RemoveResponse(idBytes)
		if result != nil {
			result(res)
			if dispatcher.EndInvoke != nil {
				dispatcher.EndInvoke()
			}
		}
	})

	if dispatcher.PreInvoke != nil && !dispatcher.PreInvoke() {
		dispatcher.register.RemoveResponse(idBytes)
		return
	}

	rpcInvoke := rpc.NewRpcInvoke()
	dataBytes := rpcInvoke.Request(idBytes)
	connection := p.GetConnection()
	if connection != nil {
		connection.Send(dataBytes)
	}
}

func (dispatcher *Dispatcher) RawInvoke(
	id interface{},
	p peer.IPeer,
	data []byte) []byte {

	if id == nil {
		return nil
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(id)
	if err != nil {
		return nil
	}
	idBytes := buf.Bytes()
	resultChan := make(chan []byte)
	defer func() {
		close(resultChan)
	}()

	rpcRegister := dispatcher.register
	rpcRegister.AddResponse(idBytes, func(p peer.IPeer, methodId []byte, data []byte) {
		if len(data) == 0 {
			return
		}

		rpcRegister.RemoveResponse(idBytes)
		resultChan <- data
	})

	if dispatcher.PreInvoke != nil && !dispatcher.PreInvoke() {
		rpcRegister.RemoveResponse(idBytes)
		return nil
	}

	rpcInvoke := rpc.NewRpcInvoke()
	dataBytes := rpcInvoke.Request(idBytes)
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

func (dispatcher *Dispatcher) AsyncRawInvoke(
	id interface{},
	p peer.IPeer,
	data []byte,
	result func(data []byte)) {

	if id == nil {
		return
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(id)
	if err != nil {
		return
	}

	idBytes := buf.Bytes()
	dispatcher.register.AddResponse(idBytes, func(p peer.IPeer, methodId []byte, data []byte) {
		if len(data) == 0 {
			return
		}

		dispatcher.register.RemoveResponse(idBytes)
		if result != nil {
			result(data)
			if dispatcher.EndInvoke != nil {
				dispatcher.EndInvoke()
			}
		}
	})

	if dispatcher.PreInvoke != nil && !dispatcher.PreInvoke() {
		dispatcher.register.RemoveResponse(idBytes)
		return
	}

	rpcInvoke := rpc.NewRpcInvoke()
	dataBytes := rpcInvoke.Request(idBytes)
	connection := p.GetConnection()
	if connection != nil {
		connection.Send(append(dataBytes, data...))
	}
}
