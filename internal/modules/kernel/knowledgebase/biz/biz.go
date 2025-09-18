package biz

import (
	"path/filepath"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/task"

	iiec "github.com/muidea/magicCodeAgent/internal/initiators/embedtext/pkg/common"
	"github.com/muidea/magicCodeAgent/internal/modules/base/biz"
	"github.com/muidea/magicCodeAgent/internal/modules/kernel/knowledgebase/dao"
	imkkc "github.com/muidea/magicCodeAgent/internal/modules/kernel/knowledgebase/pkg/common"
)

type KnowledgeBase struct {
	biz.Base

	daoPtr          *dao.KnowledgeBase
	embedTextWorker iiec.EmbedTextWorker
	baseStaticPath  string
}

func New(
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine,
	embedTextWorker iiec.EmbedTextWorker,
	basePath string,
) *KnowledgeBase {
	ptr := &KnowledgeBase{
		Base:            biz.New(imkkc.KnowledgeBaseModule, eventHub, backgroundRoutine),
		daoPtr:          dao.New(filepath.Join(basePath, imkkc.DataPath)),
		embedTextWorker: embedTextWorker,
		baseStaticPath:  filepath.Join(basePath, imkkc.StaticPath),
	}

	ptr.SubscribeFunc(imkkc.HybridSearchWarehouseKnowledge, ptr.searchKnowledge)

	return ptr
}

func (s *KnowledgeBase) Run() (err *cd.Error) {
	s.daoPtr.Open()
	return
}

func (s *KnowledgeBase) Stop() {
	s.daoPtr.Close()
}

func (s *KnowledgeBase) searchKnowledge(ev event.Event, re event.Result) {
	paramVal := ev.Data()
	if paramVal == nil {
		if re != nil {
			re.Set(nil, cd.NewError(cd.IllegalParam, "illegal param"))
		}
		return
	}

	paramPtr := paramVal.(*imkkc.HybridSearchKnowledgeParam)
	chunkList, chunkErr := s.HybridSearchKnowledge(
		paramPtr.Warehouse,
		paramPtr.Text,
		paramPtr.TopK,
		paramPtr.BusinessContext,
		paramPtr.Similarity,
		paramPtr.ContextType)
	if chunkErr != nil {
		if re != nil {
			re.Set(nil, cd.NewError(cd.Unexpected, chunkErr.Error()))
		}
		return
	}

	if re != nil {
		re.Set(chunkList, nil)
	}
}
