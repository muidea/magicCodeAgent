package module

import (
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/framework/plugin/module"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicCodeAgent/internal/modules/blocks/module/biz"
	"github.com/muidea/magicCodeAgent/internal/modules/blocks/module/pkg/common"
)

func init() {
	module.Register(New())
}

func New() *Module {
	return &Module{}
}

type Module struct {
	bizPtr *biz.Module
}

func (s *Module) ID() string {
	return common.ModuleModule
}

func (s *Module) Setup(eventHub event.Hub, backgroundRoutine task.BackgroundRoutine) (err *cd.Error) {

	s.bizPtr = biz.New(eventHub, backgroundRoutine)
	return
}

func (s *Module) Run() (err *cd.Error) {
	return
}

func (s *Module) Teardown() {
}
