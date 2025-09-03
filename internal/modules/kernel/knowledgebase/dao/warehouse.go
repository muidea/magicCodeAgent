package dao

import (
	"encoding/json"
	"fmt"

	"github.com/muidea/magicCommon/foundation/log"

	iiec "github.com/muidea/magicCodeAgent/internal/initiators/embedtext/pkg/common"
	imkkc "github.com/muidea/magicCodeAgent/internal/modules/kernel/knowledgebase/pkg/common"
)

func (s *KnowledgeBase) NewWarehouse(warehouse string) error {
	warehouseTableName := fmt.Sprintf("tbl_%s", warehouse)
	_, err := s.dbInstancePtr.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s(
			id BIGINT AUTO_INCREMENT,
			text TEXT,
			metadata TEXT,
			filePath TEXT,
			embedding FLOAT[]
		);
	`, warehouseTableName))
	if err != nil {
		log.Errorf("create table %s error:%s", warehouseTableName, err.Error())
		return err
	}

	_, err = s.dbInstancePtr.Exec(fmt.Sprintf(`CREATE INDEX IF NOT EXISTS %s_idx_embedding ON %s USING vss(embedding) WITH (metric='cosine')`, warehouseTableName, warehouseTableName))
	if err != nil {
		log.Errorf("create index %s_idx_embedding error:%s", warehouseTableName, err.Error())
		return err
	}
	_, err = s.dbInstancePtr.Exec(fmt.Sprintf(`CREATE INDEX IF NOT EXISTS %s_idx_filePath ON %s(filePath)`, warehouseTableName, warehouseTableName))
	if err != nil {
		log.Errorf("create index %s_idx_filePath error:%s", warehouseTableName, err.Error())
		return err
	}

	return nil
}

func (s *KnowledgeBase) DropWarehouse(warehouse string) error {
	warehouseTableName := fmt.Sprintf("tbl_%s", warehouse)
	_, err := s.dbInstancePtr.Exec(fmt.Sprintf(`DROP TABLE IF EXISTS %s`, warehouseTableName))
	if err != nil {
		log.Errorf("drop table %s error:%s", warehouseTableName, err.Error())
		return err
	}

	return nil
}

func (s *KnowledgeBase) CheckWarehouse(warehouse string) error {
	warehouseTableName := fmt.Sprintf("tbl_%s", warehouse)
	_, err := s.dbInstancePtr.Exec(fmt.Sprintf(`SELECT * FROM %s LIMIT 1`, warehouseTableName))
	if err != nil {
		log.Errorf("check table %s error:%s", warehouseTableName, err.Error())
		return err
	}
	return nil
}

func (s *KnowledgeBase) ListWarehouse() (ret []string, err error) {
	return
}

func (s *KnowledgeBase) AddChunk(warehouse string, chunkPtr *imkkc.Chunk, embedTextWorker iiec.EmbedTextWorker) error {
	warehouseTableName := fmt.Sprintf("tbl_%s", warehouse)
	meta, _ := json.Marshal(chunkPtr.Metadata)
	vectorVal, vectorErr := embedTextWorker.EmbedText(chunkPtr.Text)
	if vectorErr != nil {
		log.Errorf("embed text failed, text:%s, err:%v", chunkPtr.Text, vectorErr)
		return vectorErr
	}
	vecLit := floatSliceToArrayLiteral(vectorVal)
	// insert: pass text and metadata as parameters, vec inline to avoid driver binding issues
	sqlStmt := fmt.Sprintf(`INSERT INTO %s (text, metadata, filePath, embedding) VALUES (?, ?, ?, %s)`, warehouseTableName, vecLit)
	_, err := s.dbInstancePtr.Exec(sqlStmt, chunkPtr.Text, string(meta), chunkPtr.FilePath)
	return err
}

func (s *KnowledgeBase) AddChunksBulk(warehouse string, chunkPtrs []*imkkc.Chunk, embedTextWorker iiec.EmbedTextWorker) (err error) {
	warehouseTableName := fmt.Sprintf("tbl_%s", warehouse)
	txVal, txErr := s.dbInstancePtr.Begin()
	if txErr != nil {
		err = txErr
		log.Errorf("BeginTransaction failed, s.dbInstancePtr.Begin error:%s", txErr.Error())
		return
	}
	defer func() {
		if err == nil {
			err = txVal.Commit()
			return
		}
		err = txVal.Rollback()
	}()

	stmtVal, stmtErr := txVal.Prepare(fmt.Sprintf(`INSERT INTO %s (text, metadata, filePath, embedding) VALUES (?, ?, ?, ?)`, warehouseTableName))
	if stmtErr != nil {
		err = stmtErr
		return
	}
	defer stmtVal.Close()
	for _, chunk := range chunkPtrs {
		meta, _ := json.Marshal(chunk.Metadata)
		vectorVal, vectorErr := embedTextWorker.EmbedText(chunk.Text)
		if vectorErr != nil {
			log.Errorf("embed text failed, text:%s, err:%v", chunk.Text, vectorErr)
			err = vectorErr
			return
		}

		vecLit := floatSliceToArrayLiteral(vectorVal)
		_, stmtErr = stmtVal.Exec(chunk.Text, string(meta), chunk.FilePath, vecLit)
		if stmtErr != nil {
			err = stmtErr
			return
		}
	}

	return
}

func (s *KnowledgeBase) DelChunks(warehouse, filePath string) (err error) {
	warehouseTableName := fmt.Sprintf("tbl_%s", warehouse)
	sql := fmt.Sprintf("DELETE FROM %s WHERE filePath = ?", warehouseTableName)
	_, err = s.dbInstancePtr.Exec(sql, filePath)
	return
}

func (s *KnowledgeBase) SearchByVector(warehouse string, vector []float32, topK int) ([]*imkkc.Chunk, error) {
	warehouseTableName := fmt.Sprintf("tbl_%s", warehouse)
	vecLit := floatSliceToArrayLiteral(vector)
	sqlStmt := fmt.Sprintf(`
		SELECT id, text, metadata, filePath, distance
		FROM vss_search('%s', 'embedding', %s, %d);
	`, warehouseTableName, vecLit, topK)
	rows, err := s.dbInstancePtr.Query(sqlStmt)
	if err != nil {
		log.Errorf("query table %s records failed, error:%s", warehouseTableName, err.Error())
		return nil, err
	}
	defer rows.Close()
	var ret []*imkkc.Chunk
	for rows.Next() {
		var id int64
		var text string
		var metadata string
		var filePath string
		var distance float32
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
			Vector:   vector,
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

func (s *KnowledgeBase) SearchByText(warehouse, text string, topK int, embedTextWorker iiec.EmbedTextWorker) ([]*imkkc.Chunk, error) {
	vector, err := embedTextWorker.EmbedText(text)
	if err != nil {
		log.Errorf("embed text failed, error:%s", err.Error())
		return nil, err
	}
	return s.SearchByVector(warehouse, vector, topK)
}
