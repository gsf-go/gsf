package serialization

type ISerializablePacket interface {
	ToBinaryWriter(writer IEndianBinaryWriter)
	FromBinaryReader(reader IEndianBinaryReader)
}
