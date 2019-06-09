package network

type NetEvent struct {
	OnConnected    func(connection IConnection)
	OnMessage      func(connection IConnection, data []byte)
	OnDisconnected func(connection IConnection, reason string)
	OnError        func(connection IConnection, err error)
}

func NewNetEvent() *NetEvent {
	return &NetEvent{}
}
