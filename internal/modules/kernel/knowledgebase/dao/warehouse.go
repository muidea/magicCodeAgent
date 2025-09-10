package dao

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
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

func (s *KnowledgeBase) HybridSearchByText(warehouse, text string, topK int, threshold float32, embedTextWorker iiec.EmbedTextWorker) ([]*imkkc.Chunk, error) {
	warehouseVectTableName := s.getWarehouseVectTableName(warehouse)
	keywords, err := embedTextWorker.EmbedKeywords(text, topK)
	if err != nil {
		log.Errorf("embed keywords failed: %s", err.Error())
		return nil, err
	}

	vectorWeight, ftsWeight := AutoAdjustWeights(text, ContextMixed, keywords)
	vector, err := embedTextWorker.EmbedText(text)
	if err != nil {
		log.Errorf("embed text failed: %s", err.Error())
		return nil, err
	}
	vecLit := floatSliceToArrayLiteral(vector)

	log.Infof("warehouse:%s, keywords:%v, vector weight:%f, fts weight:%f, topK:%d, threshold:%g", warehouse, keywords, vectorWeight, ftsWeight, topK, threshold)

	vectorSQL := fmt.Sprintf(`
		SELECT id, text, metadata, filePath, array_cosine_distance(embedding, %s) AS dist
		FROM %s
		WHERE array_cosine_distance(embedding, %s) < %g
		LIMIT %d
	`, vecLit, warehouseVectTableName, vecLit, threshold, topK)

	rows, err := s.dbInstancePtr.Query(vectorSQL)
	if err != nil {
		log.Errorf("vector query failed: %s", err.Error())
		return nil, err
	}
	defer rows.Close()

	type weightedChunk struct {
		chunk      *imkkc.Chunk
		vectorDist float32
		ftsScore   float32
		score      float32
	}

	chunkMap := make(map[int64]*weightedChunk)
	for rows.Next() {
		var id int64
		var text, metadata, filePath string
		var distance float64
		if err := rows.Scan(&id, &text, &metadata, &filePath, &distance); err != nil {
			log.Errorf("scan row failed: %s", err.Error())
			return nil, err
		}

		ch := &imkkc.Chunk{
			ID:       id,
			Text:     text,
			FilePath: filePath,
			Distance: distance,
			Metadata: make(map[string]string),
		}
		if err := json.Unmarshal([]byte(metadata), &ch.Metadata); err != nil {
			log.Warnf("unmarshal metadata failed: %s", err.Error())
		}

		chunkMap[id] = &weightedChunk{
			chunk:      ch,
			vectorDist: float32(distance),
			ftsScore:   0,
			score:      0,
		}
	}

	if len(keywords) > 0 {
		likeClauses := make([]string, 0, len(keywords))
		for _, kw := range keywords {
			likeClauses = append(likeClauses, fmt.Sprintf("text LIKE '%%%s%%'", kw))
		}
		ftsSQL := fmt.Sprintf(`
			SELECT id, text, metadata, filePath
			FROM %s
			WHERE %s
		`, warehouseVectTableName, strings.Join(likeClauses, " OR "))

		ftsRows, err := s.dbInstancePtr.Query(ftsSQL)
		if err != nil {
			log.Warnf("FTS query failed: %s", err.Error())
		} else {
			defer ftsRows.Close()
			for ftsRows.Next() {
				var id int64
				var text, metadata, filePath string
				if err := ftsRows.Scan(&id, &text, &metadata, &filePath); err != nil {
					log.Warnf("scan FTS row failed: %s", err.Error())
					continue
				}
				if wc, exists := chunkMap[id]; exists {
					wc.ftsScore = 1.0 // 简单赋值1.0表示匹配
				} else {
					chunkPtr := &imkkc.Chunk{
						ID:       id,
						Text:     text,
						FilePath: filePath,
						Metadata: make(map[string]string),
					}
					if err := json.Unmarshal([]byte(metadata), &chunkPtr.Metadata); err != nil {
						log.Warnf("unmarshal metadata failed: %s", err.Error())
					}
					chunkMap[id] = &weightedChunk{
						chunk:      chunkPtr,
						vectorDist: 1.0,
						ftsScore:   1.0,
						score:      0,
					}
				}
			}
		}
	}

	for _, wc := range chunkMap {
		wc.score = vectorWeight*(1.0-wc.vectorDist) + ftsWeight*wc.ftsScore
	}

	wcList := make([]*weightedChunk, 0, len(chunkMap))
	for _, wc := range chunkMap {
		wcList = append(wcList, wc)
	}

	sort.SliceStable(wcList, func(i, j int) bool {
		return wcList[i].score > wcList[j].score
	})

	finalList := make([]*imkkc.Chunk, 0, topK)
	for i := 0; i < len(wcList) && i < topK; i++ {
		finalList = append(finalList, wcList[i].chunk)
	}

	return finalList, nil
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
