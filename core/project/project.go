package project

import (
	"fmt"
	"github.com/muidea/magicToolkit/core/content"
	"github.com/muidea/magicToolkit/model"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/muidea/magicOrm/provider/remote"
	"github.com/muidea/magicToolkit/core/util"
)

// Project core define
type Project struct {
	projectName string
	rootPath    string
	modelDefine *model.Define
	objectList  []*remote.Object
}

// New new core
func New(rootPath string, modelDefine *model.Define, objectList []*remote.Object) *Project {
	return &Project{rootPath: rootPath, modelDefine: modelDefine, objectList: objectList}
}

func (s *Project) getFile(filePath string) (ret *os.File, err error) {
	existsFlag := util.Exists(filePath)
	if !existsFlag {
		err = fmt.Errorf("file does not exist")
		return
	}

	ret, err = os.Open(filePath)
	return
}

func (s *Project) cloneContent(filePath, fileName string) (err error) {
	log.Print(fileName)

	contentVal := map[string]interface{}{"ProjectName": s.projectName}
	factoryPtr := util.NewFactory(s.rootPath, fileName, contentVal)

	items := strings.Split(fileName, ".")
	fileName = strings.Join(items[:len(items)-1], ".")
	if fileName == "main.go" {
		fileName = path.Join("cmd", s.projectName, fileName)
	}

	err = factoryPtr.Create(filePath, fileName, true)
	return
}

// Create create
func (s *Project) Create(projectName string) (err error) {
	s.projectName = projectName

	err = s.walkTemplate("")
	return
}

func (s *Project) walkTemplate(relativePath string) (err error) {
	baseTemplatePath, baseTemplateErr := util.GetTemplatePath(util.ProjectTemplatePath)
	if baseTemplateErr != nil {
		err = baseTemplateErr
		return
	}

	fullPath := path.Join(baseTemplatePath, relativePath)
	existFlag := util.Exists(fullPath)
	if !existFlag {
		err = fmt.Errorf("template file does not exist, fileName:%s", fullPath)
		return
	}

	fdList, fdErr := ioutil.ReadDir(fullPath)
	if fdErr != nil {
		err = fdErr
		return
	}

	for _, fd := range fdList {
		fileName := fd.Name()
		if fd.IsDir() {
			err = s.walkTemplate(path.Join(relativePath, fileName))
			if err != nil {
				return
			}

			continue
		}

		items := strings.Split(fileName, ".")
		extName := items[len(items)-1]
		if extName == "template" {
			err = s.cloneContent(relativePath, fd.Name())
			if err != nil {
				return
			}

			continue
		}

		if fileName == "package.tpl" {
			packagePtr := content.NewPackage(baseTemplatePath, s.rootPath, s.projectName, s.modelDefine.Module)
			packagePtr.Create(relativePath, fd.Name())
			continue
		}
	}

	return
}
