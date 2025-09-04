package biz

import (
	"os"
	"path/filepath"

	fp "github.com/muidea/magicCommon/foundation/path"
)

func (s *KnowledgeBase) NewWarehouse(warehouse string) error {
	warehousePath := filepath.Join(s.baseStaticPath, warehouse)
	if fp.Exist(warehousePath) {
		return nil
	}

	_ = os.MkdirAll(warehousePath, 0755)
	return s.daoPtr.NewWarehouse(warehouse)
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

func (s *KnowledgeBase) ListWarehouse() (ret []string, err error) {
	entryList, entryErr := os.ReadDir(s.baseStaticPath)
	if entryErr != nil {
		err = entryErr
		return
	}

	for _, entry := range entryList {
		if entry.IsDir() {
			ret = append(ret, entry.Name())
		}
	}
	return
}

func (s *KnowledgeBase) ExistWorkspace(workspace string) bool {
	warehousePath := filepath.Join(s.baseStaticPath, workspace)
	_, err := os.Stat(warehousePath)
	return err == nil
}
