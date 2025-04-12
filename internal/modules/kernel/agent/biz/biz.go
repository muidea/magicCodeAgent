package biz

import (
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicCodeAgent/internal/modules/base/biz"
	"github.com/muidea/magicCodeAgent/internal/modules/kernel/agent/pkg/common"
)

type Agent struct {
	biz.Base
}

func New(
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine,
) *Agent {
	ptr := &Agent{
		Base: biz.New(common.AgentModule, eventHub, backgroundRoutine),
	}

	return ptr
}

func (s *Agent) Run() (err *cd.Error) {
	return
}
