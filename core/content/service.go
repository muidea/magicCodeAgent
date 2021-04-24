package content

import (
	"fmt"
	"path"
	"strings"

	_ "github.com/muidea/magicToolkit/core/content/client"
	_ "github.com/muidea/magicToolkit/core/content/common"

	"github.com/muidea/magicToolkit/core/util"
	"github.com/muidea/magicToolkit/model"
)

type Service struct {
	templatePath  string
	projectPath   string
	projectName   string
	serviceDefine *model.Module
}

func NewService(templatePath, projectPath, projectName string, serviceDefine *model.Module) *Service {
	return &Service{templatePath: templatePath, projectPath: projectPath, projectName: projectName, serviceDefine: serviceDefine}
}

func (s *Service) getPackageName(relativePath string) string {
	items := strings.Split(relativePath, "/")
	return items[len(items)-1]
}

func (s *Service) Create(relativePath, templateFile string) (err error) {
	packageName := s.getPackageName(relativePath)

	serviceVal := map[string]interface{}{"PackageName": packageName, "ProjectName": s.projectName, "Service": s.serviceDefine}
	factoryPtr := util.NewFileFactory(s.projectPath, templateFile, serviceVal)

	fileName := fmt.Sprintf("%s.go", strings.ToLower(s.serviceDefine.Name))
	err = factoryPtr.Create(relativePath, fileName, true)
	if err != nil {
		return
	}

	modelTemplate := path.Join(s.templatePath, relativePath, util.ModelTemplateFile)
	if !util.Exists(modelTemplate) {
		return
	}

	for _, val := range s.serviceDefine.Child {
		modelPtr := NewModel(s.templatePath, s.projectPath, s.projectName, s.serviceDefine.Name, val)
		err = modelPtr.Create(relativePath, modelTemplate)
		if err != nil {
			return
		}
	}

	return
}
