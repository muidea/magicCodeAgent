package biz

import (
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicToolkit/internal/modules/base/biz"
	"github.com/muidea/magicToolkit/internal/modules/kernel/monitor/pkg/common"
)

type Monitor struct {
	biz.Base
}

func New(
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine,
) *Monitor {
	ptr := &Monitor{
		Base: biz.New(common.MonitorModule, eventHub, backgroundRoutine),
	}

	return ptr
}

func (s *Monitor) Initialize() (err *cd.Error) {
	return
}
