package dao

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/marcboeker/go-duckdb/v2"

	"github.com/muidea/magicCommon/foundation/log"
	fp "github.com/muidea/magicCommon/foundation/path"

	"github.com/muidea/magicCodeAgent/internal/config"
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
			id INTEGER,
			embedding FLOAT[],
			content TEXT,
			tags TEXT[]
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

func (s *KnowledgeBase) Load(workspaceName, content string, tags []string) error {
	return nil
}

func (s *KnowledgeBase) Insert(workspaceName string, embedding []float32, content string, tags []string) error {
	tableName := fmt.Sprintf("tbl_%s", workspaceName)
	_, err := s.dbInstancePtr.Exec(fmt.Sprintf(`
		INSERT INTO %s(id, embedding, content, tags) VALUES(NULL, $1, $2, $3);
	`, tableName), embedding, content, tags)
	return err
}
