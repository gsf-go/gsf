package network

import (
	"gsc/bytestream"
)

type HandleCallback func(connection IConnection, data []byte)

type IHandle interface {
	ReadHandle(packet *Packet, connection IConnection, handleCallback HandleCallback) uint16
	WriteHandle(data []byte) []byte
}

type Handle struct {
}

func NewHandle() *Handle {
	return &Handle{}
}

func (handle *Handle) ReadHandle(
	packet *Packet,
	connection IConnection,
	handleCallback HandleCallback) uint16 {

	count := uint16(len(packet.Buffer))
	verifyLength := uint16(2)

	if count > verifyLength {
		byteArray := bytestream.NewByteReader2(packet.Buffer)
		length := uint16(0)
		byteArray.Read(&length)

		if length <= count {
			if handleCallback != nil {
				handleCallback(connection, packet.Buffer[verifyLength:])
			}

			copy(packet.Buffer[0:count-length], packet.Buffer[length:count])
			packet.Buffer = packet.Buffer[0 : count-length]
			count = handle.ReadHandle(packet, connection, handleCallback)
		}
	}
	return count
}

func (handle *Handle) WriteHandle(data []byte) []byte {
	byteArray := bytestream.NewByteWriter3()
	byteArray.Write(data)
	return byteArray.ToBytes()
}
