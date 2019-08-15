package module

import (
	"github.com/sf-go/gsf/src/gsc/serialization"
	"github.com/sf-go/gsf/src/gsf/service"
	"github.com/sf-go/gsf/src/gsm/component"
	"github.com/sf-go/gsf/src/gsm/controller"
	"github.com/sf-go/gsf/src/gsm/invoker"
	"github.com/sf-go/gsf/src/gsm/peer"
)

type Module struct {
	controllers []controller.IController
	invoker     invoker.IInvoker
}

func NewModule() *Module {
	return &Module{
		controllers: make([]controller.IController, 0),
	}
}

func (module *Module) GetInvoker() invoker.IInvoker {
	return module.invoker
}

func (module *Module) AddController(controller controller.IController) {
	module.controllers = append(module.controllers, controller)
}

func (module *Module) AddModel(name string,
	generate func(name string, peer peer.IPeer) serialization.ISerializablePacket) {
	serialization.PacketManagerInstance.AddPacket(name,
		func(name string, args ...interface{}) serialization.ISerializablePacket {
			return generate(name, args[0].(peer.IPeer))
		})
}

func (module *Module) AddComponent(template component.IComponent) {

	name := template.GetObjectId()
	module.invoker.FixRegister("Get_"+name,
		func(peer peer.IPeer, args ...interface{}) []interface{} {
			return peer.GetComponent(name).GetterCallback(args[0].(string))
		})

	module.invoker.FixRegister("Set_"+name,
		func(peer peer.IPeer, args ...interface{}) []interface{} {
			return []interface{}{
				peer.GetComponent(name).SetterCallback(args...),
			}
		})
}

func (module *Module) Initialize(service service.IService) {
	module.invoker = service.GetInvoker()
}

func (module *Module) InitializeFinish(service service.IService) {

}

func (module *Module) Connected(peer peer.IPeer) {

}

func (module *Module) Disconnected(peer peer.IPeer) {

}
