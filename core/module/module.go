package module

import (
	"fmt"
	"github.com/muidea/magicToolkit/core/util"
	"github.com/muidea/magicToolkit/factory"
	"io/ioutil"
	"log"
	"path"
	"strings"
)

type Module struct {
	rootPath    string
	projectName string
	moduleName  string
}

func New(rootPath string) *Module {
	return &Module{rootPath: rootPath}
}

func (s *Module) Create(projectName, moduleName string) (err error) {
	s.projectName = projectName
	s.moduleName = moduleName
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
		func() {
			if err != nil {
				log.Printf("init project failed, relativePath:%s,file:%s", relativePath, fileName)
			}
		}()
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
			err = s.cloneContent(baseTemplatePath, relativePath, fd.Name())
			if err != nil {
				return
			}

			continue
		}
	}

	return
}

func (s *Module) cloneContent(templatePath, filePath, fileName string) (err error) {
	log.Print(fileName)

	contentVal := map[string]interface{}{
		"ProjectName": s.projectName,
		"ModuleName":  s.moduleName,
		"ModulePath":  util.ToCamelCase(s.moduleName),
	}
	factoryPtr := factory.NewFileFactory(templatePath, s.rootPath, fileName, contentVal)

	items := strings.Split(fileName, ".")
	fileName = strings.Join(items[:len(items)-1], ".")
	if fileName == "common.go" {
		fileName = path.Join("common", util.ToCamelCase(s.moduleName)+".go")
	}

	err = factoryPtr.Create(filePath, fileName, true)
	return
}
