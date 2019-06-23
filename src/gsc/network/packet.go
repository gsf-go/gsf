package network

type Packet struct {
	Config     *NetConfig
	Buffer     []byte
	Offset     uint16
	Connection IConnection
}
