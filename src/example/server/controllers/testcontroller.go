package controllers

import (
	"gsc/logger"
	"gsm/controller"
)

type TestController struct {
	*controller.Controller
}

func NewTestController() *TestController {
	return &TestController{}
}

func (testController *TestController) Initialize() {
	testController.Register("Test", func() interface{} {
		return testController.Test
	})
}

func (testController *TestController) Test(num int) {
	logger.Log.Debug(string(num))
}
