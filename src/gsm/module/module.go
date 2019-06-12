package module

import (
	"gsc/serialization"
	"gsf/peer"
	"gsf/service"
	"gsm/controller"
)

type Module struct {
	controllers []controller.IController
}

func NewModule() *Module {
	return &Module{
		controllers: make([]controller.IController, 0),
	}
}

func (module *Module) AddController(controller controller.IController) {
	module.controllers = append(module.controllers, controller)
	controller.Initialize()
}

func (module *Module) AddModel(name string, generate func() serialization.ISerializablePacket) {
	serialization.GetPacketManagerInstance().AddPacket(name, generate)
}

func (module *Module) Initialize(service service.IService) {

}

func (module *Module) InitializeFinish(service service.IService) {

}

func (module *Module) Accepted(peer peer.IPeer) {

}

func (module *Module) Connected(peer peer.IPeer) {

}

func (module *Module) Disconnected(peer peer.IPeer) {

}
