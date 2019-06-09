package socket

import (
	"gsc/network"
	"gsf/peer"
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
