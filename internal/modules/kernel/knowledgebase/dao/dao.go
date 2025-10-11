package dao

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/marcboeker/go-duckdb/v2"

	"github.com/muidea/magicCommon/foundation/log"
	fp "github.com/muidea/magicCommon/foundation/path"
)

type KnowledgeBase struct {
	dbInstancePtr *sql.DB
	baseDataPath  string
}

func New(dataPath string) *KnowledgeBase {
	return &KnowledgeBase{
		baseDataPath: dataPath,
	}
}

func (s *KnowledgeBase) Open() error {
	dataPath := s.baseDataPath
	if !fp.Exist(dataPath) {
		os.MkdirAll(dataPath, 0755)
	}
	dbPtr, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		log.Errorf("open knowledge db failed, dataPath:%s, err:%v", dataPath, err)
		return err
	}

	_, err = dbPtr.Exec(`INSTALL vss; LOAD vss;`)
	if err != nil {
		log.Errorf("init knowledge db vss failed, err:%v", err)
		return err
	}

	_, err = dbPtr.Exec(`SET hnsw_enable_experimental_persistence = true;`)
	if err != nil {
		log.Errorf("init knowledge db failed, err:%v", err)
		return err
	}

	dsn := fmt.Sprintf("%s/knowledge.db", dataPath)
	sql := fmt.Sprintf(`ATTACH '%s' AS kb; USE 'kb';`, dsn)
	_, err = dbPtr.Exec(sql)
	if err != nil {
		log.Errorf("attach local knowledge.db failed, error:%s", err.Error())
		return err
	}

	err = s.initKnowledge(dbPtr)
	if err != nil {
		log.Errorf("init knowledge db failed, err:%v", err)
		return err
	}

	err = s.initWarehouse(dbPtr)
	if err != nil {
		log.Errorf("init warehouse db failed, err:%v", err)
		return err
	}

	s.dbInstancePtr = dbPtr
	return nil
}

func (s *KnowledgeBase) Close() {
	if s.dbInstancePtr == nil {
		return
	}
	_ = s.dbInstancePtr.Close()
}
