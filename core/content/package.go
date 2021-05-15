package content

import (
	"fmt"
	"github.com/muidea/magicToolkit/factory"
	"path"
	"strings"

	"github.com/muidea/magicToolkit/core/util"
	"github.com/muidea/magicToolkit/model"
)

type Package struct {
	templatePath string
	projectPath  string
	projectName  string
	apiVersion   string
	serviceList  []*model.Module
}

func NewPackage(templatePath, projectPath, projectName, apiVersion string, serviceList []*model.Module) *Package {
	return &Package{templatePath: templatePath, projectPath: projectPath, projectName: projectName, apiVersion: apiVersion, serviceList: serviceList}
}

func (s *Package) getPackageName(relativePath string) string {
	items := strings.Split(relativePath, "/")
	return items[len(items)-1]
}

func (s *Package) Create(relativePath, templateFile string) (err error) {
	packageName := s.getPackageName(relativePath)

	packageFileName := fmt.Sprintf("%s.go", packageName)
	packageVal := map[string]interface{}{"PackageName": packageName, "ProjectName": s.projectName, "ServiceList": s.serviceList}
	packageFactory := factory.NewFileFactory(s.templatePath, s.projectPath, templateFile, packageVal)
	err = packageFactory.Create(relativePath, packageFileName, true)
	if err != nil {
		return
	}

	serviceTemplate := path.Join(s.templatePath, relativePath, factory.ServiceTemplateFile)
	if util.Exists(serviceTemplate) {
		for _, val := range s.serviceList {
			val.ApiVersion = s.apiVersion
			servicePtr := NewService(s.templatePath, s.projectPath, s.projectName, val)
			err = servicePtr.Create(relativePath, factory.ServiceTemplateFile)
			if err != nil {
				return
			}
		}
	}

	modelTemplate := path.Join(s.templatePath, relativePath, factory.ModelTemplateFile)
	if util.Exists(modelTemplate) {
		for _, val := range s.serviceList {
			serviceModel := false
			val.ApiVersion = s.apiVersion
			for _, sv := range val.Child {
				if sv.Name == val.Name {
					serviceModel = true
				}

				sv.ApiVersion = s.apiVersion
				modelPtr := NewModel(s.templatePath, s.projectPath, s.projectName, val.Name, sv)
				err = modelPtr.Create(relativePath, factory.ModelTemplateFile)
				if err != nil {
					return
				}
			}

			// if no service model, special create it
			if !serviceModel {
				modelPtr := NewModel(s.templatePath, s.projectPath, s.projectName, val.Name, val)
				err = modelPtr.Create(relativePath, factory.ModelTemplateFile)
				if err != nil {
					return
				}
			}
		}
	}
	return
}
