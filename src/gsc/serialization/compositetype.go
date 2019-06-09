package serialization

//  1位(指针) 1位(符号) 4位 (Value) | 1位(指针) 1位(符号) 4位 (Key类型) | 1位(指针)  1位(符号) 4位(类型)
type CompositeType struct {
	buffer []uint8
}

func NewCompositeType(buffer []uint8) *CompositeType {
	return &CompositeType{
		buffer: buffer,
	}
}

func (compositeType *CompositeType) Append(value uint8) {
	compositeType.buffer = append(compositeType.buffer, value)
}

func (compositeType *CompositeType) Pop() uint8 {
	if len(compositeType.buffer) <= 0 {
		panic("Empty!")
	}

	length := len(compositeType.buffer)
	defer func() {
		compositeType.buffer = compositeType.buffer[:length-1]
	}()
	return compositeType.buffer[length-1]
}

func (compositeType *CompositeType) GetType() []uint8 {
	return compositeType.buffer
}
