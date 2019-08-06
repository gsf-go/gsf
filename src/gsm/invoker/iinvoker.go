package invoker

import (
	"github.com/sf-go/gsf/src/gsm/peer"
)

type IInvoker interface {
	Dispatch(peer peer.IPeer, data []byte)

	Register(id string,
		before func() interface{},
		handle func() interface{},
		After func() interface{})

	RawRegister(id string,
		before func(peer peer.IPeer, data []byte) bool,
		handle func(peer peer.IPeer, data []byte) []byte,
		After func(peer peer.IPeer, data []byte))

	Invoke(id string,
		peer peer.IPeer,
		function func() []interface{}) []interface{}

	RawInvoke(id string,
		p peer.IPeer,
		data []byte) []byte

	AsyncInvoke(id string,
		peer peer.IPeer,
		args func() []interface{},
		result func([]interface{}))

	AsyncRawInvoke(id string,
		p peer.IPeer,
		data []byte,
		result func(data []byte))
}
