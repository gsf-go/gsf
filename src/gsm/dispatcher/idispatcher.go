package dispatcher

import "github.com/sf-go/gsf/src/gsm/peer"

type IDispatcher interface {
	Dispatch(peer peer.IPeer, data []byte)

	Register(id interface{},
		before func() interface{},
		handle func() interface{},
		After func() interface{})

	RawRegister(id interface{},
		before func(peer peer.IPeer, data []byte) bool,
		handle func(peer peer.IPeer, data []byte) []byte,
		After func(peer peer.IPeer, data []byte))

	Invoke(id interface{},
		peer peer.IPeer,
		function func() []interface{}) []interface{}

	RawInvoke(id interface{},
		p peer.IPeer,
		data []byte) []byte

	AsyncInvoke(id interface{},
		peer peer.IPeer,
		args func() []interface{},
		result func([]interface{}))

	AsyncRawInvoke(id interface{},
		p peer.IPeer,
		data []byte,
		result func(data []byte))
}
