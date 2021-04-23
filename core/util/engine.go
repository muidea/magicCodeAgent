package util

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/muidea/magicToolkit/model"
)

const (
	ProjectTemplate = "project"
	UnitTemplate    = "unit"
)

type Factory struct {
	projectName  string
	templateName string
}

func NewFactory(projectName, templateName string) *Factory {
	return &Factory{projectName: projectName, templateName: templateName}
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

func (s *Factory) getPackageName(relativePath string) string {
	items := strings.Split(relativePath, "/")
	return items[len(items)-1]
}

func (s *Factory) Create(define *model.Define, relativePath, fileName string, forceFlag bool) (err error) {
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

			data := map[string]interface{}{"ProjectName": s.projectName, "PackageName": s.getPackageName(relativePath), "Define": define}
			err = templateHandle.Execute(fileHandle, data)
		}()
		if err != nil {
			return
		}

		cmd := exec.Command("gofmt", "-w", filePath)
		err = cmd.Run()
	}
	return
}
