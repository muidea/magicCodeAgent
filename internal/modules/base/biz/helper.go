package biz

import (
	"io/fs"
	"path"
	"strings"

	"github.com/muidea/magicCodeAgent/internal/config"
	tc "github.com/muidea/magicCodeAgent/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/factory"
	"github.com/muidea/magicCodeAgent/pkg/util"
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/foundation/log"
)

type Helper struct {
	templateBasePath string

	workingPath        string
	projectName        string
	projectDescription string
	moduleName         string
	modulePath         string
	enableService      bool
	enableDao          bool
}

func NewHelper(templateBasePath string) *Helper {
	ret := &Helper{
		templateBasePath: templateBasePath,
	}
	return ret
}

func (s *Helper) Create(prefixPath string, createParam *tc.CreateParam) (err *cd.Error) {
	s.workingPath = createParam.WorkingPath
	s.projectName = createParam.ProjectName
	s.projectDescription = createParam.ProjectDescription
	s.moduleName = util.ToTitle(createParam.ModuleName)
	s.modulePath = path.Join(prefixPath, util.ToLower(createParam.ModuleName))
	s.enableService = createParam.EnableService
	s.enableDao = createParam.EnableDao

	fs.WalkDir(config.GetTemplateFS(), s.templateBasePath, func(filePath string, fileEntry fs.DirEntry, err error) error {
		if s.templateBasePath == filePath || fileEntry.IsDir() || err != nil {
			return nil
		}

		relativePath, fileName := path.Split(filePath[len(s.templateBasePath)+1:])
		err = s.cloneModuleContent(s.templateBasePath, relativePath, fileName)
		if err != nil {
			log.Errorf("cloneModuleContent, filePath:%s, fileName:%s err:%v", relativePath, fileName, err)
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
		"ProjectName":        s.projectName,
		"ProjectDescription": s.projectDescription,
		"ModuleName":         s.moduleName,
		"ModulePath":         s.modulePath,
		"EnableService":      s.enableService,
		"EnableDao":          s.enableDao,
	}
	factoryPtr := factory.NewFileFactory(templateBasePath, relativePath, templateFileName, contentVal)

	err = factoryPtr.Create(s.workingPath, path.Join(util.ToLower(s.moduleName), relativePath), fileName, true)
	return
}
