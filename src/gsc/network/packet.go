package network

type Packet struct {
	Config *NetConfig
	Buffer []byte
}
