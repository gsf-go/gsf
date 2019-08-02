package controller

import (
	"github.com/sf-go/gsf/src/gsm/dispatcher"
)

type RawController struct {
	rpcDispatcher *dispatcher.RawDispatcher
}

func NewRawController() *RawController {
	contoller := &RawController{
		rpcDispatcher: dispatcher.NewRawDispatcher(),
	}
	contoller.Initialize(contoller.rpcDispatcher)
	return contoller
}

func (controller *RawController) Initialize(dispatcher dispatcher.IDispatcher) {

}
