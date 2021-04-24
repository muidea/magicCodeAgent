package util

import (
	"bytes"
	"fmt"
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

type FileFactory struct {
	rootPath     string
	templateName string
	contentData  map[string]interface{}
}

func RegisterFunc(name string, funcPtr interface{}) {
	fu.ValidateFunc(funcPtr)

	functionMap[name] = funcPtr
}

func NewFileFactory(rootPath, templateName string, contentData map[string]interface{}) *FileFactory {
	return &FileFactory{rootPath: rootPath, templateName: templateName, contentData: contentData}
}

func (s *FileFactory) getTemplateName(relativePath string) string {
	return fmt.Sprintf("%s/%s", relativePath, s.templateName)
}

func (s *FileFactory) getTemplateContent(relativePath string) (ret string, err error) {
	contentVal, contentErr := GetTemplateFileContent(relativePath, s.templateName)
	if contentErr != nil {
		err = contentErr
		return
	}

	ret = contentVal
	return
}

func (s *FileFactory) Create(relativePath, fileName string, forceFlag bool) (err error) {
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

			templateHandle, templateErr := template.New(s.getTemplateName(relativePath)).Funcs(functionMap).Parse(contentVal)
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

type ValueFactory struct {
	templateName string
	contentData  map[string]interface{}
}

func NewValueFactory(templateName string, contentData map[string]interface{}) *ValueFactory {
	return &ValueFactory{templateName: templateName, contentData: contentData}
}

func (s *ValueFactory) getTemplateName(relativePath string) string {
	return fmt.Sprintf("%s/%s", relativePath, s.templateName)
}

func (s *ValueFactory) getTemplateContent(relativePath string) (ret string, err error) {
	contentVal, contentErr := GetTemplateFileContent(relativePath, s.templateName)
	if contentErr != nil {
		err = contentErr
		return
	}

	ret = contentVal
	return
}

func (s *ValueFactory) Create(relativePath string) (ret string, err error) {
	byteBuffer := bytes.NewBufferString("")
	contentVal, contentErr := s.getTemplateContent(relativePath)
	if contentErr != nil {
		err = contentErr
		return
	}

	templateHandle, templateErr := template.New(s.getTemplateName(relativePath)).Funcs(functionMap).Parse(contentVal)
	if templateErr != nil {
		err = templateErr
		return
	}

	err = templateHandle.Execute(byteBuffer, s.contentData)
	if err != nil {
		return
	}

	ret = byteBuffer.String()
	return
}
