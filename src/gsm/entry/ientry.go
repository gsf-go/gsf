package entry

import (
	"github.com/gsf/gsf/src/gsm/module"
)

type IEntry interface {
	Main() module.IModule
}
