package socket

import (
	"gsc/network"
	"gsc/rpc"
	"gsf/peer"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestServerSocket(t *testing.T) {
	config := network.NewNetConfig()
	config.BufferSize = 50
	config.Address = "127.0.0.1"
	config.Port = 8888
	config.ConnectTimeout = 3

	rpc.GetRpcRegisterInstance().Add("Func",
		func(args []reflect.Value) []reflect.Value {
			Func := func(num int, str string) string {
				return strconv.Itoa(num) + " " + str
			}
			return reflect.ValueOf(Func).Call(args)
		})

	serverSocket := NewServerSocket()
	serverSocket.Start(config)

	clientSocket := NewClientSocket()
	clientSocket.OnConnected = func(peer peer.IPeer) {
		connection := peer.GetConnection()
		request := rpc.NewRpcInvoke()
		ret := request.Request("Func", 100, "xxxxx")
		connection.Send(ret)
	}

	clientSocket.OnMessage = func(peer peer.IPeer, data []byte) {

		rpc.GetRpcRegisterInstance().Remove("Func")
		rpc.GetRpcRegisterInstance().Add("Func", func(values []reflect.Value) []reflect.Value {
			Func := func(str string) {
				t.Log(str)
			}
			return reflect.ValueOf(Func).Call(values)
		})

		response := rpc.NewRpcResponse()
		response.Response(data)
	}
	clientSocket.Connect(config)
	time.Sleep(time.Second * 5)
}
