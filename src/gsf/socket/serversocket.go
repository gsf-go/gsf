package socket

import (
	"github.com/sf-go/gsf/src/gsc/network"
	"github.com/sf-go/gsf/src/gsc/network/server"
	"github.com/sf-go/gsf/src/gsc/pool"
	"github.com/sf-go/gsf/src/gsm/dispatcher"
	"github.com/sf-go/gsf/src/gsm/peer"
)

type ServerSocket struct {
	*Event

	server      server.IServer
	objectPool  *pool.ObjectPool
	stop        func()
	dispatchers map[string]dispatcher.IDispatcher
}

func NewServerSocket(dispatcher dispatcher.IDispatcher) *ServerSocket {
	return &ServerSocket{
		server: server.NewTcpServer(network.NewHandle()),
		objectPool: pool.NewObjectPool(func() interface{} {
			return peer.NewPeer()
		}),
		Event:      NewEvent(),
		dispatcher: dispatcher,
	}
}

func (serverSocket *ServerSocket) addListener(
	connection network.IConnection,
	server *server.TcpServer) {

	p := serverSocket.objectPool.GetObject().(peer.IPeer)
	p.SetConnection(connection)

	server.OnMessage = func(connection network.IConnection, data []byte) {
		if serverSocket.OnMessage != nil {
			serverSocket.OnMessage(p, data)
			return
		}

		serverSocket.dispatcher.Dispatch(p, data)
	}

	server.OnError = func(connection network.IConnection, err error) {

	}

	server.OnDisconnected = func(connection network.IConnection, reason string) {
		if serverSocket.OnDisconnected != nil {
			serverSocket.OnDisconnected(p)
		}
	}

	if serverSocket.OnConnected != nil {
		serverSocket.OnConnected(p)
	}
}

func (serverSocket *ServerSocket) Start(config *network.NetConfig) {
	tcpServer := serverSocket.server.(*server.TcpServer)
	tcpServer.OnConnected = func(connection network.IConnection) {
		serverSocket.addListener(connection, tcpServer)
	}
	tcpServer.Accept(config)
	serverSocket.stop = func() {
		tcpServer.Close()
	}
}

func (serverSocket *ServerSocket) Stop() {
	if serverSocket.stop != nil {
		serverSocket.stop()
	}
}
