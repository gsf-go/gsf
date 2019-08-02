package controller

import "github.com/sf-go/gsf/src/gsm/dispatcher"

type Controller struct {
	Dispatcher *dispatcher.IDispatcher
}

func NewController() *Controller {
	return &Controller{}
}

func (controller *Controller) Initialize(dispatcher dispatcher.IDispatcher) {

}
