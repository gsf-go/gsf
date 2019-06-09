package client

import (
	"gsc/network"
)

type IClient interface {
	Connect(config *network.NetConfig)
}
