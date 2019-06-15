package network

import (
	"github.com/sf-go/gsf/src/gsc/logger"
)

type IConnection interface {
	Send(buffer []byte)
	Close()
}

type Connection struct {
	write        func([]byte) (int, error)
	handle       IHandle
	closeDone    chan<- func() (IConnection, error, string)
	messageQueue *messageQueue
}

func (connection *Connection) Close() {
	connection.messageQueue.Close()
	connection.closeDone <- func() (IConnection, error, string) {
		return connection, nil, "close"
	}
}

func NewConnection(
	write func([]byte) (int, error),
	handle IHandle,
	closeDone chan<- func() (IConnection, error, string)) *Connection {
	return &Connection{
		write:        write,
		handle:       handle,
		closeDone:    closeDone,
		messageQueue: newMessageQueue(),
	}
}

func (connection *Connection) Send(buffer []byte) {
	if connection.handle == nil {
		return
	}

	connection.messageQueue.GetWriteChan() <- func() {

		data := connection.handle.WriteHandle(buffer)
		bytesSender, err := connection.write(data)
		if err != nil {
			connection.closeDone <- func() (IConnection, error, string) {
				return connection, err, ""
			}
		}

		logger.Log.Debug("Total:%d Send:%d", len(buffer), bytesSender)
	}
}
