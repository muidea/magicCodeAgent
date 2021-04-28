package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/muidea/magicBatis/def"
	commonDef "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/foundation/net"
	"github.com/muidea/magicToolkit/model"

	"github.com/muidea/magicOrm/provider/remote"
)

const (
	filterAction = "filter"
	queryAction  = "query"
	createAction = "create"
	updateAction = "update"
	deleteAction = "delete"
)

func init() {
	RegisterFunc("toLower", ToLower)
	RegisterFunc("toTitle", ToTitle)
	RegisterFunc("extractAction", ExtractActionVal)
	RegisterFunc("formatAction", FormatAction)
	RegisterFunc("formatParam", FormatParam)
	RegisterFunc("formatResult", FormatResult)
}

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
	return strings.Title(val)
}

func validAction(actionUrl string) ([]string, string) {
	items := strings.Split(actionUrl, "/")
	len := len(items)
	idVal := items[len-1]
	actionVal := items[len-2]
	others := []string{}
	if len >= 3 {
		others = items[:len-2]
	}

	if actionVal == "query" && idVal == "" {
		actionVal = "filter"
	}

	return others, actionVal
}

func ExtractActionVal(actionUrl string) string {
	_, val := validAction(actionUrl)
	return val
}

func FormatAction(serviceName, modelName, actionUrl string) string {
	others, action := validAction(actionUrl)
	items := []string{strings.Title(action)}
	items = append(items, strings.Title(serviceName))
	if serviceName != modelName {
		items = append(items, strings.Title(modelName))
	}
	for _, val := range others {
		items = append(items, strings.Title(val))
	}

	return strings.Join(items, "")
}

func FormatParam(serviceName, modelName, actionUrl string) string {
	_, action := validAction(actionUrl)
	switch action {
	case filterAction:
		return "filter *def.QueryFilter"
	case queryAction:
		return "id int"
	case createAction:
		return fmt.Sprintf("ptr *model.%s", modelName)
	case updateAction:
		return fmt.Sprintf("ptr *model.%s", modelName)
	case deleteAction:
		return "id int"
	default:
	}

	return ""
}

func FormatResult(serviceName, modelName, actionUrl string) string {
	_, action := validAction(actionUrl)
	switch action {
	case filterAction:
		return fmt.Sprintf("ret []*model.%s,err error", modelName)
	case queryAction:
		return fmt.Sprintf("ret *model.%s,err error", modelName)
	case createAction:
		return fmt.Sprintf("ret *model.%s,err error", modelName)
	case updateAction:
		return fmt.Sprintf("ret *model.%s,err error", modelName)
	case deleteAction:
		return fmt.Sprintf("ret *model.%s,err error", modelName)
	default:
	}

	return ""
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
	existsFlag := Exists(templatePath)
	if !existsFlag {
		err = fmt.Errorf("toolkit template does not exist")
		return
	}

	ret = templatePath
	return
}

func getTemplateContent(templatePath, filePath, fileName string) (ret string, err error) {
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

	contentVal, contentErr := ioutil.ReadFile(fullFilePath)
	if contentErr != nil {
		err = contentErr
		return
	}

	ret = string(contentVal)
	return
}

func GetTemplateFileContent(filePath, fileName string) (ret string, err error) {
	templatePath, templateErr := GetTemplatePath(ProjectTemplatePath)
	if templateErr != nil {
		err = templateErr
		return
	}

	ret, err = getTemplateContent(templatePath, filePath, fileName)
	return
}

func GetTemplateUnitContent(filePath, unitName string) (ret string, err error) {
	templatePath, templateErr := GetTemplatePath(UnitTemplatePath)
	if templateErr != nil {
		err = templateErr
		return
	}

	ret, err = getTemplateContent(templatePath, filePath, unitName)
	return
}

func GetModelDefine(projectName, definePath string) (ret *model.Define, err error) {
	defineFile, defineErr := os.Open(definePath)
	if defineErr != nil {
		return
	}
	defer defineFile.Close()

	byteContent, byteErr := ioutil.ReadAll(defineFile)
	if byteErr != nil {
		err = byteErr
		return
	}

	contentDefine := &model.Define{}
	err = json.Unmarshal(byteContent, contentDefine)
	if err != nil {
		return
	}

	contentDefine.ProjectName = projectName

	ret = contentDefine
	return
}

func QueryObject(projectName, batisService string) (ret []*remote.Object, err error) {
	result := &def.QueryModelResult{}
	url := strings.Join([]string{batisService, "/model/query/"}, "")
	url = fmt.Sprintf("%s?source=%s", url, projectName)

	httpClient := &http.Client{}
	_, err = net.HTTPGet(httpClient, url, result)
	if err != nil {
		return
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("query project objects failed, reason:%s", result.Reason)
		return
	}

	ret = result.Model
	return
}
