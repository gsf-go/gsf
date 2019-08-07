package peer

import (
	"github.com/sf-go/gsf/src/gsc/network"
	"github.com/sf-go/gsf/src/gsm/component"
)

type Peer struct {
	connection network.IConnection
	components map[string]component.IComponent
}

func (peer *Peer) Range(foreach func(key string, component component.IComponent) bool) {
	for key, value := range peer.components {
		if !foreach(key, value) {
			break
		}
	}
}

func NewPeer() *Peer {
	return &Peer{
		components: make(map[string]component.IComponent),
	}
}

func (peer *Peer) GetConnection() network.IConnection {
	return peer.connection
}

func (peer *Peer) SetConnection(connection network.IConnection) {
	peer.connection = connection
}

func (peer *Peer) AddComponent(component component.IComponent) {
	peer.components[component.GetObjectId()] = component
}

func (peer *Peer) GetComponent(componentName string) component.IComponent {
	cpt, ok := peer.components[componentName]
	if ok {
		return cpt
	}
	return nil
}

func (peer *Peer) RemoveComponent(componentName string) {
	delete(peer.components, componentName)
}

func (peer *Peer) HasComponent(componentName string) bool {
	_, ok := peer.components[componentName]
	return ok
}
