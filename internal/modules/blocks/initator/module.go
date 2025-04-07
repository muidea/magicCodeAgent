package module

import (
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/framework/plugin/module"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicToolkit/internal/modules/blocks/initator/biz"
	"github.com/muidea/magicToolkit/internal/modules/blocks/initator/pkg/common"
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
	return
}

func (s *Initator) Teardown() {
}
