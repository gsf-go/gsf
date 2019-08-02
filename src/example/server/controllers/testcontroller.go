package controllers

import (
	"github.com/sf-go/gsf/src/example/server/components"
	"github.com/sf-go/gsf/src/example/server/models"
	"github.com/sf-go/gsf/src/gsc/logger"
	"github.com/sf-go/gsf/src/gsm/dispatcher"
	"github.com/sf-go/gsf/src/gsm/peer"
)

type TestController struct {
	dispatcher dispatcher.IDispatcher
}

func (testController *TestController) GetName() string {
	return "TestController"
}

func NewTestController(dispatcher dispatcher.IDispatcher) *TestController {
	controller := &TestController{
		dispatcher: dispatcher,
	}
	controller.Initialize()
	return controller
}

func (testController *TestController) Initialize() {
	testController.dispatcher.Register([]byte("Test"),
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

	testController.dispatcher.Register([]byte("Test2"), func() interface{} {
		return testController.Test2
	}, nil, nil)
}

func (testController *TestController) Test2(num int32, peer peer.IPeer) bool {
	return true
}

func (testController *TestController) Test(num int, testmodel *models.TestModel, peer peer.IPeer) int {
	logger.Log.Debug(testmodel.Name)

	cmpt := peer.GetComponent("UserComponent").(*components.UserComponent)
	logger.Log.Debug(cmpt.Account)
	logger.Log.Debug(cmpt.Password)
	logger.Log.Debug(testmodel.Name)
	return num + 10000
}
