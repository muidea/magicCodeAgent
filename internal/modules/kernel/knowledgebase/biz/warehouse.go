package biz

import (
	"os"
	"path/filepath"

	"github.com/muidea/magicCommon/foundation/log"
	fp "github.com/muidea/magicCommon/foundation/path"

	imkkc "github.com/muidea/magicCodeAgent/internal/modules/kernel/knowledgebase/pkg/common"
)

func (s *KnowledgeBase) NewWarehouse(warehouse, description string) error {
	warehousePath := filepath.Join(s.baseStaticPath, warehouse)
	if fp.Exist(warehousePath) {
		return nil
	}

	_ = os.MkdirAll(warehousePath, 0755)
	return s.daoPtr.NewWarehouse(warehouse, description)
}

func (s *KnowledgeBase) DropWarehouse(warehouse string) error {
	warehousePath := filepath.Join(s.baseStaticPath, warehouse)
	if !fp.Exist(warehousePath) {
		return nil
	}

	if err := s.daoPtr.DropWarehouse(warehouse); err != nil {
		return err
	}

	if err := s.daoPtr.CleanWarehouseKnowledge(warehouse); err != nil {
		return err
	}

	_ = os.RemoveAll(warehousePath)
	return nil
}

func (s *KnowledgeBase) ListWarehouse() (ret []*imkkc.WarehouseItem, err error) {
	return s.daoPtr.ListWarehouse()
}

func (s *KnowledgeBase) ExistWarehouse(warehouse string) bool {
	if err := s.daoPtr.CheckWarehouse(warehouse); err != nil {
		return false
	}

	warehousePath := filepath.Join(s.baseStaticPath, warehouse)
	_, err := os.Stat(warehousePath)
	return err == nil
}

func (s *KnowledgeBase) SearchKnowledge(warehouse, text string, topK int, threshold float32) (ret []*imkkc.Chunk, err error) {
	log.Infof("search knowledgewarehouse, warehouse:%s, text:%s, topK:%d, threshold:%g", warehouse, text, topK, threshold)
	return s.daoPtr.SearchByText(warehouse, text, topK, threshold, s.embedTextWorker)
}

func (s *KnowledgeBase) SearchKnowledgeBulk(warehouse string, texts []string, topK int, threshold float32) (ret []*imkkc.Chunk, err error) {
	log.Infof("search knowledgewarehouse, warehouse:%s, text:%v, topK:%d, threshold:%g", warehouse, texts, topK, threshold)
	return s.daoPtr.SearchByTextBulk(warehouse, texts, topK, threshold, s.embedTextWorker)
}

func (s *KnowledgeBase) HybridSearchKnowledge(warehouse string, text string, topK int, threshold float32, contextType imkkc.ContextType) (ret []*imkkc.Chunk, err error) {
	return s.daoPtr.HybridSearchByText(warehouse, text, topK, threshold, contextType, s.embedTextWorker)
}
