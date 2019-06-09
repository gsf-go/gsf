package socket

import "gsc/network"

type IClientSocket interface {
	Connect(config *network.NetConfig)
}
