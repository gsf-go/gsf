package network

type NetConfig struct {
	Enable         bool   `csv:"开关"`
	BufferSize     int32  `csv:"缓冲区尺寸"`
	Address        string `csv:"IP地址"`
	Port           int32  `csv:"IP端口"`
	ConnectTimeout int32  `csv:"连接延迟"`
}

func NewNetConfig() *NetConfig {
	return &NetConfig{}
}
