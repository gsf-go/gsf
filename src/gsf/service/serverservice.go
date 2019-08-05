package service

import (
	"github.com/sf-go/gsf/src/gsc/network"
	"github.com/sf-go/gsf/src/gsf/socket"
	"github.com/sf-go/gsf/src/gsm/invoker"
	"github.com/sf-go/gsf/src/gsm/peer"
	"sync"
)

type ServerService struct {
	serverSocket    socket.IServerSocket
	event           *sync.Map
	messageCallback []func(peer peer.IPeer, data []byte) bool
	dispatcher      invoker.IInvoker
}

func NewServerService(dispatcher invoker.IInvoker, serverSocket socket.IServerSocket) *ServerService {
	return &ServerService{
		serverSocket: serverSocket,
		event:        new(sync.Map),
		dispatcher:   dispatcher,
	}
}

func (service *ServerService) GetInvoker() invoker.IInvoker {
	return service.dispatcher
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
	callback func(peer peer.IPeer, data []byte) bool) {

	service.messageCallback = append(service.messageCallback, callback)
	serverSocket, ok := service.serverSocket.(*socket.ServerSocket)
	if ok {
		serverSocket.OnMessage = func(peer peer.IPeer, data []byte) {
			for _, item := range service.messageCallback {
				if item(peer, data) {
					return
				}
			}
		}
	}
}
