package module

import (
	"github.com/gsf/gsf/src/gsf/peer"
	"github.com/gsf/gsf/src/gsf/service"
)

type IModule interface {
	Initialize(service service.IService)
	InitializeFinish(service service.IService)

	Connected(peer peer.IPeer)
	Disconnected(peer peer.IPeer)
}
