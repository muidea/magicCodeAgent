package biz

import (
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicCodeAgent/internal/modules/base/biz"
	"github.com/muidea/magicCodeAgent/internal/modules/kernel/debug/pkg/common"
)

type Debug struct {
	biz.Base
}

func New(
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine,
) *Debug {
	ptr := &Debug{
		Base: biz.New(common.DebugModule, eventHub, backgroundRoutine),
	}

	return ptr
}

func (s *Debug) Run() (err *cd.Error) {
	return
}
