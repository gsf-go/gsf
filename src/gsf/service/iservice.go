package service

import (
	"github.com/sf-go/gsf/src/gsm/dispatcher"
	"github.com/sf-go/gsf/src/gsm/peer"
)

type IService interface {
	GetDispatcher() dispatcher.IDispatcher
	AddEventListener(eventType string, callback func(peer peer.IPeer, args ...interface{}))
	PostNotification(eventType string, peer peer.IPeer, args ...interface{})
	RemoveEventListener(eventType string)
	SetHandler(callback func(peer peer.IPeer, data []byte) bool)
}
