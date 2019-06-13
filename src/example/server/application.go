package main

import (
	"github.com/gsf/gsf/src/example/server/modules"
	"github.com/gsf/gsf/src/gsc/crypto"
	"github.com/gsf/gsf/src/gsc/logger"
	"github.com/gsf/gsf/src/gsc/network"
	"github.com/gsf/gsf/src/gsm/module"
)

type Application struct {
}

func NewApplication() *Application {
	return &Application{}
}

func (application *Application) RegisterModule(moduleManager *module.ModuleManager) {
	moduleManager.AddModule("TestServerModule", modules.NewTestServerModule())
}

func (application *Application) SetLogConfig(config *logger.LogConfig) {

}

func (application *Application) SetNetConfig(config *network.NetConfig) {
	config.BufferSize = 50
	config.Address = "127.0.0.1"
	config.Port = 8889
	config.ConnectTimeout = 3
}

func (application *Application) SetCryptoConfig(config *crypto.CryptoConfig) {

}
