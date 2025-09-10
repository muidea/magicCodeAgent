package dao

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/muidea/magicCommon/foundation/log"

	iiec "github.com/muidea/magicCodeAgent/internal/initiators/embedtext/pkg/common"
	imkkc "github.com/muidea/magicCodeAgent/internal/modules/kernel/knowledgebase/pkg/common"
)

func (s *KnowledgeBase) getWarehouseTableName() string {
	return "kb.tbl_warehouse"
}
func (s *KnowledgeBase) getWarehouseIndexName() string {
	return "kb_tbl_warehouse"
}

func (s *KnowledgeBase) getWarehouseVectTableName(warehouse string) string {
	return fmt.Sprintf("kb.tbl_%s_vect", warehouse)
}

func (s *KnowledgeBase) getWarehouseVectIndexName(warehouse string) string {
	return fmt.Sprintf("kb_tbl_%s_vect", warehouse)
}

func (s *KnowledgeBase) initWarehouse(dbPtr *sql.DB) error {
	warehouseTableName := s.getWarehouseTableName()
	sql := fmt.Sprintf(`
	CREATE SEQUENCE IF NOT EXISTS %s_seq;
	CREATE TABLE IF NOT EXISTS %s (
		id BIGINT PRIMARY KEY DEFAULT nextval('%s_seq'),
		name TEXT,
		description TEXT,
		createdAt BIGINT,
	);
	`, warehouseTableName, warehouseTableName, warehouseTableName)
	_, err := dbPtr.Exec(sql)
	if err != nil {
		log.Errorf("create table %s, error: %v", warehouseTableName, err)
		return err
	}
	_, err = dbPtr.Exec(fmt.Sprintf(`
	CREATE INDEX IF NOT EXISTS %s_idx_warehouse ON %s (name);`, s.getWarehouseIndexName(), warehouseTableName))
	if err != nil {
		log.Errorf("create index %s_idx_warehouse error: %v", warehouseTableName, err)
		return err
	}

	return nil
}

func (s *KnowledgeBase) addWarehouse(warehouse, description string) error {
	warehouseTableName := s.getWarehouseTableName()
	_, err := s.dbInstancePtr.Exec(fmt.Sprintf(`
	INSERT INTO %s (name, description, createdAt) VALUES (?, ?, ?)`, warehouseTableName), warehouse, description, time.Now().UTC().UnixMilli())
	if err != nil {
		log.Errorf("insert table %s failed, error:%s", warehouseTableName, err.Error())
		return err
	}
	return nil
}

func (s *KnowledgeBase) deleteWarehouse(warehouse string) error {
	warehouseTableName := s.getWarehouseTableName()
	_, err := s.dbInstancePtr.Exec(fmt.Sprintf(`DELETE FROM %s WHERE name = ?`, warehouseTableName), warehouse)
	if err != nil {
		log.Errorf("delete table %s failed, error:%s", warehouseTableName, err.Error())
		return err
	}
	return nil
}

func (s *KnowledgeBase) CheckWarehouse(warehouse string) error {
	warehouseTableName := s.getWarehouseTableName()
	rowPtr := s.dbInstancePtr.QueryRow(fmt.Sprintf(`SELECT * FROM %s WHERE name = ?`, warehouseTableName), warehouse)
	if rowPtr.Err() != nil {
		log.Errorf("check table %s error:%s", warehouseTableName, rowPtr.Err().Error())
		return rowPtr.Err()
	}
	return nil
}

func (s *KnowledgeBase) ListWarehouse() ([]*imkkc.WarehouseItem, error) {
	warehouseTableName := s.getWarehouseTableName()
	rows, err := s.dbInstancePtr.Query(fmt.Sprintf(`SELECT name, description, createdAt FROM %s`, warehouseTableName))
	if err != nil {
		log.Errorf("query table %s records failed, error:%s", warehouseTableName, err.Error())
		return nil, err
	}
	defer rows.Close()
	var ret []*imkkc.WarehouseItem
	for rows.Next() {
		var name string
		var description string
		var createdAt int64
		err := rows.Scan(&name, &description, &createdAt)
		if err != nil {
			log.Errorf("scan row failed, error:%s", err.Error())
			return nil, err
		}
		ret = append(ret, &imkkc.WarehouseItem{
			Name:        name,
			Description: description,
			CreatedAt:   createdAt,
		})
	}
	return ret, nil
}

