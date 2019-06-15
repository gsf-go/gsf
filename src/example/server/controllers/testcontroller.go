package controllers

import (
	"github.com/gsf/gsf/src/example/server/models"
	"github.com/gsf/gsf/src/gsc/logger"
	"github.com/gsf/gsf/src/gsf/peer"
	"github.com/gsf/gsf/src/gsm/controller"
)

type TestController struct {
	*controller.Controller
}

func NewTestController() *TestController {
	return &TestController{
		Controller: controller.NewController(),
	}
}

func (testController *TestController) Initialize() {
	testController.Controller.Initialize()
	testController.Register("Test", func() interface{} {
		return testController.Test
	}, func() interface{} {
		return func(num int, testmodel *models.TestModel, peer peer.IPeer) bool {
			logger.Log.Debug("xxxxxxxxxxxxxxxxxxxx")
			return false
		}
	}, func() interface{} {
		return func(num int, testmodel *models.TestModel, peer peer.IPeer) {
			logger.Log.Debug("oooooooooooooooooooo")
		}
	})
}

func (testController *TestController) Test(num int, testmodel *models.TestModel, peer peer.IPeer) int {
	logger.Log.Debug(testmodel.Name)
	return num + 10000
}
