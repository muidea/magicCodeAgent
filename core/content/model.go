package content

import (
	"fmt"
	"github.com/muidea/magicToolkit/core/util"
	"github.com/muidea/magicToolkit/model"
	"strings"
)

type Model struct {
	rootPath    string
	modelDefine *model.Define
}

func NewModel(rootPath string, modelDefine *model.Define) *Model {
	return &Model{rootPath: rootPath, modelDefine: modelDefine}
}

func (s *Model) getModelName(relativePath string) string {
	items := strings.Split(relativePath, "/")
	return items[len(items)-1]
}

func (s *Model) Create(relativePath, templateFile string) (err error) {
	packageName := s.getModelName(relativePath)

	fileName := fmt.Sprintf("%s.go", packageName)
	packageVal := map[string]interface{}{"ApiVersion": s.modelDefine.ApiVersion, "ProjectName": s.modelDefine.ProjectName, "ModelName": packageName}
	packageFactory := util.NewFactory(s.rootPath, templateFile, packageVal)
	packageFactory.Create(relativePath, fileName, true)

	return
}
