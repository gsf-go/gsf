package service

import (
	"gsc/network"
	"gsf/peer"
	"gsf/socket"
	"sync"
)

type ClientService struct {
	clientSocket socket.IClientSocket
	event        *sync.Map
}

func NewClientService(clientSocket socket.IClientSocket) *ClientService {
	return &ClientService{
		clientSocket: clientSocket,
		event:        new(sync.Map),
	}
}

func (service *ClientService) Connect(config *network.NetConfig) {
	service.clientSocket.Connect(config)
}

func (service *ClientService) AddEventListener(
	eventType string,
	callback func(peer peer.IPeer, args ...interface{})) {
	service.event.Store(eventType, callback)
}

func (service *ClientService) PostNotification(
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

func (service *ClientService) RemoveEventListener(eventType string) {
	service.event.Delete(eventType)
}

func (service *ClientService) SetHandler(
	opCode string,
	callback func(peer peer.IPeer, data []byte)) {

}
