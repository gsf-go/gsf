package controllers

import "github.com/sf-go/gsf/src/gsm/invoker"

type TestController struct {
	dispatcher invoker.IInvoker
}

func (controller *TestController) GetName() string {
	return "TestController"
}

func NewTestController(dispatcher invoker.IInvoker) *TestController {
	return &TestController{
		dispatcher: dispatcher,
	}
}
