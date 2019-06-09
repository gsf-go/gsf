package network

import (
	"sync"
)

type ConnectionPool struct {
	pool *sync.Map
}

func NewConnectionPool() *ConnectionPool {
	return &ConnectionPool{
		pool: new(sync.Map),
	}
}

func (connectionPool *ConnectionPool) AddConnection(address string, conn IConnection) {
	connectionPool.pool.Store(address, conn)
}

func (connectionPool *ConnectionPool) GetConnection(address string) IConnection {
	v, ok := connectionPool.pool.Load(address)
	if ok {
		return v.(IConnection)
	}
	return nil
}

func (connectionPool *ConnectionPool) Remove(address string) {
	connectionPool.pool.Delete(address)
}
