package service

import (
	"github.com/sf-go/gsf/src/gsm/invoker"
	"github.com/sf-go/gsf/src/gsm/peer"
)

type IService interface {
	GetInvoker() invoker.IInvoker
	AddEventListener(eventType string, callback func(peer peer.IPeer, args ...interface{}))
	PostNotification(eventType string, peer peer.IPeer, args ...interface{})
	RemoveEventListener(eventType string)
	SetHandler(callback func(peer peer.IPeer, data []byte) bool)
}
