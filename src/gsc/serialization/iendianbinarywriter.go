package serialization

type IEndianBinaryWriter interface {
	Write(args ...interface{})
	ToBytes() []byte
}

type EndianBinaryWriter struct {
	writer ISerializable
	buffer []byte
}

func NewEndianBinaryWriter() *EndianBinaryWriter {
	return &EndianBinaryWriter{
		writer: NewSerializable(),
	}
}

func (w *EndianBinaryWriter) Write(args ...interface{}) {
	w.buffer = w.writer.Serialize(args...)
}

func (w *EndianBinaryWriter) ToBytes() []byte {
	return w.buffer
}
