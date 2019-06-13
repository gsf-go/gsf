package controllers

import (
	"example/server/models"
	"gsc/logger"
	"gsf/peer"
	"gsm/controller"
)

type TestController struct {
	*controller.Controller
}

func NewTestController() *TestController {
	return &TestController{}
}

func (testController *TestController) Initialize() {
	testController.Controller.Initialize()
	testController.Register("Test", func() interface{} {
		return testController.Test
	})
}

func (testController *TestController) Test(num int, testmodel *models.TestModel, peer peer.IPeer) int {
	logger.Log.Debug(testmodel.Name)
	return num + 10000
}
