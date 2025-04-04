package module

import (
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/log"
	"github.com/muidea/magicCommon/framework/plugin/module"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicToolkit/internal/config"
	"github.com/muidea/magicToolkit/internal/modules/module/biz"
	"github.com/muidea/magicToolkit/internal/modules/module/pkg/common"
	"github.com/muidea/magicToolkit/pkg/util"
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
	currentWorkPath := config.CurrentWorkPath()
	err = util.VerifyModuleWorkPath(currentWorkPath)
	if err != nil {
		log.Errorf("verify work path failed, workPath:%s, err:%s", currentWorkPath, err)
		return
	}

	return
}

func (s *Module) Teardown() {
}
