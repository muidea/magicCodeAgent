package biz

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/muidea/magicCommon/foundation/log"
	fp "github.com/muidea/magicCommon/foundation/path"

	imkkc "github.com/muidea/magicCodeAgent/internal/modules/kernel/knowledgebase/pkg/common"
)

func (s *KnowledgeBase) ImportKnowledge(warehouse, filePath string) error {
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

	name := filepath.Base(filePath)
	return s.daoPtr.AddWarehouseKnowledge(warehouse, name, filePath)
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

	entryList, entryErr := os.ReadDir(warehousePath)
	if entryErr != nil {
		err = entryErr
		return
	}

	dbItemList, dbItemErr := s.daoPtr.QueryWarehouseKnowledge(warehouse)
	if dbItemErr != nil {
		err = dbItemErr
		return
	}

	for _, entry := range entryList {
		if entry.IsDir() {
			continue
		}

		filePath := filepath.Join(warehouse, entry.Name())
		existFlag := false
		for _, sv := range dbItemList {
			if filePath == sv.FilePath {
				existFlag = true
				ret = append(ret, sv)
				break
			}
		}
		if !existFlag {
			fileInfo, fileErr := entry.Info()
			if fileErr != nil {
				continue
			}

			ret = append(ret, &imkkc.KnowledgeItem{
				Name:      entry.Name(),
				FilePath:  filePath,
				Warehouse: warehouse,
				Status:    imkkc.IdleState,
				CreatedAt: fileInfo.ModTime().UTC().UnixMilli(),
			})
		}
	}

	return
}
