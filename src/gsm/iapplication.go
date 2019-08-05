package gsm

import (
	"github.com/sf-go/gsf/src/gsc/crypto"
	"github.com/sf-go/gsf/src/gsc/logger"
	"github.com/sf-go/gsf/src/gsc/network"
	"github.com/sf-go/gsf/src/gsc/rpc"
	"github.com/sf-go/gsf/src/gsf/service"
	"github.com/sf-go/gsf/src/gsf/socket"
	"github.com/sf-go/gsf/src/gsm/invoker"
	"github.com/sf-go/gsf/src/gsm/module"
	"github.com/sf-go/gsf/src/gsm/peer"
	"os"
	"os/signal"
	"syscall"
)

type IApplication interface {
	RegisterModule(moduleManager *ModuleManager)
	SetLogConfig(config *logger.LogConfig)
	SetNetConfig(config *network.NetConfig)
	SetCryptoConfig(config *crypto.CryptoConfig)
}

func RunServer(application IApplication, args []string) {
	moduleManager := NewModuleManager()
	netConfig := network.NewNetConfig()
	logConfig := logger.NewLogConfig()

	application.SetLogConfig(logConfig)
	application.SetNetConfig(netConfig)
	application.SetCryptoConfig(crypto.NewCryptoConfig())
	application.RegisterModule(moduleManager)

	rpcRegister := rpc.NewRpcRegister()
	dispatcher := invoker.NewInvoker(rpcRegister)
	serverSocket := socket.NewServerSocket(dispatcher)
	serverService := service.NewServerService(dispatcher, serverSocket)
	moduleInitialize(serverService, moduleManager)
	moduleConnectInitialize(serverSocket.Event, moduleManager)
	serverService.StartServer(netConfig)
}

func RunClient(application IApplication, args []string) {
	moduleManager := NewModuleManager()
	netConfig := network.NewNetConfig()
	logConfig := logger.NewLogConfig()

	application.SetLogConfig(logConfig)
	application.SetNetConfig(netConfig)
	application.SetCryptoConfig(crypto.NewCryptoConfig())
	application.RegisterModule(moduleManager)

	rpcRegister := rpc.NewRpcRegister()
	dispatcher := invoker.NewInvoker(rpcRegister)
	clientSocket := socket.NewClientSocket(dispatcher)
	clientService := service.NewClientService(clientSocket, dispatcher)
	moduleInitialize(clientService, moduleManager)
	moduleConnectInitialize(clientSocket.Event, moduleManager)
	clientService.Connect(netConfig)
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

func moduleConnectInitialize(event *socket.Event, moduleManager *ModuleManager) {
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

func moduleInitialize(service service.IService, moduleManager *ModuleManager) {
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
