package biz

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/muidea/magicCommon/foundation/log"
	"github.com/muidea/magicToolkit/pkg/factory"
	"github.com/muidea/magicToolkit/pkg/models"
	"github.com/muidea/magicToolkit/pkg/util"
)

type InitatorHelper struct {
	rootPath    string
	projectName string
	modulePath  string
	moduleName  string
}

func (s *InitatorHelper) Create(projectName, moduleName string) (ret *models.Module, err error) {
	s.projectName = projectName
	s.moduleName = util.ToTitle(moduleName)
	s.modulePath = util.ToLower(moduleName)

	err = s.walkTemplate("")
	if err != nil {
		log.Errorf("walkTemplate failed, err:%s", err.Error())
		return
	}

	ret = &models.Module{
		Name: moduleName,
	}
	return
}

func (s *InitatorHelper) walkTemplate(relativePath string) (err error) {
	baseTemplatePath, baseTemplateErr := util.GetModuleTemplatePath()
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
			err = s.walkTemplate(path.Join(relativePath, fileName))
			if err != nil {
				return
			}

			continue
		}

		items := strings.Split(fileName, ".")
		extName := items[len(items)-1]
		if extName == "template" {
			err = s.cloneModuleContent(baseTemplatePath, relativePath, fd.Name())
			if err != nil {
				return
			}

			continue
		}
	}

	return
}

func (s *InitatorHelper) cloneModuleContent(templatePath, filePath, fileName string) (err error) {
	log.Infof(fileName)

	contentVal := map[string]interface{}{
		"ProjectName": s.projectName,
		"ModuleName":  s.moduleName,
		"ModulePath":  s.modulePath,
	}
	factoryPtr := factory.NewFileFactory(templatePath, s.rootPath, filePath, fileName, contentVal)

	items := strings.Split(fileName, ".")
	fileName = strings.Join(items[:len(items)-1], ".")
	if fileName == "common.go" {
		fileName = util.ToLower(s.moduleName) + ".go"
		err = factoryPtr.Create("common", fileName, true)
		return
	}

	filePath = path.Join("core", "module", s.modulePath, filePath)
	err = factoryPtr.Create(filePath, fileName, true)
	return
}
