package network

type NetConfig struct {
	Enable         bool
	BufferSize     int32
	Address        string
	Port           int32
	ConnectTimeout int32
}

func NewNetConfig() *NetConfig {
	return &NetConfig{}
}
