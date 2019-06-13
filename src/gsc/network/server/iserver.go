package server

import "github.com/gsf/gsf/src/gsc/network"

type IServer interface {
	Accept(config *network.NetConfig)
	Close()
}
