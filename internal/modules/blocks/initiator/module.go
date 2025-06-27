package module

import (
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/framework/plugin/module"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicCodeAgent/internal/modules/blocks/initiator/biz"
	"github.com/muidea/magicCodeAgent/internal/modules/blocks/initiator/pkg/common"
)

func init() {
	module.Register(New())
}

func New() *Initiator {
	return &Initiator{}
}

type Initiator struct {
	bizPtr *biz.Initiator
}

func (s *Initiator) ID() string {
	return common.InitiatorModule
}

func (s *Initiator) Setup(eventHub event.Hub, backgroundRoutine task.BackgroundRoutine) (err *cd.Error) {

	s.bizPtr = biz.New(eventHub, backgroundRoutine)
	return
}

func (s *Initiator) Run() (err *cd.Error) {
	return
}

func (s *Initiator) Teardown() {
}
