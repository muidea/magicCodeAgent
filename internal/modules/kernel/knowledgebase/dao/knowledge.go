package dao

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/muidea/magicCommon/foundation/log"

	imkkc "github.com/muidea/magicCodeAgent/internal/modules/kernel/knowledgebase/pkg/common"
)

func (s *KnowledgeBase) init(dbPtr *sql.DB) error {
	knowledgeTableName := "tbl_knowledge"
	sql := fmt.Sprintf(`
	CREATE SEQUENCE IF NOT EXISTS tbl_knowledge_seq;
	CREATE TABLE IF NOT EXISTS %s (
		id BIGINT PRIMARY KEY DEFAULT nextval('tbl_knowledge_seq'),
		name TEXT,
		filePath TEXT,
		warehouse TEXT,
		createdAt BIGINT,
	);
	`, knowledgeTableName)
	_, err := dbPtr.Exec(sql)
	if err != nil {
		log.Errorf("create table %s, error: %v", knowledgeTableName, err)
		return err
	}
	_, err = dbPtr.Exec(fmt.Sprintf(`
	CREATE INDEX IF NOT EXISTS %s_idx_file_warehouse ON %s (filePath, warehouse);`, knowledgeTableName, knowledgeTableName))
	if err != nil {
		log.Errorf("create index %s_idx_file_warehouse error: %v", knowledgeTableName, err)
		return err
	}
	_, err = dbPtr.Exec(fmt.Sprintf(`
	CREATE INDEX IF NOT EXISTS %s_idx_warehouse ON %s (warehouse);`, knowledgeTableName, knowledgeTableName))
	if err != nil {
		log.Errorf("create index %s_idx_warehouse error: %v", knowledgeTableName, err)
		return err
	}

	return nil
}

func (s *KnowledgeBase) QueryWarehouseKnowledge(warehouse string) ([]*imkkc.KnowledgeItem, error) {
	knowledgeTableName := "tbl_knowledge"
	sql := fmt.Sprintf(`
	SELECT id, name, filePath, warehouse, createdAt FROM %s WHERE warehouse = ?`, knowledgeTableName)
	rowsPtr, rowsErr := s.dbInstancePtr.Query(sql, warehouse)
	if rowsErr != nil {
		log.Errorf("query warehouse %s knowledge error: %v", warehouse, rowsErr)
		return nil, rowsErr
	}
	defer rowsPtr.Close()

	ret := []*imkkc.KnowledgeItem{}
	for rowsPtr.Next() {
		var id uint64
		var name string
		var filePath string
		var warehouse string
		var createdAt int64
		rowsErr = rowsPtr.Scan(&id, &name, &filePath, &warehouse, &createdAt)
		if rowsErr != nil {
			log.Errorf("scan row error: %v", rowsErr)
			return nil, rowsErr
		}
		ret = append(ret, &imkkc.KnowledgeItem{
			ID:        id,
			Name:      name,
			FilePath:  filePath,
			Warehouse: warehouse,
			Status:    imkkc.LoadedState,
			CreatedAt: createdAt,
		})
	}

	return ret, nil
}

func (s *KnowledgeBase) AddWarehouseKnowledge(warehouse, name, filePath string) error {
	knowledgeTableName := "tbl_knowledge"
	sql := fmt.Sprintf("INSERT INTO %s(name, filePath, warehouse, createdAt) VALUES(?, ?, ?, ?)", knowledgeTableName)
	_, insertErr := s.dbInstancePtr.Exec(sql, name, filePath, warehouse, time.Now().UTC().UnixMilli())
	if insertErr != nil {
		log.Errorf("add warehouse knowledge failed, error:%s", insertErr.Error())
		return insertErr
	}
	return insertErr
}

func (s *KnowledgeBase) DelWarehouseKnowledge(warehouse, filePath string) error {
	knowledgeTableName := "tbl_knowledge"
	sql := fmt.Sprintf("DELETE FROM %s WHERE filePath = ? AND warehouse = ?", knowledgeTableName)
	_, deleteErr := s.dbInstancePtr.Exec(sql, filePath, warehouse)
	if deleteErr != nil {
		log.Errorf("delete warehouse knowledge failed, error:%s", deleteErr.Error())
		return deleteErr
	}
	return deleteErr
}

func (s *KnowledgeBase) CleanWarehouseKnowledge(warehouse string) error {
	knowledgeTableName := "tbl_knowledge"
	sql := fmt.Sprintf("DELETE FROM %s WHERE warehouse = ?", knowledgeTableName)
	_, deleteErr := s.dbInstancePtr.Exec(sql, warehouse)
	if deleteErr != nil {
		log.Errorf("delete warehouse knowledge failed, error:%s", deleteErr.Error())
		return deleteErr
	}
	return deleteErr
}
