package dispatcher

import "github.com/sf-go/gsf/src/gsm/peer"

type RawDispatcher struct {
}

func NewRawDispatcher() *RawDispatcher {
	return &RawDispatcher{}
}

func (dispatcher *RawDispatcher) Dispatch(peer peer.IPeer, data []byte) {
	panic("implement me")
}

func (dispatcher *RawDispatcher) Register(messageId []byte, handle func() interface{}, before func() interface{}, After func() interface{}) {
	panic("implement me")
}

func (dispatcher *RawDispatcher) Invoke(messageId []byte, peer peer.IPeer, function func() []interface{}) []interface{} {
	panic("implement me")
}

func (dispatcher *RawDispatcher) AsyncInvoke(messageId []byte, peer peer.IPeer, args func() []interface{}, result func([]interface{})) {
	panic("implement me")
}
