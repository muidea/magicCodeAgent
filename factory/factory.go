package factory

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"text/template"

	log "github.com/cihub/seelog"

	fu "github.com/muidea/magicCommon/foundation/util"
	"github.com/muidea/magicToolkit/core/util"
)

var functionMap template.FuncMap

func init() {
	functionMap = template.FuncMap{}

	RegisterFunc("toLower", util.ToLower)
	RegisterFunc("toTitle", util.ToTitle)
	RegisterFunc("toCamelCase", util.ToCamelCase)
}

func RegisterFunc(name string, funcPtr interface{}) {
	fu.ValidateFunc(funcPtr)

	_, curOK := functionMap[name]
	if curOK {
		log.Criticalf("duplicate function,name:%s", name)
	}

	functionMap[name] = funcPtr
}

func UnRegisterFunc(name string) {
	delete(functionMap, name)
}

type FileFactory struct {
	templatePath string
	rootPath     string
	relativePath string
	templateName string
	contentData  map[string]interface{}
}

func NewFileFactory(templatePath, rootPath, relativePath, templateName string, contentData map[string]interface{}) *FileFactory {
	return &FileFactory{templatePath: templatePath, rootPath: rootPath, relativePath: relativePath, templateName: templateName, contentData: contentData}
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

			templateVal, templateErr := s.getTemplateContent(s.relativePath)
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

type MemoryFactory struct {
	templateVal string
}

func NewMemoryFactor(templateVal string) *MemoryFactory {
	return &MemoryFactory{templateVal: templateVal}
}

func (s *MemoryFactory) Create(data map[string]interface{}) (ret string, err error) {
	templateHandle, templateErr := template.New("").Funcs(functionMap).Parse(s.templateVal)
	if templateErr != nil {
		err = templateErr
		return
	}

	byteBuff := &bytes.Buffer{}
	err = templateHandle.Execute(byteBuff, data)
	if err != nil {
		return
	}

	ret = byteBuff.String()
	return
}
