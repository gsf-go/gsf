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

			udpServer.OnDisconnected(connection, reason)

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
			packet := &network.Packet{
				Config: config,
				Buffer: byteArray,
			}
			n, addr, err := packetConn.ReadFrom(byteArray)
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

			packet.Connection = connection
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
			err = packetConn.SetWriteDeadline(deadline)
			if err != nil {
				errChan <- func() (network.IConnection, error, string) {
					return connection, err, "Error"
				}
				return
			}

			packet.Offset = uint16(n)
			go func() {
				defer func() {
					bufferPool.Recycle(buffer)
				}()

				udpServer.handle.ReadHandle(packet, udpServer.post)
			}()
		}
	}
}

func (udpServer *udpServer) post(connection network.IConnection, data []byte) {
	if udpServer.OnMessage != nil {
		udpServer.OnMessage(connection, data)
	}
}

func (udpServer *udpServer) Close(reason string) {
	udpServer._cancel()
}
