package content

import (
	"fmt"
	"strings"

	"github.com/muidea/magicToolkit/core/util"
	"github.com/muidea/magicToolkit/model"
)

type Package struct {
	rootPath    string
	projectName string
	serviceList []*model.Module
}

func NewPackage(rootPath, projectName string, serviceList []*model.Module) *Package {
	return &Package{rootPath: rootPath, projectName: projectName, serviceList: serviceList}
}

func (s *Package) getPackageName(relativePath string) string {
	items := strings.Split(relativePath, "/")
	return items[len(items)-1]
}

func (s *Package) Create(relativePath, templateFile string) (err error) {
	packageName := s.getPackageName(relativePath)

	fileName := fmt.Sprintf("%s.go", packageName)
	packageVal := map[string]interface{}{"PackageName": packageName, "ProjectName": s.projectName, "ServiceList": s.serviceList}
	packageFactory := util.NewFactory(s.rootPath, templateFile, packageVal)
	packageFactory.Create(relativePath, fileName, true)

	return
}
