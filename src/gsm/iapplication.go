package gsm

import (
	"github.com/sf-go/gsf/src/gsc/crypto"
	"github.com/sf-go/gsf/src/gsc/logger"
	"github.com/sf-go/gsf/src/gsc/network"
	"github.com/sf-go/gsf/src/gsf/service"
	"github.com/sf-go/gsf/src/gsf/socket"
	"github.com/sf-go/gsf/src/gsm/module"
	"github.com/sf-go/gsf/src/gsm/peer"
	"os"
	"os/signal"
	"syscall"
)

type IApplication interface {
	RegisterModule(moduleManager *module.ModuleManager)
	SetLogConfig(config *logger.LogConfig)
	SetNetConfig(config *network.NetConfig)
	SetCryptoConfig(config *crypto.CryptoConfig)
}

func RunServer(application IApplication, args []string) {
	moduleManager := module.NewModuleManager()
	netConfig := network.NewNetConfig()
	logConfig := logger.NewLogConfig()

	application.SetLogConfig(logConfig)
	logger.Log.SetConfig(logConfig)
	application.SetNetConfig(netConfig)
	application.SetCryptoConfig(crypto.NewCryptoConfig())
	application.RegisterModule(moduleManager)

	serverSocket := socket.NewServerSocket()
	s := service.NewServerService(serverSocket)
	moduleInitialize(s, moduleManager)
	moduleConnectInitialize(serverSocket.Event, moduleManager)
	s.StartServer(netConfig)
}

func RunClient(application IApplication, args []string) {
	moduleManager := module.NewModuleManager()
	netConfig := network.NewNetConfig()
	logConfig := logger.NewLogConfig()

	application.SetLogConfig(logConfig)
	logger.Log.SetConfig(logConfig)
	application.SetNetConfig(netConfig)
	application.SetCryptoConfig(crypto.NewCryptoConfig())
	application.RegisterModule(moduleManager)

	clientSocket := socket.NewClientSocket()
	s := service.NewClientService(clientSocket)
	moduleInitialize(s, moduleManager)
	moduleConnectInitialize(clientSocket.Event, moduleManager)
	s.Connect(netConfig)
}

func shutdown(callback func()) {
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)
	signal.Notify(s, syscall.SIGTERM)
	go func() {
		<-s
		if callback != nil {
			callback()
		}
		os.Exit(0)
	}()
}

func moduleConnectInitialize(event *socket.Event, moduleManager *module.ModuleManager) {
	event.OnConnected = func(peer peer.IPeer) {
		moduleManager.Range(func(key, value interface{}) bool {
			m := value.(module.IModule)
			m.Connected(peer)
			return true
		})
	}

	event.OnDisconnected = func(peer peer.IPeer) {
		moduleManager.Range(func(key, value interface{}) bool {
			m := value.(module.IModule)
			m.Disconnected(peer)
			return true
		})
	}
}

func moduleInitialize(service service.IService, moduleManager *module.ModuleManager) {
	moduleManager.Range(func(key, value interface{}) bool {
		m := value.(module.IModule)
		m.Initialize(service)
		return true
	})

	moduleManager.Range(func(key, value interface{}) bool {
		m := value.(module.IModule)
		m.InitializeFinish(service)
		return true
	})
}
