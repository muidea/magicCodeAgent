package biz

import (
	"io/fs"
	"path"
	"strings"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicToolkit/internal/config"
	tc "github.com/muidea/magicToolkit/pkg/common"
	"github.com/muidea/magicToolkit/pkg/factory"
	"github.com/muidea/magicToolkit/pkg/util"
)

type Helper struct {
	templateBasePath string

	workingPath   string
	projectName   string
	moduleName    string
	modulePath    string
	enableService bool
	enableDao     bool
}

func NewHelper(templateBasePath string) *Helper {
	ret := &Helper{
		templateBasePath: templateBasePath,
	}
	return ret
}

func (s *Helper) Create(projectName, prefixPath string, moduleParam *tc.CreateParam) (err *cd.Error) {
	s.workingPath = moduleParam.WorkingPath
	s.projectName = projectName
	s.moduleName = util.ToTitle(moduleParam.Name)
	s.modulePath = path.Join(prefixPath, util.ToLower(moduleParam.Name))
	s.enableService = moduleParam.EnableService
	s.enableDao = moduleParam.EnableDao

	fs.WalkDir(config.GetTemplateFS(), s.templateBasePath, func(filePath string, fileEntry fs.DirEntry, err error) error {
		if s.templateBasePath == filePath || fileEntry.IsDir() || err != nil {
			return nil
		}

		relativePath, fileName := path.Split(filePath[len(s.templateBasePath)+1:])
		err = s.cloneModuleContent(s.templateBasePath, relativePath, fileName)
		if err != nil {
			return err
		}

		return nil
	})

	return
}

func (s *Helper) cloneModuleContent(templateBasePath, relativePath, templateFileName string) (err error) {
	items := strings.Split(templateFileName, ".")
	fileName := strings.Join(items[:len(items)-1], ".")
	fileExt := items[len(items)-1]
	if fileExt != "template" {
		return
	}

	contentVal := map[string]interface{}{
		"ProjectName":   s.projectName,
		"ModuleName":    s.moduleName,
		"ModulePath":    s.modulePath,
		"EnableService": s.enableService,
		"EnableDao":     s.enableDao,
	}
	factoryPtr := factory.NewFileFactory(templateBasePath, relativePath, templateFileName, contentVal)

	err = factoryPtr.Create(s.workingPath, relativePath, fileName, true)
	return
}
