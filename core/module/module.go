package module

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	log "github.com/cihub/seelog"

	"github.com/muidea/magicToolkit/core/util"
	"github.com/muidea/magicToolkit/factory"
	"github.com/muidea/magicToolkit/model"
)

type Module struct {
	rootPath    string
	projectName string
	moduleName  string
}

func New(rootPath string) *Module {
	return &Module{rootPath: rootPath}
}

func (s *Module) Create(projectName, moduleName string) (ret *model.Module, err error) {
	s.projectName = projectName
	s.moduleName = moduleName

	err = s.walkTemplate("")
	if err != nil {
		log.Errorf("walkTemplate failed, err:%s", err.Error())
		return
	}

	ret = &model.Module{
		Name: moduleName,
	}
	return
}

func (s *Module) walkTemplate(relativePath string) (err error) {
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
			err = s.cloneModuleContent(baseTemplatePath, relativePath, fd.Name())
			if err != nil {
				return
			}

			continue
		}
	}

	return
}

func (s *Module) cloneModuleContent(templatePath, filePath, fileName string) (err error) {
	log.Infof(fileName)

	contentVal := map[string]interface{}{
		"ProjectName": s.projectName,
		"ModuleName":  s.moduleName,
		"ModulePath":  util.ToCamelCase(s.moduleName),
	}
	factoryPtr := factory.NewFileFactory(templatePath, s.rootPath, filePath, fileName, contentVal)

	items := strings.Split(fileName, ".")
	fileName = strings.Join(items[:len(items)-1], ".")
	if fileName == "common.go" {
		fileName = util.ToCamelCase(s.moduleName) + ".go"
		err = factoryPtr.Create("common", fileName, true)
		return
	}

	filePath = path.Join("core", "module", util.ToCamelCase(s.moduleName), filePath)
	err = factoryPtr.Create(filePath, fileName, true)
	return
}
