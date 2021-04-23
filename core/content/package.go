package content

import (
	"fmt"
	"github.com/muidea/magicToolkit/core/util"
	"github.com/muidea/magicToolkit/model"
	"strings"
)

type Package struct {
	rootPath    string
	modelDefine *model.Define
}

func NewPackage(rootPath string, modelDefine *model.Define) *Package {
	return &Package{rootPath: rootPath, modelDefine: modelDefine}
}

func (s *Package) getPackageName(relativePath string) string {
	items := strings.Split(relativePath, "/")
	return items[len(items)-1]
}

func (s *Package) Create(relativePath, templateFile string) (err error) {
	packageName := s.getPackageName(relativePath)

	fileName := fmt.Sprintf("%s.go", packageName)
	packageVal := map[string]interface{}{"ApiVersion": s.modelDefine.ApiVersion, "ProjectName": s.modelDefine.ProjectName, "PackageName": packageName}
	packageFactory := util.NewFactory(s.rootPath, templateFile, packageVal)
	packageFactory.Create(relativePath, fileName, true)

	return
}
