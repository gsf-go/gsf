package serialization

import (
	"sync"
)

//go:generate ../../../gen/Singleton.exe -struct=PacketManager -out=./packetmanager.go

var PacketManagerInstance *PacketManager
var PacketManagerOnce sync.Once

func GetPacketManagerInstance() *PacketManager {
	PacketManagerOnce.Do(func() {
		PacketManagerInstance = NewPacketManager()
	})
	return PacketManagerInstance
}

type PacketManager struct {
	packets *sync.Map
}

func NewPacketManager() *PacketManager {
	return &PacketManager{
		packets: new(sync.Map),
	}
}

func (packetManager *PacketManager) AddPacket(name string, generate func(args ...interface{}) ISerializablePacket) {
	packetManager.packets.Store(name, generate)
}

func (packetManager *PacketManager) GetPacket(name string) func(args ...interface{}) ISerializablePacket {
	v, ok := packetManager.packets.Load(name)
	if ok {
		return v.(func(...interface{}) ISerializablePacket)
	}
	return nil
}

func (packetManager *PacketManager) Remove(name string) {
	packetManager.packets.Delete(name)
}
