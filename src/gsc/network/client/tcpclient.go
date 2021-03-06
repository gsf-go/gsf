package client

import (
	"context"
	"github.com/sf-go/gsf/src/gsc/logger"
	"github.com/sf-go/gsf/src/gsc/network"
	"net"
	"strconv"
	"time"
)

type TcpClient struct {
	*network.NetEvent

	context context.Context
	handle  network.IHandle

	_close  func()
	_cancel func()
}

func NewTcpClient(
	handle network.IHandle) *TcpClient {

	ctx, cancel := context.WithCancel(context.Background())
	return &TcpClient{
		context:  ctx,
		handle:   handle,
		_cancel:  cancel,
		NetEvent: network.NewNetEvent(),
	}
}

func (tcpClient *TcpClient) Connect(config *network.NetConfig) {

	address := config.Address + ":" + strconv.Itoa(int(config.Port))
	conn, err := net.DialTimeout("tcp", address,
		time.Second*time.Duration(config.ConnectTimeout))

	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	errDone := make(chan func() (network.IConnection, error, string), 1)

	connection := network.NewConnection(
		conn.Write,
		tcpClient.handle,
		errDone)

	if err, ok := err.(net.Error); ok && err.Timeout() {
		if tcpClient.OnConnected != nil {
			tcpClient.OnConnected(nil)
		}
		return
	}

	tcpClient._close = func() {
		if conn != nil {
			_ = conn.Close()
		}
	}

	go func() {
		tcpClient.read(
			config,
			conn,
			make([]byte, config.BufferSize),
			errDone,
			connection)
	}()

	go func() {
		select {
		case callback := <-errDone:
			connection, err, reason := callback()

			tcpClient.OnDisconnected(connection, reason)

			if tcpClient.OnError != nil {
				tcpClient.OnError(connection, err)
			}

			tcpClient.close(reason)
		}
	}()

	if tcpClient.OnConnected != nil {
		tcpClient.OnConnected(connection)
	}
}

func (tcpClient *TcpClient) read(
	config *network.NetConfig,
	conn net.Conn,
	buffer []byte,
	errChan chan<- func() (network.IConnection, error, string),
	connection network.IConnection) {

	packet := &network.Packet{
		Config:     config,
		Buffer:     buffer,
		Connection: connection,
	}

	for {
		select {
		case <-tcpClient.context.Done():
			return
		default:
			n, err := conn.Read(buffer[packet.Offset:])
			if err != nil {
				errChan <- func() (network.IConnection, error, string) {
					return connection, err, "Error"
				}
				return
			}

			if n == 0 {
				errChan <- func() (network.IConnection, error, string) {
					return connection, nil, "EOF"
				}
				return
			}

			deadline := time.Now().Add(15 * time.Second)
			err = conn.SetWriteDeadline(deadline)
			if err != nil {
				errChan <- func() (network.IConnection, error, string) {
					return connection, err, "Error"
				}
				return
			}

			packet.Offset += uint16(n)
			tcpClient.handle.ReadHandle(packet, tcpClient.post)
		}
	}

}

func (tcpClient *TcpClient) close(reason string) {
	tcpClient._cancel()
	if tcpClient._close != nil {
		tcpClient._close()
	}

	logger.Log.Debug(reason)
	tcpClient.OnMessage = nil
	tcpClient.OnConnected = nil
	tcpClient.OnDisconnected = nil
	tcpClient.OnError = nil
}

func (tcpClient *TcpClient) post(connection network.IConnection, data []byte) {
	if tcpClient.OnMessage != nil {
		tcpClient.OnMessage(connection, data)
	}
}
