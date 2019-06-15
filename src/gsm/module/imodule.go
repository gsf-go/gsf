package module

import (
	"github.com/sf-go/gsf/src/gsf/peer"
	"github.com/sf-go/gsf/src/gsf/service"
)

type IModule interface {
	Initialize(service service.IService)
	InitializeFinish(service service.IService)

	Connected(peer peer.IPeer)
	Disconnected(peer peer.IPeer)
}
