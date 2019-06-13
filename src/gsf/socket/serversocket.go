package socket

import (
	"github.com/gsf/gsf/src/gsc/network"
	"github.com/gsf/gsf/src/gsc/network/server"
	"github.com/gsf/gsf/src/gsc/pool"
	"github.com/gsf/gsf/src/gsf/peer"
)

type ServerSocket struct {
	*Event

	server     server.IServer
	objectPool *pool.ObjectPool
	stop       func()
}

func NewServerSocket() *ServerSocket {
	return &ServerSocket{
		server: server.NewTcpServer(network.NewHandle()),
		objectPool: pool.NewObjectPool(func() interface{} {
			return peer.NewPeer()
		}),
		Event: NewEvent(),
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

		GetInvokerInstance().Invoke(p, data)
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
