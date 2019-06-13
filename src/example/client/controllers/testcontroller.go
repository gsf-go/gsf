package controllers

import (
	"github.com/gsf/gsf/src/gsm/controller"
)

type TestController struct {
	*controller.Controller
}

func NewTestController() *TestController {
	return &TestController{}
}

func (testController *TestController) Initialize() {
	testController.Controller.Initialize()
}
