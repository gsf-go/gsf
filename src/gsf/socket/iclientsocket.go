package socket

import "github.com/gsf/gsf/src/gsc/network"

type IClientSocket interface {
	Connect(config *network.NetConfig)
}
