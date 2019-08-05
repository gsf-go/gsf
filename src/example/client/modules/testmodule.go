package modules

import (
	"github.com/sf-go/gsf/src/example/client/components"
	"github.com/sf-go/gsf/src/example/client/controllers"
	"github.com/sf-go/gsf/src/example/client/models"
	"github.com/sf-go/gsf/src/gsc/logger"
	"github.com/sf-go/gsf/src/gsc/serialization"
	"github.com/sf-go/gsf/src/gsf/service"
	"github.com/sf-go/gsf/src/gsm/invoker"
	"github.com/sf-go/gsf/src/gsm/module"
	"github.com/sf-go/gsf/src/gsm/peer"
)

type TestClientModule struct {
	*module.Module
	dispatcher invoker.IInvoker
}

func NewTestClientModule() *TestClientModule {
	return &TestClientModule{
		Module: module.NewModule(),
	}
}

func (testModule *TestClientModule) Initialize(service service.IService) {
	testModule.Module.Initialize(service)

	testModule.dispatcher = service.GetInvoker()
	testModule.AddController(controllers.NewTestController(testModule.dispatcher))
	testModule.AddModel("TestModel", func(name string, peer peer.IPeer) serialization.ISerializablePacket {
		return models.NewTestModel()
	})
	logger.Log.Debug("Initialize")
}

func (testModule *TestClientModule) Connected(peer peer.IPeer) {

	component := components.NewUserComponent()
	component.SetValue("Account", "account")
	component.SetValue("Password", "123456")
	peer.AddComponent(component)

	//result := testModule.dispatcher.Invoke("Test", peer, func() []interface{} {
	//	return []interface{}{
	//		10000,
	//		&models.TestModel{
	//			Name: "wwj",
	//			Age:  500,
	//		},
	//	}
	//})
	//
	//logger.Log.Debug(strconv.Itoa(result[0].(int)))
	//
	//testModule.dispatcher.AsyncInvoke("Test", peer, func() []interface{} {
	//	return []interface{}{
	//		10000,
	//		&models.TestModel{
	//			Name: "wwj",
	//			Age:  500,
	//		},
	//	}
	//}, func(result []interface{}) {
	//	logger.Log.Debug(strconv.Itoa(result[0].(int)))
	//})

	testModule.dispatcher.AsyncRawInvoke("Test3", peer, []byte("11111"), func(result []byte) {
		logger.Log.Debug(string(result))
	})
}

func (testModule *TestClientModule) InitializeFinish(service service.IService) {
	testModule.Module.InitializeFinish(service)

	logger.Log.Debug("InitializeFinish")
}
