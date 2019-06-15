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
	testModule.AddModel("TestModel", func() serialization.ISerializablePacket {
		return new(models.TestModel)
	})
	logger.Log.Debug("Initialize")
}

func (testModule *TestServerModule) InitializeFinish(service service.IService) {
	testModule.Module.InitializeFinish(service)

	logger.Log.Debug("InitializeFinish")
}

func (testModule *TestServerModule) Connected(peer peer.IPeer) {
	logger.Log.Debug("Connected")
	peer.AddComponent("User", components.NewUserComponent())
}
