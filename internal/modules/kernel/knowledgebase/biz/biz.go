package biz

import (
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicCodeAgent/internal/modules/base/biz"
	"github.com/muidea/magicCodeAgent/internal/modules/kernel/knowledgebase/dao"
	imkkc "github.com/muidea/magicCodeAgent/internal/modules/kernel/knowledgebase/pkg/common"
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
		Base:   biz.New(imkkc.KnowledgeBaseModule, eventHub, backgroundRoutine),
		daoPtr: dao.New(),
	}

	return ptr
}

func (s *KnowledgeBase) Run() (err *cd.Error) {
	s.daoPtr.Open()
	return
}

func (s *KnowledgeBase) Stop() {
	s.daoPtr.Close()
}
