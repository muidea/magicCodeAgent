package biz

import (
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicCodeAgent/internal/modules/base/biz"
	"github.com/muidea/magicCodeAgent/internal/modules/blocks/debugger/pkg/common"
)

type Debugger struct {
	biz.Base
}

func New(
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine,
) *Debugger {
	ptr := &Debugger{
		Base: biz.New(common.DebuggerModule, eventHub, backgroundRoutine),
	}

	return ptr
}

func (s *Debugger) Run() (err *cd.Error) {
	return
}
