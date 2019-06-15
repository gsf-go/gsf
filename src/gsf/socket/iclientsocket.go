package socket

import "github.com/sf-go/gsf/src/gsc/network"

type IClientSocket interface {
	Connect(config *network.NetConfig)
}
