package controller

import (
	"gsf/peer"
)

type IController interface {
	Initialize()

	Register(messageId string, function func() interface{})
	Invoke(messageId string, peer peer.IPeer, function func() []interface{})
}
