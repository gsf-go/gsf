package server

import (
	"context"
	"gsc/logger"
	"gsc/network"
	"gsc/pool"
	"net"
	"strconv"
	"time"
)

type udpServer struct {
	*network.NetEvent

	handle         network.IHandle
	context        context.Context
	connectionPool *network.ConnectionPool

	_cancel func()
}

func NewUdpServer(
	handle network.IHandle) *udpServer {
	ctx, cancel := context.WithCancel(context.Background())
	return &udpServer{
		handle:         handle,
		_cancel:        cancel,
		context:        ctx,
		connectionPool: network.NewConnectionPool(),
		NetEvent:       network.NewNetEvent(),
	}
}

func (udpServer *udpServer) Accept(config *network.NetConfig) {

	bufferPool := pool.NewBytePool(config.BufferSize)
	address := config.Address + ":" + strconv.Itoa(int(config.Port))
	packetConn, err := net.ListenPacket("udp", address)

	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	errDone := make(chan func() (network.IConnection, error, string), 1)

	go udpServer.handleClient(config, bufferPool, errDone, packetConn)

	go func() {
		select {
		case callback := <-errDone:
			connection, err, reason := callback()

			if err != nil {
				udpServer.OnDisconnected(connection, reason)
			} else {
				udpServer.OnDisconnected(connection, "Error")
			}

			if udpServer.OnError != nil {
				udpServer.OnError(connection, err)
			}
		}
	}()
}

func (udpServer *udpServer) handleClient(
	config *network.NetConfig,
	bufferPool *pool.BytePool,
	errChan chan<- func() (network.IConnection, error, string),
	packetConn net.PacketConn) {

	for {
		select {
		case <-udpServer.context.Done():
			return
		default:
			buffer := bufferPool.GetBuffer()

			byteArray := buffer.Bytes()
			offset := uint16(0)

			n, addr, err := packetConn.ReadFrom(byteArray[offset:])
			connection := udpServer.connectionPool.GetConnection(addr.String())
			if connection == nil {
				connection = network.NewConnection(
					func(data []byte) (i int, e error) {
						return packetConn.WriteTo(data, addr)
					},
					udpServer.handle,
					errChan)

				if udpServer.OnConnected != nil {
					udpServer.OnConnected(connection)
				}

				udpServer.connectionPool.AddConnection(addr.String(), connection)
			}

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

			deadline := time.Now().Add(15 * time.Second)
			err = packetConn.SetWriteDeadline(deadline)
			if err != nil {
				errChan <- func() (network.IConnection, error, string) {
					return connection, err, ""
				}
				return
			}

			offset += udpServer.handleData(config, connection,
				byteArray[0:uint16(n)+offset])
		}
	}
}

func (udpServer *udpServer) handleData(
	config *network.NetConfig,
	connection network.IConnection,
	buffer []byte) uint16 {

	packet := &network.Packet{
		Config: config,
		Buffer: buffer,
	}

	return udpServer.handle.ReadHandle(packet, connection, udpServer.post)
}

func (udpServer *udpServer) post(connection network.IConnection, data []byte) {
	if udpServer.OnMessage != nil {
		udpServer.OnMessage(connection, data)
	}
}

func (udpServer *udpServer) Close(reason string) {
	udpServer._cancel()
}
