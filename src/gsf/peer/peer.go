package peer

import (
	"github.com/gsf/gsf/src/gsc/network"
	"github.com/gsf/gsf/src/gsm/component"
	"sync"
)

type Peer struct {
	connection network.IConnection
	components *sync.Map
}

func (peer *Peer) Range(foreach func(key string, component component.IComponent) bool) {
	peer.components.Range(func(key, value interface{}) bool {
		return foreach(key.(string), value.(component.IComponent))
	})
}

func NewPeer() *Peer {
	return &Peer{
		components: new(sync.Map),
	}
}

func (peer *Peer) GetConnection() network.IConnection {
	return peer.connection
}

func (peer *Peer) SetConnection(connection network.IConnection) {
	peer.connection = connection
}

func (peer *Peer) AddComponent(component component.IComponent) {
	peer.components.Store(component.GetName(), component)
}

func (peer *Peer) GetComponent(componentName string) component.IComponent {
	cpt, ok := peer.components.Load(componentName)
	if ok {
		return cpt.(component.IComponent)
	}
	return nil
}

func (peer *Peer) RemoveComponent(componentName string) {
	peer.components.Delete(componentName)
}

func (peer *Peer) HasComponent(componentName string) bool {
	_, ok := peer.components.Load(componentName)
	return ok
}
