package gsm

import (
	"github.com/sf-go/gsf/src/gsm/module"
	"sync"
)

type ModuleManager struct {
	modules *sync.Map
}

func NewModuleManager() *ModuleManager {
	return &ModuleManager{
		modules: new(sync.Map),
	}
}

func (moduleManager *ModuleManager) AddModule(name string, module module.IModule) {
	moduleManager.modules.Store(name, module)
}

func (moduleManager *ModuleManager) RemoveModule(name string) {
	moduleManager.modules.Delete(name)
}

func (moduleManager *ModuleManager) GetModule(name string) module.IModule {
	m, ok := moduleManager.modules.Load(name)
	if ok {
		return m.(module.IModule)
	}
	return nil
}

func (moduleManager *ModuleManager) Range(f func(key, value interface{}) bool) {
	moduleManager.modules.Range(f)
}
