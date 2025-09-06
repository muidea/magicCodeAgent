package biz

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/muidea/magicCommon/foundation/log"
	fp "github.com/muidea/magicCommon/foundation/path"

	imkkc "github.com/muidea/magicCodeAgent/internal/modules/kernel/knowledgebase/pkg/common"
)

func (s *KnowledgeBase) InsertKnowlege(warehouse, name, description, filePath string) error {
	if err := s.importKnowledge(warehouse, filePath); err != nil {
		return err
	}

	return s.daoPtr.AddWarehouseKnowledge(warehouse, name, description, filePath)
}

func (s *KnowledgeBase) importKnowledge(warehouse, filePath string) error {
	_ = s.daoPtr.DelWarehouseKnowledge(warehouse, filePath)
	_ = s.daoPtr.DelChunks(warehouse, filePath)

	if err := s.daoPtr.CheckWarehouse(warehouse); err != nil {
		log.Errorf("check warehouse %s failed, error:%s", warehouse, err.Error())
		return err
	}

	finalFilePath := filepath.Join(s.baseStaticPath, filePath)
	mdBytes, err := os.ReadFile(finalFilePath)
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
	finalFilePath := filepath.Join(s.baseStaticPath, filePath)
	_ = os.Remove(finalFilePath)
	return nil
}

func (s *KnowledgeBase) ListKnowledge(warehouse string) (ret []*imkkc.KnowledgeItem, err error) {
	warehousePath := filepath.Join(s.baseStaticPath, warehouse)
	if !fp.Exist(warehousePath) {
		err = fmt.Errorf("warehouse:%s not exist", warehouse)
		return
	}

	dbItemList, dbItemErr := s.daoPtr.QueryWarehouseKnowledge(warehouse)
	if dbItemErr != nil {
		err = dbItemErr
		return
	}

	ret = dbItemList
	return
}

func (s *KnowledgeBase) GetKnowledgeContent(filePath string) (ret string, err error) {
	finalFilePath := filepath.Join(s.baseStaticPath, filePath)
	mdBytes, mdErr := os.ReadFile(finalFilePath)
	if mdErr != nil {
		err = mdErr
		log.Errorf("open file failed, filePath:%s, err:%v", filePath, err)
		return
	}

	ret = string(mdBytes)
	return
}
