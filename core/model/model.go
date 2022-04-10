package model

import (
	"fmt"
	"github.com/muidea/magicToolkit/core/model/frontend"
	"io/ioutil"
	"path"

	log "github.com/cihub/seelog"

	"github.com/muidea/magicToolkit/core/model/backend"
	mu "github.com/muidea/magicToolkit/core/model/util"
	"github.com/muidea/magicToolkit/core/util"
	"github.com/muidea/magicToolkit/parser"
)

type Model struct {
	rootPath    string
	projectName string
	bothFlag    bool
}

func New(rootPath string, bothFlag bool) *Model {
	return &Model{rootPath: rootPath, bothFlag: bothFlag}
}

func (s *Model) Create(projectName string) (err error) {
	s.projectName = projectName

	modelPath := path.Join(s.rootPath, "model")
	existFlag := util.Exists(modelPath)
	if !existFlag {
		return
	}

	fdList, fdErr := ioutil.ReadDir(modelPath)
	if fdErr != nil {
		err = fdErr
		return
	}

	for _, fd := range fdList {
		fileName := fd.Name()
		fullPath := path.Join(modelPath, fileName)
		if !util.Exists(fullPath) || mu.IsInnerFile(fileName) {
			continue
		}

		packagePtr, packageErr := parser.Parse(fullPath)
		if packageErr != nil {
			log.Errorf("parse failed, file:%s, err:%s", fullPath, packageErr)
			continue
		}

		backend := backend.New(s.rootPath, s.projectName)
		err := backend.Create(packagePtr)
		if err != nil {
			log.Errorf("create %s backend failed, err:%s", packagePtr.Name, err.Error())
		}

		if s.bothFlag {
			frontendRoot := fmt.Sprintf("%s/web/src", s.rootPath)
			frontend := frontend.New(frontendRoot, s.projectName)
			err = frontend.Create(packagePtr)
			if err != nil {
				log.Errorf("create %s frontend failed, err:%s", packagePtr.Name, err.Error())
			}
		}
	}

	return
}
