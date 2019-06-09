package pool

import (
	"sync"
)

type ObjectPool struct {
	pool *sync.Pool
}

func NewObjectPool(generate func() interface{}) *ObjectPool {
	return &ObjectPool{
		pool: &sync.Pool{
			New: generate,
		},
	}
}

func (objectPool *ObjectPool) GetObject() interface{} {
	return objectPool.pool.Get()
}

func (objectPool *ObjectPool) Recycle(obj interface{}) {
	objectPool.pool.Put(obj)
}
