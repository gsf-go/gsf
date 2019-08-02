package socket

import (
	"github.com/sf-go/gsf/src/gsc/network"
	"github.com/sf-go/gsf/src/gsc/network/client"
	"github.com/sf-go/gsf/src/gsc/pool"
	"github.com/sf-go/gsf/src/gsm/dispatcher"
	"github.com/sf-go/gsf/src/gsm/peer"
)

type ClientSocket struct {
	*Event

	client     client.IClient
	objectPool *pool.ObjectPool
	dispatcher dispatcher.IDispatcher
}

func NewClientSocket(dispatcher dispatcher.IDispatcher) *ClientSocket {
	return &ClientSocket{
		client: client.NewTcpClient(network.NewHandle()),
		objectPool: pool.NewObjectPool(func() interface{} {
			return peer.NewPeer()
		}),
		Event:      NewEvent(),
		dispatcher: dispatcher,
	}
}

func (clientSocket *ClientSocket) addListener(
	connection network.IConnection,
	client *client.TcpClient) {

	p := clientSocket.objectPool.GetObject().(peer.IPeer)
	p.SetConnection(connection)

	client.OnMessage = func(connection network.IConnection, data []byte) {
		if clientSocket.OnMessage != nil {
			clientSocket.OnMessage(p, data)
			return
		}

		clientSocket.dispatcher.Dispatch(p, data)
	}

	client.OnError = func(connection network.IConnection, err error) {

	}

	client.OnDisconnected = func(connection network.IConnection, reason string) {
		if clientSocket.OnDisconnected != nil {
			clientSocket.OnDisconnected(p)
		}
	}

	if clientSocket.OnConnected != nil {
		clientSocket.OnConnected(p)
	}
}

func (clientSocket *ClientSocket) Connect(config *network.NetConfig) {
	tcpClient := clientSocket.client.(*client.TcpClient)
	tcpClient.OnConnected = func(connection network.IConnection) {
		clientSocket.addListener(connection, tcpClient)
	}
	tcpClient.Connect(config)
}
