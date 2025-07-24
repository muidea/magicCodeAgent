package biz

import (
	"io/fs"
	"path"
	"strings"

	tc "github.com/muidea/magicCodeAgent/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/factory"
	"github.com/muidea/magicCodeAgent/pkg/util"
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/foundation/log"
)

type Helper interface {
	Create(createParam *tc.CreateParam) (err *cd.Error)
}

type helperImpl struct {
	templateRootFS   fs.FS
	templateBasePath string

	workingPath        string
	projectName        string
	projectDescription string
	moduleName         string
	modulePath         string
	initiatorName      string
	initiatorPath      string
}

func NewHelper(rootFS fs.FS, templateBasePath string) Helper {
	ret := &helperImpl{
		templateRootFS:   rootFS,
		templateBasePath: templateBasePath,
	}
	return ret
}

func (s *helperImpl) Create(createParam *tc.CreateParam) (err *cd.Error) {
	s.workingPath = createParam.WorkingPath
	s.projectName = createParam.ProjectName
	s.projectDescription = createParam.ProjectDescription
	s.moduleName = util.ToTitle(createParam.ModuleName)
	s.modulePath = path.Join(createParam.ModulePath, util.ToLower(createParam.ModuleName))
	s.initiatorName = util.ToTitle(createParam.InitiatorName)
	s.initiatorPath = path.Join(createParam.InitiatorPath, util.ToLower(createParam.InitiatorName))

	fs.WalkDir(s.templateRootFS, s.templateBasePath, func(filePath string, fileEntry fs.DirEntry, err error) error {
		if s.templateBasePath == filePath || fileEntry.IsDir() || err != nil {
			return nil
		}

		relativePath, fileName := path.Split(filePath[len(s.templateBasePath)+1:])
		err = s.cloneModuleContent(relativePath, fileName)
		if err != nil {
			log.Errorf("cloneModuleContent, filePath:%s, fileName:%s err:%v", relativePath, fileName, err)
			return err
		}

		return nil
	})

	return
}

func (s *helperImpl) cloneModuleContent(relativePath, templateFileName string) (err error) {
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
		"InitiatorName":      s.initiatorName,
		"InitiatorPath":      s.initiatorPath,
	}
	factoryPtr := factory.NewFileFactory(s.templateRootFS, s.templateBasePath, relativePath, templateFileName, contentVal)

	err = factoryPtr.Create(s.workingPath, relativePath, fileName, true)
	return
}
