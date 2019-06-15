package server

import "github.com/sf-go/gsf/src/gsc/network"

type IServer interface {
	Accept(config *network.NetConfig)
	Close()
}
