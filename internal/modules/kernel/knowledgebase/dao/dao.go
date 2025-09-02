package dao

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/marcboeker/go-duckdb/v2"

	"github.com/muidea/magicCommon/foundation/log"
	fp "github.com/muidea/magicCommon/foundation/path"

	"github.com/muidea/magicCodeAgent/internal/config"

	imkkc "github.com/muidea/magicCodeAgent/internal/modules/kernel/knowledgebase/pkg/common"
)

type KnowledgeBase struct {
	dbInstancePtr *sql.DB
}

func New() *KnowledgeBase {
	return &KnowledgeBase{}
}

func (s *KnowledgeBase) Open() error {
	dataPath := filepath.Join(config.RootConfigDir(), "data")
	if !fp.Exist(dataPath) {
		os.MkdirAll(dataPath, 0755)
	}
	dsn := fmt.Sprintf("file:%s/knowledge.db?access_mode=read_write", dataPath)
	db, err := sql.Open("duckdb", dsn)
	if err != nil {
		log.Errorf("open knowledge db failed, err:%v", err)
		return err
	}

	_, err = db.Exec(`INSTALL vss; LOAD vss;`)
	if err != nil {
		log.Errorf("init knowledge db failed, err:%v", err)
		return err
	}

	s.dbInstancePtr = db
	return nil
}

func (s *KnowledgeBase) Close() {
	_ = s.dbInstancePtr.Close()
}

func (s *KnowledgeBase) NewWorkspace(workspaceName string) error {
	tableName := fmt.Sprintf("tbl_%s", workspaceName)
	_, err := s.dbInstancePtr.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s(
			id BIGINT AUTO_INCREMENT,
			text TEXT,
			metadata TEXT,
			embedding FLOAT[]
		);
	`, tableName))
	if err != nil {
		log.Errorf("create table %s error:%s", tableName, err.Error())
		return err
	}

	_, err = s.dbInstancePtr.Exec(fmt.Sprintf(`CREATE INDEX IF NOT EXISTS %s_idx ON %s USING vss(embedding) WITH (metric='cosine')`, tableName, tableName))
	if err != nil {
		log.Errorf("create index %s error:%s", tableName, err.Error())
		return err
	}

	return nil
}

func (s *KnowledgeBase) DropWorkspace(workspaceName string) error {
	tableName := fmt.Sprintf("tbl_%s", workspaceName)
	_, err := s.dbInstancePtr.Exec(fmt.Sprintf(`DROP TABLE IF EXISTS %s`, tableName))
	if err != nil {
		log.Errorf("drop table %s error:%s", tableName, err.Error())
		return err
	}

	return nil
}

func (s *KnowledgeBase) AddChunk(workspaceName string, chunkPtr *imkkc.Chunk) error {
	tableName := fmt.Sprintf("tbl_%s", workspaceName)
	meta, _ := json.Marshal(chunkPtr.Metadata)
	vecLit := floatSliceToArrayLiteral(chunkPtr.Vector)
	// insert: pass text and metadata as parameters, vec inline to avoid driver binding issues
	sqlStmt := fmt.Sprintf(`INSERT INTO %s (text, metadata, embedding) VALUES (?, ?, %s)`, tableName, vecLit)
	_, err := s.dbInstancePtr.Exec(sqlStmt, chunkPtr.Text, string(meta))
	return err
}

func (s *KnowledgeBase) AddChunksBulk(workspaceName string, chunkPtrs []*imkkc.Chunk) error {
	tableName := fmt.Sprintf("tbl_%s", workspaceName)
	tx, err := s.dbInstancePtr.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(fmt.Sprintf(`INSERT INTO %s (text, metadata, embedding) VALUES (?, ?, ?)`, tableName))
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()
	for _, chunk := range chunkPtrs {
		meta, _ := json.Marshal(chunk.Metadata)
		vecLit := floatSliceToArrayLiteral(chunk.Vector)
		_, err := stmt.Exec(chunk.Text, string(meta), vecLit)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (s *KnowledgeBase) SearchByVector(workspaceName string, vector []float32, topK int) ([]*imkkc.Chunk, error) {
	tableName := fmt.Sprintf("tbl_%s", workspaceName)
	vecLit := floatSliceToArrayLiteral(vector)
	sqlStmt := fmt.Sprintf(`
		SELECT id, text, metadata, distance
		FROM vss_search('%s', 'embedding', %s, %d);
	`, tableName, vecLit, topK)
	rows, err := s.dbInstancePtr.Query(sqlStmt)
	if err != nil {
		log.Errorf("query table %s records failed, error:%s", tableName, err.Error())
		return nil, err
	}
	defer rows.Close()
	var ret []*imkkc.Chunk
	for rows.Next() {
		var id int64
		var text string
		var metadata string
		var distance float32
		err = rows.Scan(&id, &text, &metadata, &distance)
		if err != nil {
			log.Errorf("scan row failed, error:%s", err.Error())
			return nil, err
		}
		chunk := &imkkc.Chunk{
			ID:       id,
			Text:     text,
			Metadata: make(map[string]string),
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
