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

	imkkc "github.com/muidea/magicCodeAgent/internal/modules/kernel/knowledgebase/pkg/common"
)

type KnowledgeBase struct {
	dbInstancePtr *sql.DB
}

func New() *KnowledgeBase {
	return &KnowledgeBase{}
}

func (s *KnowledgeBase) Open() error {
	dataPath := filepath.Join(config.RootConfigDir(), imkkc.BasePath, imkkc.DataPath)
	if !fp.Exist(dataPath) {
		os.MkdirAll(dataPath, 0755)
	}
	dsn := fmt.Sprintf("%s/knowledge.db?access_mode=read_write", dataPath)
	dbPtr, err := sql.Open("duckdb", dsn)
	if err != nil {
		log.Errorf("open knowledge db failed, dataPath:%s, err:%v", dataPath, err)
		return err
	}

	_, err = dbPtr.Exec(`INSTALL vss; LOAD vss;`)
	if err != nil {
		log.Errorf("init knowledge db failed, err:%v", err)
		return err
	}

	err = s.init(dbPtr)
	if err != nil {
		log.Errorf("init knowledge db failed, err:%v", err)
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
