package server

import (
	"gsc/network"
	"gsc/network/client"
	"testing"
	"time"
)

func TestTcp(t *testing.T) {
	config := network.NewNetConfig()
	config.BufferSize = 10
	config.Address = "127.0.0.1"
	config.Port = 8888
	config.ConnectTimeout = 3

	handle := network.NewHandle()

	tcpServer := NewTcpServer(handle)
	tcpServer.Accept(config)
	tcpServer.OnMessage = func(connection network.IConnection, data []byte) {
		t.Logf(string(data))
		connection.Send([]byte("response"))
	}

	tcpClient := client.NewTcpClient(handle)
	tcpClient.OnMessage = func(connecton network.IConnection, data []byte) {
		t.Logf(string(data))
	}

	tcpClient.OnConnected = func(connection network.IConnection) {
		if connection != nil {
			connection.Send([]byte("test"))
		}
	}
	tcpClient.Connect(config)

	time.Sleep(time.Second * 5)
}
