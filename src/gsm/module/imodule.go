package module

import (
	"gsf/peer"
	"gsf/service"
)

type IModule interface {
	Initialize(service service.IService)
	InitializeFinish(service service.IService)

	Connected(peer peer.IPeer)
	Disconnected(peer peer.IPeer)
}
