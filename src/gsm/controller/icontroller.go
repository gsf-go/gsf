package controller

import (
	"github.com/sf-go/gsf/src/gsf/peer"
)

type IController interface {
	Initialize()
	Register(messageId string, handle func() interface{}, before func() interface{}, After func() interface{})

	Invoke(messageId string, peer peer.IPeer, function func() []interface{}) []interface{}
	AsyncInvoke(messageId string, peer peer.IPeer, args func() []interface{}, result func([]interface{}))
}
