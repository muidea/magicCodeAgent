package util

import (
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
	"path"
	"strings"
)

const (
	ProjectDefaultTemplatePath = "default"
)

func Exists(itemPath string) bool {
	_, err := os.Stat(itemPath)
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}

	return true
}

func ToLower(val string) string {
	return strings.ToLower(val)
}

func ToTitle(val string) string {
	caser := cases.Title(language.English)
	return caser.String(val)
}

func ToCamelCase(val string) string {
	if len(val) == 1 {
		return strings.ToLower(val)
	}

	return strings.ToLower(val[0:1]) + val[1:]
}

func GetProjectTemplatePath(typeVal string) (ret string, err error) {
	toolKitPath := os.Getenv("MAGICTOOLKIT_PATH")
	if toolKitPath == "" {
		err = fmt.Errorf("no found toolkitPath")
		return
	}

	templatePath := path.Join(toolKitPath, "template", "project", typeVal)
	existsFlag := Exists(templatePath)
	if !existsFlag {
		err = fmt.Errorf("toolkit project template does not exist")
		return
	}

	ret = templatePath
	return
}

func GetModuleTemplatePath() (ret string, err error) {
	toolKitPath := os.Getenv("MAGICTOOLKIT_PATH")
	if toolKitPath == "" {
		err = fmt.Errorf("no found toolkitPath")
		return
	}

	templatePath := path.Join(toolKitPath, "template", "module")
	existsFlag := Exists(templatePath)
	if !existsFlag {
		err = fmt.Errorf("toolkit module template does not exist")
		return
	}

	ret = templatePath
	return
}

func GetTemplateContent(templatePath, filePath, fileName string) (ret string, err error) {
	var fullFilePath string
	if filePath == "" {
		fullFilePath = fmt.Sprintf("%s/%s", templatePath, fileName)
	} else {
		fullFilePath = fmt.Sprintf("%s/%s/%s", templatePath, filePath, fileName)
	}

	existsFlag := Exists(fullFilePath)
	if !existsFlag {
		err = fmt.Errorf("template file not exist, file:%s", fullFilePath)
		return
	}

	contentVal, contentErr := os.ReadFile(fullFilePath)
	if contentErr != nil {
		err = contentErr
		return
	}

	ret = string(contentVal)
	return
}
