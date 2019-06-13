package peer

import (
	"github.com/gsf/gsf/src/gsc/network"
	"github.com/gsf/gsf/src/gsm/component"
)

type IPeer interface {
	GetConnection() network.IConnection
	SetConnection(connection network.IConnection)

	AddComponent(componentName string, component component.IComponent)
	GetComponent(componentName string) component.IComponent
	RemoveComponent(componentName string)
	HasComponent(componentName string) bool
}
