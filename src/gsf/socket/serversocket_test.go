package socket

import (
	"github.com/sf-go/gsf/src/gsc/network"
	"github.com/sf-go/gsf/src/gsc/rpc"
	"github.com/sf-go/gsf/src/gsm/peer"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestServerSocket(t *testing.T) {
	config := network.NewNetConfig()
	config.BufferSize = 50
	config.Address = "127.0.0.1"
	config.Port = 5678
	config.ConnectTimeout = 3

	rpc.RpcRegisterInstance.Add("Func",
		func(peer peer.IPeer, args []reflect.Value) []reflect.Value {
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

	rpc.RpcRegisterInstance.Add("#Func", func(peer peer.IPeer, values []reflect.Value) []reflect.Value {
		Func := func(str string) {
			t.Log(str)
		}
		return reflect.ValueOf(Func).Call(values)
	})

	clientSocket.OnMessage = func(peer peer.IPeer, data []byte) {
		response := rpc.NewRpcResponse()
		response.Response(data, peer)
	}
	clientSocket.Connect(config)
	time.Sleep(time.Second * 3)
}
