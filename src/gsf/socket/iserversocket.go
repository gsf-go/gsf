package socket

import (
	"github.com/gsf/gsf/src/gsc/network"
	"github.com/gsf/gsf/src/gsf/peer"
)

type Event struct {
	OnConnected    func(peer peer.IPeer)
	OnDisconnected func(peer peer.IPeer)
	OnMessage      func(peer peer.IPeer, data []byte)
}

func NewEvent() *Event {
	return &Event{}
}

type IServerSocket interface {
	Start(config *network.NetConfig)
	Stop()
}
