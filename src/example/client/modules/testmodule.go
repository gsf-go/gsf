package modules

import (
	"example/client/controllers"
	"gsc/logger"
	"gsf/peer"
	"gsf/service"
	"gsm/module"
)

type TestModule struct {
	*module.Module
}

func NewTestModule() *TestModule {
	return &TestModule{
		Module: module.NewModule(),
	}
}

func (testModule *TestModule) Initialize(service service.IService) {
	testModule.Module.Initialize(service)

	testModule.AddController(controllers.NewTestController())
	logger.Log.Debug("Initialize")
}

func (testModule *TestModule) Connected(peer peer.IPeer) {
	controller := controllers.NewTestController()
	controller.Invoke("Test", peer, func() []interface{} {
		return []interface{}{10000}
	})
}

func (testModule *TestModule) InitializeFinish(service service.IService) {
	testModule.Module.InitializeFinish(service)

	logger.Log.Debug("InitializeFinish")
}
