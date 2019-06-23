package pool

import (
	"bytes"
	"sync"
)

type BytePool struct {
	pool *sync.Pool
}

func NewBytePool(bufferSize int32) *BytePool {
	return &BytePool{
		pool: &sync.Pool{
			New: func() interface{} {
				buffer := make([]byte, bufferSize)
				return bytes.NewBuffer(buffer)
			},
		},
	}
}

func (bytePool *BytePool) GetBuffer() *bytes.Buffer {
	return bytePool.pool.Get().(*bytes.Buffer)
}

func (bytePool *BytePool) Recycle(obj interface{}) {
	//buffer := obj.(*bytes.Buffer)
	//buffer.Reset()
	bytePool.pool.Put(obj)
}
