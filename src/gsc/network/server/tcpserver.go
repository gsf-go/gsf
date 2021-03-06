package server

import (
	"context"
	"github.com/sf-go/gsf/src/gsc/logger"
	"github.com/sf-go/gsf/src/gsc/network"
	"github.com/sf-go/gsf/src/gsc/pool"
	"net"
	"strconv"
	"time"
)

type TcpServer struct {
	*network.NetEvent

	handle  network.IHandle
	context context.Context

	_cancel func()
}

func NewTcpServer(
	handle network.IHandle) *TcpServer {
	ctx, cancel := context.WithCancel(context.Background())
	return &TcpServer{
		handle:   handle,
		_cancel:  cancel,
		context:  ctx,
		NetEvent: network.NewNetEvent(),
	}
}

func (tcpServer *TcpServer) Accept(config *network.NetConfig) {

	bufferPool := pool.NewBytePool(config.BufferSize)
	address := config.Address + ":" + strconv.Itoa(int(config.Port))
	listener, err := net.Listen("tcp", address)

	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	closeDone := make(chan func() (network.IConnection, error, string), 1)

	go func() {
		for {
			select {
			case <-tcpServer.context.Done():
				return
			default:
				conn, err := listener.Accept()
				if err != nil {
					logger.Log.Error(err.Error())
					continue
				}

				go tcpServer.handleClient(
					bufferPool,
					config,
					closeDone,
					conn)
			}
		}
	}()

	go func() {
		for {
			select {
			case <-tcpServer.context.Done():
				return
			case callback := <-closeDone:
				connection, err, reason := callback()

				tcpServer.OnDisconnected(connection, reason)

				if tcpServer.OnError != nil {
					tcpServer.OnError(connection, err)
				}
			}
		}
	}()
}

func (tcpServer *TcpServer) handleClient(
	bufferPool *pool.BytePool,
	config *network.NetConfig,
	closeDone chan<- func() (network.IConnection, error, string),
	conn net.Conn) {

	connection := network.NewConnection(
		conn.Write,
		tcpServer.handle,
		closeDone)

	if tcpServer.OnConnected != nil {
		tcpServer.OnConnected(connection)
	}

	buffer := bufferPool.GetBuffer()

	defer func() {
		_ = conn.Close()
		bufferPool.Recycle(buffer)
		connection.Close()
	}()

	byteArray := buffer.Bytes()
	packet := &network.Packet{
		Config:     config,
		Buffer:     byteArray,
		Connection: connection,
	}

	for {
		logger.Log.Info("Offset " + strconv.Itoa(int(packet.Offset)))
		n, err := conn.Read(byteArray[packet.Offset:])
		logger.Log.Info("Recv %d", n)

		if err != nil {
			closeDone <- func() (network.IConnection, error, string) {
				return connection, err, "Error"
			}
			return
		}

		if n == 0 {
			closeDone <- func() (network.IConnection, error, string) {
				return connection, nil, "EOF"
			}
			return
		}

		deadline := time.Now().Add(15 * time.Second)
		err = conn.SetWriteDeadline(deadline)
		if err != nil {
			closeDone <- func() (network.IConnection, error, string) {
				return connection, err, "Error"
			}
			return
		}

		packet.Offset += uint16(n)
		tcpServer.handle.ReadHandle(packet, tcpServer.post)
	}
}

func (tcpServer *TcpServer) post(connection network.IConnection, data []byte) {
	if tcpServer.OnMessage != nil {
		tcpServer.OnMessage(connection, data)
	}
}

func (tcpServer *TcpServer) Close() {
	tcpServer._cancel()
}