func (s *KnowledgeBase) NewWarehouse(warehouse, description string) error {
	warehouseVectTableName := s.getWarehouseVectTableName(warehouse)
	_, err := s.dbInstancePtr.Exec(fmt.Sprintf(`
		CREATE SEQUENCE IF NOT EXISTS %s_seq;
		CREATE TABLE IF NOT EXISTS %s(
			id BIGINT PRIMARY KEY DEFAULT nextval('%s_seq'),
			text TEXT,
			metadata TEXT,
			filePath TEXT,
			embedding FLOAT[512]
		);
	`, warehouseVectTableName, warehouseVectTableName, warehouseVectTableName))
	if err != nil {
		log.Errorf("create table %s error:%s", warehouseVectTableName, err.Error())
		return err
	}

	sql := fmt.Sprintf(`CREATE INDEX IF NOT EXISTS %s_idx_embedding ON %s USING HNSW (embedding) WITH (metric = 'cosine');`, s.getWarehouseVectIndexName(warehouse), warehouseVectTableName)
	_, err = s.dbInstancePtr.Exec(sql)
	if err != nil {
		log.Errorf("create index %s_idx_embedding, sql:%s, error:%s", warehouseVectTableName, sql, err.Error())
		return err
	}
	_, err = s.dbInstancePtr.Exec(fmt.Sprintf(`CREATE INDEX IF NOT EXISTS %s_idx_filePath ON %s(filePath)`, s.getWarehouseVectIndexName(warehouse), warehouseVectTableName))
	if err != nil {
		log.Errorf("create index %s_idx_filePath error:%s", warehouseVectTableName, err.Error())
		return err
	}
	return s.addWarehouse(warehouse, description)
}

func (s *KnowledgeBase) DropWarehouse(warehouse string) error {
	warehouseVectTableName := s.getWarehouseVectTableName(warehouse)
	_, err := s.dbInstancePtr.Exec(fmt.Sprintf(`DROP TABLE IF EXISTS %s`, warehouseVectTableName))
	if err != nil {
		log.Errorf("drop table %s error:%s", warehouseVectTableName, err.Error())
		return err
	}

	return s.deleteWarehouse(warehouse)
}

func (s *KnowledgeBase) AddChunk(warehouse string, chunkPtr *imkkc.Chunk, embedTextWorker iiec.EmbedTextWorker) error {
	warehouseVectTableName := s.getWarehouseVectTableName(warehouse)
	meta, _ := json.Marshal(chunkPtr.Metadata)
	vectorVal, vectorErr := embedTextWorker.EmbedText(chunkPtr.Text)
	if vectorErr != nil {
		log.Errorf("embed text failed, text:%s, err:%v", chunkPtr.Text, vectorErr)
		return vectorErr
	}
	chunkPtr.Vector = vectorVal

	log.Infof("warehouse:%s,chunk:%s", warehouse, chunkPtr)
	sqlStmt := fmt.Sprintf(`INSERT INTO %s (text, metadata, filePath, embedding) VALUES (?, ?, ?, ?)`, warehouseVectTableName)
	_, err := s.dbInstancePtr.Exec(sqlStmt, chunkPtr.Text, string(meta), chunkPtr.FilePath, chunkPtr.Vector)
	return err
}

func (s *KnowledgeBase) DelChunks(warehouse, filePath string) (err error) {
	warehouseVectTableName := s.getWarehouseVectTableName(warehouse)
	sql := fmt.Sprintf("DELETE FROM %s WHERE filePath = ?", warehouseVectTableName)
	_, err = s.dbInstancePtr.Exec(sql, filePath)
	return
}

