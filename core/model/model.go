package model

import (
	"io/ioutil"
	"path"

	log "github.com/cihub/seelog"

	"github.com/muidea/magicToolkit/core/module"
	"github.com/muidea/magicToolkit/core/util"
	"github.com/muidea/magicToolkit/parser"
)

type Model struct {
	rootPath    string
	projectName string
}

func New(rootPath string) *Model {
	return &Model{rootPath: rootPath}
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
		if !util.Exists(fullPath) || fileName == "log.go" || fileName == "totalizer.go" || fileName == "helper.go" {
			continue
		}

		packagePtr, packageErr := parser.Parse(fullPath)
		if packageErr != nil {
			log.Errorf("parse failed, file:%s, err:%s", fullPath, packageErr)
		}

		modulePtr := module.New(s.rootPath)
		modulePtr.Create(s.projectName, util.ToTitle(packagePtr.FileName))
	}

	return
}
