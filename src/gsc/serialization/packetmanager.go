package serialization

import (
	"sync"
)

//go:generate ../../../gen/Singleton.exe -struct=PacketManager -out=./packetmanager.go

var PacketManagerInstance = NewPacketManager()

type PacketManager struct {
	packets *sync.Map
}

func NewPacketManager() *PacketManager {
	return &PacketManager{
		packets: new(sync.Map),
	}
}

func (packetManager *PacketManager) AddPacket(name string, generate func(name string, args ...interface{}) ISerializablePacket) {
	packetManager.packets.Store(name, generate)
}

func (packetManager *PacketManager) GetPacket(name string) func(name string, args ...interface{}) ISerializablePacket {
	v, ok := packetManager.packets.Load(name)
	if ok {
		return v.(func(string, ...interface{}) ISerializablePacket)
	}
	return nil
}

func (packetManager *PacketManager) Remove(name string) {
	packetManager.packets.Delete(name)
}
