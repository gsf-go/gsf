package server

import (
	"context"
	"github.com/gsf/gsf/src/gsc/logger"
	"github.com/gsf/gsf/src/gsc/network"
	"github.com/gsf/gsf/src/gsc/pool"
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
		select {
		case callback := <-closeDone:
			connection, err, reason := callback()

			if err != nil {
				tcpServer.OnDisconnected(connection, reason)
			} else {
				tcpServer.OnDisconnected(connection, "Error")
			}

			if tcpServer.OnError != nil {
				tcpServer.OnError(connection, err)
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
	offset := uint16(0)

	for {
		n, err := conn.Read(byteArray[offset:])

		if err != nil {
			closeDone <- func() (network.IConnection, error, string) {
				return connection, err, ""
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
				return connection, err, ""
			}
			return
		}

		offset += tcpServer.handleData(
			config,
			connection,
			byteArray[0:uint16(n)+offset])
	}
}

func (tcpServer *TcpServer) handleData(
	config *network.NetConfig,
	connection network.IConnection,
	buffer []byte) uint16 {

	packet := &network.Packet{
		Config:     config,
		Buffer:     buffer,
		Connection: connection,
	}

	return tcpServer.handle.ReadHandle(
		packet,
		tcpServer.post)
}

func (tcpServer *TcpServer) post(connection network.IConnection, data []byte) {
	if tcpServer.OnMessage != nil {
		tcpServer.OnMessage(connection, data)
	}
}

func (tcpServer *TcpServer) Close() {
	tcpServer._cancel()
}
