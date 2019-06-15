package modules

import (
	"github.com/gsf/gsf/src/example/server/components"
	"github.com/gsf/gsf/src/example/server/controllers"
	"github.com/gsf/gsf/src/example/server/models"
	"github.com/gsf/gsf/src/gsc/logger"
	"github.com/gsf/gsf/src/gsc/serialization"
	"github.com/gsf/gsf/src/gsf/peer"
	"github.com/gsf/gsf/src/gsf/service"
	"github.com/gsf/gsf/src/gsm/module"
)

type TestServerModule struct {
	*module.Module
}

func NewTestServerModule() *TestServerModule {
	return &TestServerModule{
		Module: module.NewModule(),
	}
}

func (testModule *TestServerModule) Initialize(service service.IService) {
	testModule.Module.Initialize(service)

	testModule.AddController(controllers.NewTestController())
	testModule.AddModel("TestModel", func(args ...interface{}) serialization.ISerializablePacket {
		return models.NewTestModel()
	})

	testModule.AddModel("UserComponent", func(args ...interface{}) serialization.ISerializablePacket {
		name := args[0].(string)
		p := args[1].(peer.IPeer)
		return p.GetComponent(name).(serialization.ISerializablePacket)
	})

	logger.Log.Debug("Initialize")
}

func (testModule *TestServerModule) InitializeFinish(service service.IService) {
	testModule.Module.InitializeFinish(service)

	logger.Log.Debug("InitializeFinish")
}

func (testModule *TestServerModule) Connected(peer peer.IPeer) {
	logger.Log.Debug("Connected")
	peer.AddComponent(components.NewUserComponent())
}
