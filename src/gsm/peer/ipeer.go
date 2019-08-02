package peer

import (
	"github.com/sf-go/gsf/src/gsc/network"
	"github.com/sf-go/gsf/src/gsm/component"
)

type IPeer interface {
	GetConnection() network.IConnection
	SetConnection(connection network.IConnection)

	AddComponent(component component.IComponent)
	GetComponent(componentName string) component.IComponent
	RemoveComponent(componentName string)
	HasComponent(componentName string) bool

	Range(func(key string, component component.IComponent) bool)
}
