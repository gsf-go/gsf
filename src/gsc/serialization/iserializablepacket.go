package serialization

type ISerializablePacket interface {
	ToBinaryWriter(writer ISerializable) []byte
	FromBinaryReader(reader IDeserializable)
}
