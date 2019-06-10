package controller

import (
	"gsf/peer"
)

type IController interface {
	Initialize()

	Register(messageId string, function func() interface{})
	Invoke(messageId string, peer peer.IPeer, function func() []interface{}) []interface{}
	AsyncInvoke(messageId string, peer peer.IPeer, args func() []interface{}, result func([]interface{}))
}
