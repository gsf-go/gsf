package controllers

import (
	"github.com/sf-go/gsf/src/gsm/dispatcher"
)

type TestController struct {
	dispatcher dispatcher.IDispatcher
}

func (controller *TestController) GetName() string {
	return "TestController"
}

func NewTestController(dispatcher dispatcher.IDispatcher) *TestController {
	return &TestController{
		dispatcher: dispatcher,
	}
}
