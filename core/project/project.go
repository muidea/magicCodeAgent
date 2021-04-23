package project

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/muidea/magicToolkit/core/util"
)

// Project core define
type Project struct {
	projectName string
	rootPath    string
}

// New new core
func New(rootPath string) *Project {
	return &Project{rootPath: rootPath}
}

func (s *Project) getFile(filePath string) (ret *os.File, err error) {
	existsFlag, existsErr := util.Exists(filePath)
	if !existsFlag || existsErr != nil {
		err = fmt.Errorf("file does not exist")
		return
	}

	ret, err = os.Open(filePath)
	return
}

func (s *Project) cloneTemplate(filePath, fileName string) (err error) {
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
	baseTemplatePath, baseTemplateErr := util.GetTemplatePath(util.ProjectTemplate)
	if baseTemplateErr != nil {
		err = baseTemplateErr
		return
	}

	fullPath := path.Join(baseTemplatePath, relativePath)
	existFlag, existErr := util.Exists(fullPath)
	if !existFlag || existErr != nil {
		err = fmt.Errorf("template file does not exist, fileName:%s", fullPath)
		return
	}

	fdList, fdErr := ioutil.ReadDir(fullPath)
	if fdErr != nil {
		err = fdErr
		return
	}

	for _, fd := range fdList {
		if fd.IsDir() {
			err = s.walkTemplate(path.Join(relativePath, fd.Name()))
			if err != nil {
				return
			}

			continue
		}

		items := strings.Split(fd.Name(), ".")
		extName := items[len(items)-1]
		if extName == "template" {
			err = s.cloneTemplate(relativePath, fd.Name())
			if err != nil {
				return
			}

			continue
		}
	}

	return
}
