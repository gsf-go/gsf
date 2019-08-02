package controller

import "github.com/sf-go/gsf/src/gsm/dispatcher"

type IController interface {
	Initialize(dispatcher dispatcher.IDispatcher)
}
