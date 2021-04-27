package content

import (
	"fmt"
	"strings"

	"github.com/muidea/magicToolkit/core/util"
	"github.com/muidea/magicToolkit/model"
)

type Model struct {
	templatePath string
	projectPath  string
	projectName  string
	serviceName  string
	modelDefine  *model.Module
}

func NewModel(templatePath, projectPath, projectName, serviceName string, modelDefine *model.Module) *Model {
	return &Model{templatePath: templatePath, projectPath: projectPath, projectName: projectName, serviceName: serviceName, modelDefine: modelDefine}
}

func (s *Model) getPackageName(relativePath string) string {
	items := strings.Split(relativePath, "/")
	return items[len(items)-1]
}

func (s *Model) Create(relativePath, templateFile string) (err error) {
	modelVal := map[string]interface{}{"PackageName": util.ToLower(s.serviceName), "ProjectName": s.projectName, "ServiceName": s.serviceName, "Model": s.modelDefine}
	factoryPtr := util.NewFileFactory(s.templatePath, s.projectPath, templateFile, modelVal)

	fileName := fmt.Sprintf("%s/%s.go", util.ToLower(s.serviceName), strings.ToLower(s.modelDefine.Name))
	err = factoryPtr.Create(relativePath, fileName, true)
	if err != nil {
		return
	}

	return
}
