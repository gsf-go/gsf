package modules

import (
	"github.com/sf-go/gsf/src/example/server/components"
	"github.com/sf-go/gsf/src/example/server/controllers"
	"github.com/sf-go/gsf/src/example/server/models"
	"github.com/sf-go/gsf/src/gsc/logger"
	"github.com/sf-go/gsf/src/gsc/serialization"
	"github.com/sf-go/gsf/src/gsf/service"
	"github.com/sf-go/gsf/src/gsm/module"
	"github.com/sf-go/gsf/src/gsm/peer"
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
	testModule.AddModel("TestModel", func(name string, args ...interface{}) serialization.ISerializablePacket {
		return models.NewTestModel()
	})

	testModule.AddModel("UserComponent", func(name string, args ...interface{}) serialization.ISerializablePacket {
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

func (testModule *TestServerModule) Disconnected(peer peer.IPeer) {
	logger.Log.Debug("Disconnected")
}
