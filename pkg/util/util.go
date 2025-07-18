package util

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	ProjectDefaultTemplatePath = "default"
)

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

func CheckFileExist(templatePath, filePath, fileName string) (ret bool, err error) {
	filePath = filepath.Join(templatePath, filePath)
	filePath = filepath.Join(filePath, fileName)
	_, err = os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func GetTemplateContent(templateFS fs.FS, templatePath, filePath, fileName string) (ret string, err error) {
	fullFilePath := path.Join(templatePath, filePath, fileName)
	if _, fileErr := fs.Stat(templateFS, fullFilePath); fileErr != nil {
		err = fileErr
		return
	}

	contentVal, contentErr := fs.ReadFile(templateFS, fullFilePath)
	if contentErr != nil {
		err = contentErr
		return
	}

	ret = string(contentVal)
	return
}
