package module

import (
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/log"
	"github.com/muidea/magicCommon/framework/plugin/module"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicToolkit/internal/config"
	"github.com/muidea/magicToolkit/internal/modules/initator/biz"
	"github.com/muidea/magicToolkit/internal/modules/initator/pkg/common"
	"github.com/muidea/magicToolkit/pkg/util"
)

func init() {
	module.Register(New())
}

func New() *Initator {
	return &Initator{}
}

type Initator struct {
	bizPtr *biz.Initator
}

func (s *Initator) ID() string {
	return common.InitatorModule
}

func (s *Initator) Setup(eventHub event.Hub, backgroundRoutine task.BackgroundRoutine) (err *cd.Error) {

	s.bizPtr = biz.New(eventHub, backgroundRoutine)
	return
}

func (s *Initator) Run() (err *cd.Error) {
	currentWorkPath := config.CurrentWorkPath()
	err = util.VerifyModuleWorkPath(currentWorkPath)
	if err != nil {
		log.Errorf("verify work path failed, workPath:%s, err:%s", currentWorkPath, err)
		return
	}

	return
}

func (s *Initator) Teardown() {
}
