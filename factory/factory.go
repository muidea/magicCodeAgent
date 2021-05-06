package factory

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"text/template"

	fu "github.com/muidea/magicCommon/foundation/util"
	"github.com/muidea/magicToolkit/core/util"
)

const (
	PackageTemplateFile = "package.tpl"
	ServiceTemplateFile = "service.tpl"
	ModelTemplateFile   = "model.tpl"
)

var functionMap template.FuncMap

func init() {
	functionMap = template.FuncMap{}

	RegisterFunc("toLower", util.ToLower)
	RegisterFunc("toTitle", util.ToTitle)
	RegisterFunc("toCamelCase", util.ToCamelCase)
	RegisterFunc("mergeArgs", util.MergeArgs)
	RegisterFunc("extractAction", util.ExtractActionVal)
	RegisterFunc("getActionMethod", util.GetActionMethod)
	RegisterFunc("formatAction", util.FormatAction)
	RegisterFunc("formatParam", util.FormatParam)
	RegisterFunc("formatResult", util.FormatResult)
}

func RegisterFunc(name string, funcPtr interface{}) {
	fu.ValidateFunc(funcPtr)

	_, curOK := functionMap[name]
	if curOK {
		log.Panicf("duplicate function,name:%s", name)
	}

	functionMap[name] = funcPtr
}

type FileFactory struct {
	templatePath string
	rootPath     string
	templateName string
	contentData  map[string]interface{}
}

func NewFileFactory(templatePath, rootPath, templateName string, contentData map[string]interface{}) *FileFactory {
	return &FileFactory{templatePath: templatePath, rootPath: rootPath, templateName: templateName, contentData: contentData}
}

func (s *FileFactory) getTemplateName(relativePath string) string {
	return fmt.Sprintf("%s/%s", relativePath, s.templateName)
}

func (s *FileFactory) getTemplateContent(relativePath string) (ret string, err error) {
	contentVal, contentErr := util.GetTemplateContent(s.templatePath, relativePath, s.templateName)
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
	existFlag := util.Exists(filePath)
	if !existFlag {
		err = os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			return
		}
	}

	fullFilePath := path.Join(filePath, fileName)
	existFlag = util.Exists(fullFilePath)
	if !existFlag || forceFlag {
		func() {
			fileHandle, fileErr := os.Create(fullFilePath)
			if fileErr != nil {
				err = fileErr
				return
			}
			defer fileHandle.Close()

			templateVal, templateErr := s.getTemplateContent(relativePath)
			if templateErr != nil {
				err = templateErr
				return
			}

			templateHandle, templateErr := template.New(s.getTemplateName(relativePath)).Funcs(functionMap).Parse(templateVal)
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
	contentVal, contentErr := util.GetTemplateUnitContent(relativePath, s.templateName)
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
