package controllers

import (
	"github.com/sf-go/gsf/src/example/server/components"
	"github.com/sf-go/gsf/src/example/server/models"
	"github.com/sf-go/gsf/src/gsc/logger"
	"github.com/sf-go/gsf/src/gsf/peer"
	"github.com/sf-go/gsf/src/gsm/controller"
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
			return true
		}
	}, func() interface{} {
		return func(num int, testmodel *models.TestModel, peer peer.IPeer) {
			logger.Log.Debug("oooooooooooooooooooo")
		}
	})

	testController.Register("Test2", func() interface{} {
		return testController.Test2
	}, nil, nil)
}

func (testController *TestController) Test2(num int32, peer peer.IPeer) int32 {
	return num + 10000
}

func (testController *TestController) Test(num int, testmodel *models.TestModel, peer peer.IPeer) int {
	logger.Log.Debug(testmodel.Name)

	cmpt := peer.GetComponent("UserComponent").(*components.UserComponent)
	logger.Log.Debug(cmpt.Account)
	logger.Log.Debug(cmpt.Password)
	logger.Log.Debug(testmodel.Name)
	return num + 10000
}
