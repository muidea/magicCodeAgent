package content

import (
	"fmt"
	"strings"

	"github.com/muidea/magicToolkit/core/util"
	"github.com/muidea/magicToolkit/model"
)

type Service struct {
	rootPath      string
	projectName   string
	serviceDefine *model.Module
}

func NewService(rootPath, projectName string, serviceDefine *model.Module) *Service {
	return &Service{rootPath: rootPath, projectName: projectName, serviceDefine: serviceDefine}
}

func (s *Service) getPackageName(relativePath string) string {
	items := strings.Split(relativePath, "/")
	return items[len(items)-1]
}

func (s *Service) Create(relativePath, templateFile string) (err error) {
	packageName := s.getPackageName(relativePath)

	serviceVal := map[string]interface{}{"PackageName": packageName, "ProjectName": s.projectName, "Service": s.serviceDefine}
	factoryPtr := util.NewFactory(s.rootPath, templateFile, serviceVal)

	fileName := fmt.Sprintf("%s.go", strings.ToLower(s.serviceDefine.Name))
	err = factoryPtr.Create(relativePath, fileName, true)
	if err != nil {
		return
	}

	return
}
