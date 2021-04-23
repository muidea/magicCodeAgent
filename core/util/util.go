package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/muidea/magicOrm/provider/remote"
)

func Exists(itemPath string) (bool, error) {
	_, err := os.Stat(itemPath)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return true, err
}

func GetModelName(modelPtr *remote.Object) string {
	items := strings.Split(modelPtr.GetName(), ".")
	return items[1]
}

func GetTemplatePath(typeVal string) (ret string, err error) {
	toolKitPath := os.Getenv("MAGICTOOLKIT_PATH")
	if toolKitPath == "" {
		err = fmt.Errorf("no found toolkitPath")
		return
	}

	templatePath := path.Join(toolKitPath, "template", typeVal)
	existsFlag, existsErr := Exists(templatePath)
	if !existsFlag || existsErr != nil {
		err = fmt.Errorf("toolkit template does not exist")
		return
	}

	ret = templatePath
	return
}

func GetTemplateFileContent(filePath string) (ret string, err error) {
	templatePath, templateErr := GetTemplatePath("project")
	if templateErr != nil {
		err = templateErr
		return
	}

	fullFilePath := fmt.Sprintf("%s/%s/model.tpl", templatePath, filePath)
	existsFlag, existsErr := Exists(fullFilePath)
	if !existsFlag || existsErr != nil {
		err = fmt.Errorf("template file not exist")
		return
	}

	contentVal, contentErr := ioutil.ReadFile(fullFilePath)
	if contentErr != nil {
		err = contentErr
		return
	}

	ret = string(contentVal)
	return
}