func (s *KnowledgeBase) SearchByVector(warehouse string, vector []float32, topK int, threshold float32) ([]*imkkc.Chunk, error) {
	warehouseVectTableName := s.getWarehouseVectTableName(warehouse)
	vecLit := floatSliceToArrayLiteral(vector)
	sqlStmt := fmt.Sprintf(`
		SELECT id, text, metadata, filePath, array_cosine_distance(embedding, %s) AS dist
		FROM %s
		WHERE array_cosine_distance(embedding, %s) < %g
		ORDER BY dist
		LIMIT %d;
	`, vecLit, warehouseVectTableName, vecLit, threshold, topK)

	// Execute the query without additional parameters
	rows, err := s.dbInstancePtr.Query(sqlStmt)
	if err != nil {
		log.Errorf("query table %s records failed, error:%s", warehouseVectTableName, err.Error())
		return nil, err
	}
	defer rows.Close()
	var ret []*imkkc.Chunk
	for rows.Next() {
		var id int64
		var text string
		var metadata string
		var filePath string
		var distance float64
		err = rows.Scan(&id, &text, &metadata, &filePath, &distance)
		if err != nil {
			log.Errorf("scan row failed, error:%s", err.Error())
			return nil, err
		}
		chunk := &imkkc.Chunk{
			ID:       id,
			Text:     text,
			Metadata: make(map[string]string),
			FilePath: filePath,
			Distance: distance,
		}
		err = json.Unmarshal([]byte(metadata), &chunk.Metadata)
		if err != nil {
			log.Errorf("unmarshal metadata failed, error:%s", err.Error())
			return nil, err
		}
		ret = append(ret, chunk)
	}

	return ret, nil
}

func (s *KnowledgeBase) SearchByText(warehouse, text string, topK int, threshold float32, embedTextWorker iiec.EmbedTextWorker) ([]*imkkc.Chunk, error) {
	vector, err := embedTextWorker.EmbedText(text)
	if err != nil {
		log.Errorf("embed text failed, error:%s", err.Error())
		return nil, err
	}

	return s.SearchByVector(warehouse, vector, topK, threshold)
}

func (s *KnowledgeBase) SearchByVectorBulk(warehouse string, vectorList [][]float32, topK int, threshold float32) ([]*imkkc.Chunk, error) {
	warehouseVectTableName := s.getWarehouseVectTableName(warehouse)
	var unionSQL string
	for i, vec := range vectorList {
		vecLit := floatSliceToArrayLiteral(vec)
		if i > 0 {
			unionSQL += " UNION ALL "
		}
		unionSQL += fmt.Sprintf(`
			(
				SELECT id, text, metadata, filePath, array_cosine_distance(embedding, %s) AS dist
				FROM %s
				WHERE array_cosine_distance(embedding, %s) < %g
				ORDER BY dist
				LIMIT %d
			)
		`, vecLit, warehouseVectTableName, vecLit, threshold, topK)
	}
	rows, err := s.dbInstancePtr.Query(unionSQL)
	if err != nil {
		log.Errorf("query table %s records failed, error:%s", warehouseVectTableName, err.Error())
		return nil, err
	}
	defer rows.Close()
	var ret []*imkkc.Chunk
	for rows.Next() {
		var id int64
		var text string
		var metadata string
		var filePath string
		var distance float64
		err = rows.Scan(&id, &text, &metadata, &filePath, &distance)
		if err != nil {
			log.Errorf("scan row failed, error:%s", err.Error())
			return nil, err
		}
		chunk := &imkkc.Chunk{
			ID:       id,
			Text:     text,
			Metadata: make(map[string]string),
			FilePath: filePath,
			Distance: distance,
		}
		err = json.Unmarshal([]byte(metadata), &chunk.Metadata)
		if err != nil {
			log.Errorf("unmarshal metadata failed, error:%s", err.Error())
			return nil, err
		}
		ret = append(ret, chunk)
	}

	return ret, nil
}

func (s *KnowledgeBase) SearchByTextBulk(warehouse string, textList []string, topK int, threshold float32, embedTextWorker iiec.EmbedTextWorker) (ret []*imkkc.Chunk, err error) {
	vectorList, err := embedTextWorker.EmbedTextBulk(textList)
	if err != nil {
		log.Errorf("embed text failed, error:%s", err.Error())
		return nil, err
	}
	return s.SearchByVectorBulk(warehouse, vectorList, topK, threshold)
}
