package service

import (
	"gsf/peer"
)

type IService interface {
	AddEventListener(eventType string, callback func(peer peer.IPeer, args ...interface{}))
	PostNotification(eventType string, peer peer.IPeer, args ...interface{})
	RemoveEventListener(eventType string)
	SetHandler(opCode string, callback func(peer peer.IPeer, data []byte))
}
