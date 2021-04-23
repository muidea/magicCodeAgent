package content

import (
	"fmt"
	"github.com/muidea/magicToolkit/core/util"
	"github.com/muidea/magicToolkit/model"
	"strings"
)

type Service struct {
	rootPath    string
	modelDefine *model.Define
}

func NewService(rootPath string, modelDefine *model.Define) *Service {
	return &Service{rootPath: rootPath, modelDefine: modelDefine}
}

func (s *Service) getPackageName(relativePath string) string {
	items := strings.Split(relativePath, "/")
	return items[len(items)-1]
}

func (s *Service) Create(relativePath, templateFile string) (err error) {
	packageName := s.getPackageName(relativePath)

	for _, val := range s.modelDefine.Module {
		contentVal := map[string]interface{}{"ApiVersion": s.modelDefine.ApiVersion, "ProjectName": s.modelDefine.ProjectName, "PackageName": packageName, "Module": val}
		factoryPtr := util.NewFactory(s.rootPath, templateFile, contentVal)

		fileName := fmt.Sprintf("%s.go", strings.ToLower(val.Name))
		err = factoryPtr.Create(relativePath, fileName, true)
		if err != nil {
			return
		}
	}

	return
}
