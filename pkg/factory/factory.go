package factory

import (
	"bytes"
	"io/fs"
	"os"
	"path"
	"text/template"

	"github.com/muidea/magicCommon/foundation/log"
	fp "github.com/muidea/magicCommon/foundation/path"
	fu "github.com/muidea/magicCommon/foundation/util"

	"github.com/muidea/magicCodeAgent/pkg/util"
)

type FileFactory interface {
	Create(workingPath, relativePath, fileName string, forceFlag bool) (err error)
}

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

type fileFactory struct {
	templateRootFS fs.FS
	templatePath   string
	relativePath   string
	templateName   string
	contentData    map[string]interface{}
}

func NewFileFactory(templateRootFS fs.FS, templatePath, relativePath, templateName string, contentData map[string]interface{}) FileFactory {
	return &fileFactory{
		templateRootFS: templateRootFS,
		templatePath:   templatePath,
		relativePath:   relativePath,
		templateName:   templateName,
		contentData:    contentData,
	}
}

func (s *fileFactory) getTemplateName(relativePath string) string {
	return path.Join(relativePath, s.templateName)
}

func (s *fileFactory) getTemplateContent(relativePath string) (ret string, err error) {
	fullFilePath := path.Join(s.templatePath, relativePath, s.templateName)
	if _, fileErr := fs.Stat(s.templateRootFS, fullFilePath); fileErr != nil {
		err = fileErr
		return
	}

	contentVal, contentErr := fs.ReadFile(s.templateRootFS, fullFilePath)
	if contentErr != nil {
		err = contentErr
		return
	}

	ret = string(contentVal)
	return
}

func (s *fileFactory) Create(workingPath, relativePath, fileName string, forceFlag bool) (err error) {
	filePath := path.Join(workingPath, relativePath, fileName)
	filePath, fileName = path.Split(filePath)
	existFlag := fp.Exist(filePath)
	if !existFlag {
		err = os.MkdirAll(filePath, 0755)
		if err != nil {
			return
		}
	}

	fullFilePath := path.Join(filePath, fileName)
	existFlag = fp.Exist(fullFilePath)
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
	}

	return
}

type MemoryFactory struct {
	templateVal string
}

func NewMemoryFactory(templateVal string) *MemoryFactory {
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
