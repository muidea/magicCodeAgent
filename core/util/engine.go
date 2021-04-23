package util

import (
	"fmt"
	"os"
	"os/exec"
	"text/template"
)

const (
	ProjectTemplate = "project"
	UnitTemplate    = "unit"
)

type Factory struct {
	templateName string
	contentData  map[string]interface{}
}

func NewFactory(templateName string, contentData map[string]interface{}) *Factory {
	return &Factory{templateName: templateName, contentData: contentData}
}

func (s *Factory) getTemplateContent(relativePath, templateName string) (ret string, err error) {
	contentVal, contentErr := GetTemplateFileContent(fmt.Sprintf("%s/%s", relativePath, templateName))
	if contentErr != nil {
		err = contentErr
		return
	}

	ret = contentVal
	return
}

func (s *Factory) Create(relativePath, fileName string, forceFlag bool) (err error) {
	filePath := relativePath
	existFlag, existErr := Exists(filePath)
	if existErr != nil {
		err = existErr
		return
	}
	if !existFlag {
		err = os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			return
		}
	}

	filePath = fmt.Sprintf("%s/%s.go", filePath, fileName)
	existFlag, existErr = Exists(filePath)
	if existErr != nil {
		err = fmt.Errorf("unexpect, fileName:%s", filePath)
		return
	}

	if !existFlag || forceFlag {
		func() {
			fileHandle, fileErr := os.Create(filePath)
			if fileErr != nil {
				err = fileErr
				return
			}
			defer fileHandle.Close()

			contentVal, contentErr := s.getTemplateContent(relativePath, s.templateName)
			if contentErr != nil {
				err = contentErr
				return
			}

			templateHandle, templateErr := template.New(fileName).Parse(contentVal)
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
