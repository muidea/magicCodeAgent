package views

import (
	"fmt"
	"os"

	"github.com/muidea/magicCommon/foundation/log"

	mu "github.com/muidea/magicToolkit/core/model/util"
	"github.com/muidea/magicToolkit/core/util"
	"github.com/muidea/magicToolkit/parser"
)

type Views struct {
	rootPath    string
	projectName string
}

func New(rootPath, projectName string) *Views {
	return &Views{rootPath: rootPath, projectName: projectName}
}

func (s *Views) Create(ptr *parser.Package) (err error) {
	if len(ptr.Objects) == 0 {
		return
	}

	for _, val := range ptr.Objects {
		commonPath := fmt.Sprintf("%s/pages/%s", s.rootPath, mu.FormatObjectActionValue(ptr, val))
		if !util.Exists(commonPath) {
			err = os.MkdirAll(commonPath, os.ModePerm)
			if err != nil {
				return
			}
		}

		indexer := &indexer{}
		err = indexer.index(ptr, val, commonPath)
		if err != nil {
			log.Errorf("create index.tsx failed, err:%s", err.Error())
		}

		creater := &creater{}
		err = creater.create(ptr, val, commonPath)
		if err != nil {
			log.Errorf("create create.tsx failed, err:%s", err.Error())
		}

		updater := &updater{}
		err = updater.update(ptr, val, commonPath)
		if err != nil {
			log.Errorf("create update.tsx failed, err:%s", err.Error())
		}
	}
	return
}
