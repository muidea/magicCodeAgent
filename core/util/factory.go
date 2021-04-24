package util

import (
	"os"
	"os/exec"
	"path"
	"text/template"

	fu "github.com/muidea/magicCommon/foundation/util"
)

const (
	ProjectTemplatePath = "project"
	UnitTemplatePath    = "unit"
)

const (
	PackageTemplateFile = "package.tpl"
	ServiceTemplateFile = "service.tpl"
	ModelTemplateFile   = "model.tpl"
)

var functionMap template.FuncMap

func init() {
	functionMap = template.FuncMap{}
}

type Factory struct {
	rootPath     string
	templateName string
	contentData  map[string]interface{}
}

func RegisterFunc(name string, funcPtr interface{}) {
	fu.ValidateFunc(funcPtr)

	functionMap[name] = funcPtr
}

func NewFactory(rootPath, templateName string, contentData map[string]interface{}) *Factory {
	return &Factory{rootPath: rootPath, templateName: templateName, contentData: contentData}
}

func (s *Factory) getTemplateContent(relativePath string) (ret string, err error) {
	contentVal, contentErr := GetTemplateFileContent(relativePath, s.templateName)
	if contentErr != nil {
		err = contentErr
		return
	}

	ret = contentVal
	return
}

func (s *Factory) Create(relativePath, fileName string, forceFlag bool) (err error) {
	filePath := path.Join(s.rootPath, relativePath, fileName)
	filePath, fileName = path.Split(filePath)
	existFlag := Exists(filePath)
	if !existFlag {
		err = os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			return
		}
	}

	fullFilePath := path.Join(filePath, fileName)
	existFlag = Exists(fullFilePath)
	if !existFlag || forceFlag {
		func() {
			fileHandle, fileErr := os.Create(fullFilePath)
			if fileErr != nil {
				err = fileErr
				return
			}
			defer fileHandle.Close()

			contentVal, contentErr := s.getTemplateContent(relativePath)
			if contentErr != nil {
				err = contentErr
				return
			}

			templateHandle, templateErr := template.New(fileName).Funcs(functionMap).Parse(contentVal)
			if templateErr != nil {
				err = templateErr
				return
			}

			err = templateHandle.Execute(fileHandle, s.contentData)
		}()
		if err != nil {
			return
		}

		cmd := exec.Command("gofmt", "-w", filePath)
		err = cmd.Run()
	}

	return
}
