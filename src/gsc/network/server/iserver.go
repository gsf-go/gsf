package server

import "gsc/network"

type IServer interface {
	Accept(config *network.NetConfig)
	Close()
}
