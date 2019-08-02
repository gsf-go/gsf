package dispatcher

import "github.com/sf-go/gsf/src/gsm/peer"

type IDispatcher interface {
	Dispatch(peer peer.IPeer, data []byte)

	Register(messageId []byte,
		before func() interface{},
		handle func() interface{},
		After func() interface{})

	RawRegister(messageId []byte,
		before func(peer peer.IPeer, data []byte) bool,
		handle func(peer peer.IPeer, data []byte) []interface{},
		After func(peer peer.IPeer, data []byte))

	Invoke(messageId []byte,
		peer peer.IPeer,
		function func() []interface{}) []interface{}

	AsyncInvoke(messageId []byte,
		peer peer.IPeer,
		args func() []interface{},
		result func([]interface{}))
}
