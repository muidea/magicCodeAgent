package project

import (
	"fmt"
	"github.com/muidea/magicToolkit/model"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/muidea/magicToolkit/core/util"
	"github.com/muidea/magicToolkit/factory"
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
	existsFlag := util.Exists(filePath)
	if !existsFlag {
		err = fmt.Errorf("file does not exist")
		return
	}

	ret, err = os.Open(filePath)
	return
}

func (s *Project) cloneContent(templatePath, filePath, fileName string) (err error) {
	log.Print(fileName)

	contentVal := map[string]interface{}{"ProjectName": s.projectName}
	factoryPtr := factory.NewFileFactory(templatePath, s.rootPath, filePath, fileName, contentVal)

	items := strings.Split(fileName, ".")
	fileName = strings.Join(items[:len(items)-1], ".")
	if fileName == "main.go" {
		fileName = path.Join("cmd", s.projectName, fileName)
	}

	err = factoryPtr.Create(filePath, fileName, true)
	return
}

// Create create
func (s *Project) Create(projectName, projectType string) (ret *model.Define, err error) {
	s.projectName = projectName

	err = s.walkTemplate("", projectType)
	if err != nil {
		return
	}

	ret = &model.Define{
		ProjectName: projectName,
		ProjectType: projectType,
		Module:      []*model.Module{},
	}
	return
}

func (s *Project) walkTemplate(relativePath, projectType string) (err error) {
	if projectType == "" {
		projectType = util.ProjectDefaultTemplatePath
	}

	baseTemplatePath, baseTemplateErr := util.GetProjectTemplatePath(projectType)
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
		func() {
			if err != nil {
				log.Printf("init project failed, relativePath:%s,file:%s", relativePath, fileName)
			}
		}()
		if fd.IsDir() {
			err = s.walkTemplate(path.Join(relativePath, fileName), projectType)
			if err != nil {
				return
			}

			continue
		}

		items := strings.Split(fileName, ".")
		extName := items[len(items)-1]
		if extName == "template" {
			err = s.cloneContent(baseTemplatePath, relativePath, fd.Name())
			if err != nil {
				return
			}

			continue
		}
	}

	return
}
