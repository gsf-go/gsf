package peer

import (
	"gsc/network"
	"gsm/component"
)

type IPeer interface {
	GetConnection() network.IConnection
	SetConnection(connection network.IConnection)

	AddComponent(componentName string, component component.IComponent)
	GetComponent(componentName string) component.IComponent
	RemoveComponent(componentName string)
	HasComponent(componentName string) bool
}
