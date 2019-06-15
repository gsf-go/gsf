package client

import (
	"context"
	"github.com/sf-go/gsf/src/gsc/logger"
	"github.com/sf-go/gsf/src/gsc/network"
	"net"
	"strconv"
	"time"
)

type udpClient struct {
	*network.NetEvent

	context context.Context
	handle  network.IHandle

	_cancel func()
	_close  func()
}

func NewUdpClient(
	handle network.IHandle) *udpClient {

	ctx, cancel := context.WithCancel(context.Background())
	return &udpClient{
		context:  ctx,
		handle:   handle,
		_cancel:  cancel,
		NetEvent: network.NewNetEvent(),
	}
}

func (udpClient *udpClient) Connect(config *network.NetConfig) {

	address := config.Address + ":" + strconv.Itoa(int(config.Port))
	conn, err := net.DialTimeout(
		"udp",
		address,
		time.Second*time.Duration(config.ConnectTimeout))

	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	errDone := make(chan func() (network.IConnection, error, string), 1)

	connection := network.NewConnection(
		conn.Write,
		udpClient.handle,
		errDone)

	if err, ok := err.(net.Error); ok && err.Timeout() {
		if udpClient.OnConnected != nil {
			udpClient.OnConnected(nil)
		}
		return
	}

	if udpClient.OnConnected != nil {
		udpClient.OnConnected(connection)
	}

	udpClient._close = func() {
		if conn != nil {
			_ = conn.Close()
		}

		if udpClient.OnDisconnected != nil {
			udpClient.OnDisconnected(connection, "close")
		}
	}

	go func() {
		udpClient.read(
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

			if err != nil {
				udpClient.OnDisconnected(connection, reason)
			} else {
				udpClient.OnDisconnected(connection, "Error")
			}

			if udpClient.OnError != nil {
				udpClient.OnError(connection, err)
			}

			udpClient.close(reason)
		}
	}()
}

func (udpClient *udpClient) read(
	config *network.NetConfig,
	conn net.Conn,
	buffer []byte,
	errChan chan<- func() (network.IConnection, error, string),
	connection network.IConnection) {

	offset := uint16(0)

	for {
		select {
		case <-udpClient.context.Done():
			return
		default:
			err := conn.SetReadDeadline(time.Now().Add(15 * time.Second))
			if err != nil {
				errChan <- func() (network.IConnection, error, string) {
					return connection, err, ""
				}
				return
			}

			n, err := conn.Read(buffer)
			if err != nil {
				errChan <- func() (network.IConnection, error, string) {
					return connection, err, ""
				}
				return
			}

			if n == 0 {
				errChan <- func() (network.IConnection, error, string) {
					return connection, nil, "EOF"
				}
				return
			}

			offset += udpClient.handleData(
				config,
				connection,
				buffer[0:uint16(n)+offset])
		}
	}

}

func (udpClient *udpClient) close(reason string) {
	udpClient._cancel()
	if udpClient._close != nil {
		udpClient._close()
	}

	logger.Log.Debug("断开连接！")
	udpClient.OnMessage = nil
	udpClient.OnConnected = nil
	udpClient.OnDisconnected = nil
	udpClient.OnError = nil
}

func (udpClient *udpClient) handleData(
	config *network.NetConfig,
	connection network.IConnection,
	buffer []byte) uint16 {

	packet := &network.Packet{
		Config:     config,
		Buffer:     buffer,
		Connection: connection,
	}

	return udpClient.handle.ReadHandle(packet, udpClient.post)
}

func (udpClient *udpClient) post(connection network.IConnection, data []byte) {
	if udpClient.OnMessage != nil {
		udpClient.OnMessage(connection, data)
	}
}
