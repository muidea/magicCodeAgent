package biz

import (
	"os"

	"github.com/muidea/magicCommon/foundation/log"

	imkkc "github.com/muidea/magicCodeAgent/internal/modules/kernel/knowledgebase/pkg/common"
)

func (s *KnowledgeBase) ImportKnowledge(warehouse, filePath string) error {
	if err := s.daoPtr.CheckWarehouse(warehouse); err != nil {
		log.Errorf("check warehouse %s failed, error:%s", warehouse, err.Error())
		return err
	}

	mdBytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Errorf("open file failed, filePath:%s, err:%v", filePath, err)
		return err
	}

	chunks := imkkc.SplitMarkdown(filePath, string(mdBytes), 300, 60)
	err = s.daoPtr.AddChunksBulk(warehouse, chunks, s.embedTextWorker)
	if err != nil {
		log.Errorf("add chunks failed, warehouseName:%s, err:%v", warehouse, err)
		return err
	}

	return nil
}

func (s *KnowledgeBase) DeleteKnowledge(warehouse, filePath string) error {
	_ = s.daoPtr.DelChunks(warehouse, filePath)
	_ = s.daoPtr.DelWarehouseKnowledge(warehouse, filePath)
	return nil
}

func (s *KnowledgeBase) ListKnowledge(warehouse string) (ret []string, err error) {
	return
}
