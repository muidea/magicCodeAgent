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

func (s *KnowledgeBase) SearchKnowledge(warehouse, text string, topK int) (ret []*imkkc.Chunk, err error) {
	log.Infof("search knowledgewarehouse, warehouse:%s, text:%s, topK:%d", warehouse, text, topK)
	return s.daoPtr.SearchByText(warehouse, text, topK, s.embedTextWorker)
}
