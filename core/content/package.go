package content

import (
	"fmt"
	"path"
	"strings"

	"github.com/muidea/magicToolkit/core/util"
	"github.com/muidea/magicToolkit/model"
)

type Package struct {
	templatePath string
	projectPath  string
	projectName  string
	serviceList  []*model.Module
}

func NewPackage(templatePath, projectPath, projectName string, serviceList []*model.Module) *Package {
	return &Package{templatePath: templatePath, projectPath: projectPath, projectName: projectName, serviceList: serviceList}
}

func (s *Package) getPackageName(relativePath string) string {
	items := strings.Split(relativePath, "/")
	return items[len(items)-1]
}

func (s *Package) Create(relativePath, templateFile string) (err error) {
	packageName := s.getPackageName(relativePath)

	packageFileName := fmt.Sprintf("%s.go", packageName)
	packageVal := map[string]interface{}{"PackageName": packageName, "ProjectName": s.projectName, "ServiceList": s.serviceList}
	packageFactory := util.NewFileFactory(s.projectPath, templateFile, packageVal)
	err = packageFactory.Create(relativePath, packageFileName, true)
	if err != nil {
		return
	}

	serviceTemplate := path.Join(s.templatePath, relativePath, util.ServiceTemplateFile)
	if !util.Exists(serviceTemplate) {
		return
	}

	for _, val := range s.serviceList {
		servicePtr := NewService(s.templatePath, s.projectPath, s.projectName, val)
		err = servicePtr.Create(relativePath, serviceTemplate)
		if err != nil {
			return
		}
	}
	return
}
