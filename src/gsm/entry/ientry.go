package entry

import (
	"github.com/sf-go/gsf/src/gsm/module"
)

type IEntry interface {
	Main() module.IModule
}
