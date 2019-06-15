package service

import (
	"github.com/sf-go/gsf/src/gsc/network"
	"github.com/sf-go/gsf/src/gsf/peer"
	"github.com/sf-go/gsf/src/gsf/socket"
	"sync"
)

type ServerService struct {
	serverSocket socket.IServerSocket
	event        *sync.Map
}

func NewServerService(serverSocket socket.IServerSocket) *ServerService {
	return &ServerService{
		serverSocket: serverSocket,
		event:        new(sync.Map),
	}
}

func (service *ServerService) StartServer(config *network.NetConfig) {
	service.serverSocket.Start(config)
}

func (service *ServerService) StopServer() {
	service.serverSocket.Stop()
}

func (service *ServerService) AddEventListener(
	eventType string,
	callback func(peer peer.IPeer, args ...interface{})) {
	service.event.Store(eventType, callback)
}

func (service *ServerService) PostNotification(
	eventType string,
	p peer.IPeer,
	args ...interface{}) {
	evt, ok := service.event.Load(eventType)
	if ok {
		callback := evt.(func(peer.IPeer, ...interface{}))
		if callback != nil {
			callback(p, args...)
		}
	}
}

func (service *ServerService) RemoveEventListener(eventType string) {
	service.event.Delete(eventType)
}

func (service *ServerService) SetHandler(
	opCode string,
	callback func(peer peer.IPeer, data []byte)) {

}
