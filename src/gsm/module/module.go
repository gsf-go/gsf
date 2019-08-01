package module

import (
	"github.com/sf-go/gsf/src/gsc/serialization"
	"github.com/sf-go/gsf/src/gsf/peer"
	"github.com/sf-go/gsf/src/gsf/service"
	"github.com/sf-go/gsf/src/gsm/controller"
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

func (module *Module) AddModel(name string, generate func(args ...interface{}) serialization.ISerializablePacket) {
	serialization.PacketManagerInstance.AddPacket(name, generate)
}

func (module *Module) Initialize(service service.IService) {

}

func (module *Module) InitializeFinish(service service.IService) {

}

func (module *Module) Connected(peer peer.IPeer) {

}

func (module *Module) Disconnected(peer peer.IPeer) {

}
