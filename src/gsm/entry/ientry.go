package entry

import (
	"gsm/module"
)

type IEntry interface {
	Main() module.IModule
}
