package client

import (
	"github.com/gsf/gsf/src/gsc/network"
)

type IClient interface {
	Connect(config *network.NetConfig)
}
