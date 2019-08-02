package modules

import (
	"github.com/sf-go/gsf/src/example/client/components"
	"github.com/sf-go/gsf/src/example/client/controllers"
	"github.com/sf-go/gsf/src/example/client/models"
	"github.com/sf-go/gsf/src/gsc/logger"
	"github.com/sf-go/gsf/src/gsc/serialization"
	"github.com/sf-go/gsf/src/gsf/service"
	"github.com/sf-go/gsf/src/gsm/dispatcher"
	"github.com/sf-go/gsf/src/gsm/module"
	"github.com/sf-go/gsf/src/gsm/peer"
	"strconv"
)

type TestClientModule struct {
	*module.Module
	dispatcher dispatcher.IDispatcher
}

func NewTestClientModule() *TestClientModule {
	return &TestClientModule{
		Module: module.NewModule(),
	}
}

func (testModule *TestClientModule) Initialize(service service.IService) {
	testModule.Module.Initialize(service)

	testModule.dispatcher = service.GetDispatcher()
	testModule.AddController(controllers.NewTestController(testModule.dispatcher))
	testModule.AddModel("TestModel", func(name string, args ...interface{}) serialization.ISerializablePacket {
		return models.NewTestModel()
	})
	logger.Log.Debug("Initialize")
}

func (testModule *TestClientModule) Connected(peer peer.IPeer) {

	component := components.NewUserComponent()
	component.SetValue("Account", "account")
	component.SetValue("Password", "123456")
	peer.AddComponent(component)

	result := testModule.dispatcher.Invoke([]byte("Test"), peer, func() []interface{} {
		return []interface{}{
			10000,
			&models.TestModel{
				Name: "wwj",
				Age:  500,
			},
		}
	})

	logger.Log.Debug(strconv.Itoa(result[0].(int)))

	testModule.dispatcher.AsyncInvoke([]byte("Test"), peer, func() []interface{} {
		return []interface{}{
			10000,
			&models.TestModel{
				Name: "wwj",
				Age:  500,
			},
		}
	}, func(result []interface{}) {
		logger.Log.Debug(strconv.Itoa(result[0].(int)))
	})
}

func (testModule *TestClientModule) InitializeFinish(service service.IService) {
	testModule.Module.InitializeFinish(service)

	logger.Log.Debug("InitializeFinish")
}
