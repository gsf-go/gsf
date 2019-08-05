package controllers

import (
	"github.com/sf-go/gsf/src/example/server/models"
	"github.com/sf-go/gsf/src/gsc/logger"
	"github.com/sf-go/gsf/src/gsm/invoker"
	"github.com/sf-go/gsf/src/gsm/peer"
)

type TestController struct {
	dispatcher invoker.IInvoker
}

func (testController *TestController) GetName() string {
	return "TestController"
}

func NewTestController(dispatcher invoker.IInvoker) *TestController {
	controller := &TestController{
		dispatcher: dispatcher,
	}
	controller.Initialize()
	return controller
}

func (testController *TestController) Initialize() {
	testController.dispatcher.Register("Test",
		func() interface{} {
			return func(num int, testmodel *models.TestModel, peer peer.IPeer) bool {
				logger.Log.Debug("xxxxxxxxxxxxxxxxxxxx")
				return true
			}
		},
		func() interface{} {
			return testController.Test
		}, func() interface{} {
			return func(num int, testmodel *models.TestModel, peer peer.IPeer) {
				logger.Log.Debug("oooooooooooooooooooo")
			}
		})

	testController.dispatcher.Register("Test2", nil, func() interface{} {
		return testController.Test2
	}, nil)

	testController.dispatcher.RawRegister("Test3", nil, testController.Test3, nil)
}

func (testController *TestController) Test3(p peer.IPeer, data []byte) []byte {
	logger.Log.Info(string(data))
	return []byte("Hello!")
}

func (testController *TestController) Test2(num int32, peer peer.IPeer) bool {
	return true
}

func (testController *TestController) Test(num int, testmodel *models.TestModel, peer peer.IPeer) int {
	logger.Log.Debug(testmodel.Name)
	return num + 10000
}
