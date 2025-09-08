package dao

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/muidea/magicCommon/foundation/log"

	imkkc "github.com/muidea/magicCodeAgent/internal/modules/kernel/knowledgebase/pkg/common"
)

func (s *KnowledgeBase) getKnowledgeTableName() string {
	return "kb.tbl_knowledge"
}

func (s *KnowledgeBase) getKnowledgeIndexName() string {
	return "kb_tbl_knowledge"
}

func (s *KnowledgeBase) initKnowledge(dbPtr *sql.DB) error {
	knowledgeTableName := s.getKnowledgeTableName()
	sql := fmt.Sprintf(`
	CREATE SEQUENCE IF NOT EXISTS %s_seq;
	CREATE TABLE IF NOT EXISTS %s (
		id BIGINT PRIMARY KEY DEFAULT nextval('%s_seq'),
		name TEXT,
		description TEXT,
		filePath TEXT,
		warehouse TEXT,
		createdAt BIGINT,
	);
	`, knowledgeTableName, knowledgeTableName, knowledgeTableName)
	_, err := dbPtr.Exec(sql)
	if err != nil {
		log.Errorf("create table %s, error: %v", knowledgeTableName, err)
		return err
	}
	_, err = dbPtr.Exec(fmt.Sprintf(`
	CREATE INDEX IF NOT EXISTS %s_idx_file_warehouse ON %s (filePath, warehouse);`, s.getKnowledgeIndexName(), knowledgeTableName))
	if err != nil {
		log.Errorf("create index %s_idx_file_warehouse error: %v", knowledgeTableName, err)
		return err
	}
	_, err = dbPtr.Exec(fmt.Sprintf(`
	CREATE INDEX IF NOT EXISTS %s_idx_warehouse ON %s (warehouse);`, s.getKnowledgeIndexName(), knowledgeTableName))
	if err != nil {
		log.Errorf("create index %s_idx_warehouse error: %v", knowledgeTableName, err)
		return err
	}

	return nil
}

func (s *KnowledgeBase) QueryWarehouseKnowledge(warehouse string) ([]*imkkc.KnowledgeItem, error) {
	knowledgeTableName := s.getKnowledgeTableName()
	sql := fmt.Sprintf(`
	SELECT id, name, description, filePath, warehouse, createdAt FROM %s WHERE warehouse = ?`, knowledgeTableName)
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
		var description string
		var filePath string
		var warehouse string
		var createdAt int64
		rowsErr = rowsPtr.Scan(&id, &name, &description, &filePath, &warehouse, &createdAt)
		if rowsErr != nil {
			log.Errorf("scan row error: %v", rowsErr)
			return nil, rowsErr
		}
		ret = append(ret, &imkkc.KnowledgeItem{
			ID:          id,
			Name:        name,
			Description: description,
			FilePath:    filePath,
			Warehouse:   warehouse,
			Status:      imkkc.LoadedState,
			CreatedAt:   createdAt,
		})
	}

	return ret, nil
}

func (s *KnowledgeBase) AddWarehouseKnowledge(warehouse, name, description, filePath string) error {
	knowledgeTableName := s.getKnowledgeTableName()
	sql := fmt.Sprintf("INSERT INTO %s(name, description, filePath, warehouse, createdAt) VALUES(?, ?, ?, ?, ?)", knowledgeTableName)
	_, insertErr := s.dbInstancePtr.Exec(sql, name, description, filePath, warehouse, time.Now().UTC().UnixMilli())
	if insertErr != nil {
		log.Errorf("add warehouse knowledge failed, error:%s", insertErr.Error())
		return insertErr
	}
	return insertErr
}

func (s *KnowledgeBase) DelWarehouseKnowledge(warehouse, filePath string) error {
	knowledgeTableName := s.getKnowledgeTableName()
	sql := fmt.Sprintf("DELETE FROM %s WHERE filePath = ? AND warehouse = ?", knowledgeTableName)
	_, deleteErr := s.dbInstancePtr.Exec(sql, filePath, warehouse)
	if deleteErr != nil {
		log.Errorf("delete warehouse knowledge failed, error:%s", deleteErr.Error())
		return deleteErr
	}
	return deleteErr
}

func (s *KnowledgeBase) CleanWarehouseKnowledge(warehouse string) error {
	knowledgeTableName := s.getKnowledgeTableName()
	sql := fmt.Sprintf("DELETE FROM %s WHERE warehouse = ?", knowledgeTableName)
	_, deleteErr := s.dbInstancePtr.Exec(sql, warehouse)
	if deleteErr != nil {
		log.Errorf("delete warehouse knowledge failed, error:%s", deleteErr.Error())
		return deleteErr
	}
	return deleteErr
}
