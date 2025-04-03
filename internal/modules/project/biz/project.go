package biz

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/muidea/magicToolkit/pkg/factory"
	"github.com/muidea/magicToolkit/pkg/models"
	"github.com/muidea/magicToolkit/pkg/util"
)

// ProjectHelper core define
type ProjectHelper struct {
	projectName string
	rootPath    string
}

func (s *ProjectHelper) getFile(filePath string) (ret *os.File, err error) {
	existsFlag := util.Exists(filePath)
	if !existsFlag {
		err = fmt.Errorf("file does not exist")
		return
	}

	ret, err = os.Open(filePath)
	return
}

func (s *ProjectHelper) cloneContent(templatePath, filePath, fileName string) (err error) {
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
func (s *ProjectHelper) Create(projectName, projectType string) (ret *models.Define, err error) {
	s.projectName = projectName
	if projectType == "" {
		projectType = util.ProjectDefaultTemplatePath
	}

	err = s.walkTemplate("", projectType)
	if err != nil {
		return
	}

	ret = &models.Define{
		ProjectName: projectName,
		ProjectType: projectType,
		Module:      []*models.Module{},
	}
	return
}

func (s *ProjectHelper) walkTemplate(relativePath, projectType string) (err error) {
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

	fdList, fdErr := os.ReadDir(fullPath)
	if fdErr != nil {
		err = fdErr
		return
	}

	for _, fd := range fdList {
		fileName := fd.Name()
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
