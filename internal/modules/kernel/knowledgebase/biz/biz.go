package biz

import (
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicCodeAgent/internal/modules/base/biz"
	"github.com/muidea/magicCodeAgent/internal/modules/kernel/knowledgebase/dao"
	"github.com/muidea/magicCodeAgent/internal/modules/kernel/knowledgebase/pkg/common"
)

type KnowledgeBase struct {
	biz.Base

	daoPtr *dao.KnowledgeBase
}

func New(
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine,
) *KnowledgeBase {
	ptr := &KnowledgeBase{
		Base: biz.New(common.KnowledgeBaseModule, eventHub, backgroundRoutine),
	}

	return ptr
}

func (s *KnowledgeBase) Run() (err *cd.Error) {
	return
}
