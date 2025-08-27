package biz

import (
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicCodeAgent/internal/modules/base/biz"
	"github.com/muidea/magicCodeAgent/internal/modules/kernel/knowledge/pkg/common"
)

type Knowledge struct {
	biz.Base
}

func New(
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine,
) *Knowledge {
	ptr := &Knowledge{
		Base: biz.New(common.KnowledgeModule, eventHub, backgroundRoutine),
	}

	return ptr
}

func (s *Knowledge) Run() (err *cd.Error) {
	return
}
