package client

import (
	"github.com/sf-go/gsf/src/gsc/network"
)

type IClient interface {
	Connect(config *network.NetConfig)
}
