package controllers

import (
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
	testController.Register("Test", func() interface{} {
		return testController.Test
	})
}

func (testController *TestController) Test(num int, peer peer.IPeer) int {
	return num + 10000
}
