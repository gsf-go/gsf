package network

import (
	"github.com/sf-go/gsf/src/gsc/bytestream"
)

type HandleCallback func(connection IConnection, data []byte)

type IHandle interface {
	ReadHandle(packet *Packet, handleCallback HandleCallback)
	WriteHandle(data []byte) []byte
}

type Handle struct {
}

func NewHandle() *Handle {
	return &Handle{}
}

func (handle *Handle) ReadHandle(
	packet *Packet,
	handleCallback HandleCallback) {

	verifyLength := uint16(2)

	if packet.Offset > verifyLength {
		byteArray := bytestream.NewByteReader2(packet.Buffer)
		length := uint16(0)
		byteArray.Read(&length)
		total := length + verifyLength

		if total <= packet.Offset {
			if handleCallback != nil {
				handleCallback(packet.Connection, packet.Buffer[verifyLength:])
			}

			copy(packet.Buffer[0:packet.Offset-total], packet.Buffer[total:packet.Offset])
			packet.Offset = packet.Offset - total
			if packet.Offset > 0 {
				handle.ReadHandle(packet, handleCallback)
			}
		}
	}
}

func (handle *Handle) WriteHandle(data []byte) []byte {
	byteArray := bytestream.NewByteWriter3()
	byteArray.Write(data)
	return byteArray.ToBytes()
}
