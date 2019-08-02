package dispatcher

import "github.com/sf-go/gsf/src/gsm/peer"

type IDispatcher interface {
	Dispatch(peer peer.IPeer, data []byte)

	Register(messageId []byte, handle func() interface{}, before func() interface{}, After func() interface{})

	Invoke(messageId []byte, peer peer.IPeer, function func() []interface{}) []interface{}
	AsyncInvoke(messageId []byte, peer peer.IPeer, args func() []interface{}, result func([]interface{}))
}
