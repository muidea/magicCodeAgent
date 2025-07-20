package biz

import (
	"bytes"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/log"
	fp "github.com/muidea/magicCommon/foundation/path"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicCodeAgent/internal/config"
	"github.com/muidea/magicCodeAgent/internal/modules/base/biz"
	imbpcfg "github.com/muidea/magicCodeAgent/internal/modules/blocks/project/config"
	imbpc "github.com/muidea/magicCodeAgent/internal/modules/blocks/project/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
)

type Project struct {
	biz.Base
}

func New(
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine,
) *Project {
	ret := &Project{
		Base: biz.New(imbpc.ProjectModule, eventHub, backgroundRoutine),
	}

	ret.SubscribeFunc(imbpc.CreateNewProject, ret.create)
	ret.SubscribeFunc(imbpc.RefreshProject, ret.refresh)
	ret.SubscribeFunc(imbpc.DestroyProject, ret.destroy)

	return ret
}

func (s *Project) verifyWorkingPath(workingPath string) *cd.Error {
	if workingPath == "" {
		return cd.NewError(cd.IllegalParam, "working path is empty")
	}

	if !filepath.IsAbs(workingPath) {
		return cd.NewError(cd.IllegalParam, "working path must be absolute path")
	}

	if vErr := imbpcfg.Validate(); vErr != nil {
		return cd.NewError(cd.IllegalParam, vErr.Error())
	}

	normalizedPath := filepath.Clean(workingPath)
	codeRootPath := normalizedPath
	if isEmpty, isErr := fp.IsDirEmpty(codeRootPath); isErr != nil || !isEmpty {
		return cd.NewError(cd.IllegalParam, "working path must be empty directory")
	}
	return nil
}

func (s *Project) create(ev event.Event, re event.Result) {
	_, workspaceOK := ev.Data().(*models.Workspace)
	if !workspaceOK {
		log.Errorf("create project failed, illegal create project param")
		if re != nil {
			re.Set(nil, cd.NewError(cd.Unexpected, "illegal create project param"))
		}
		return
	}
	templatePtr := ev.GetData("template").(*models.Template)

	workingPath := config.CurrentWorkingPath()
	err := s.verifyWorkingPath(workingPath)
	if err != nil {
		log.Errorf("create project failed, verify project working path failed, err: %v", err)
		if re != nil {
			re.Set(nil, err)
		}
		return
	}

	projectTemplateBasePath := path.Join(models.TemplatePath, templatePtr.Category, templatePtr.Name, "template", "project")
	createrHelper := biz.NewHelper(templatePtr.HostRootFS, projectTemplateBasePath)
	err = createrHelper.Create(&common.CreateParam{
		WorkingPath:        workingPath,
		ProjectName:        imbpcfg.GetName(),
		ProjectDescription: imbpcfg.GetDescription(),
	})
	if err != nil {
		log.Errorf("create project failed, error: %v", err)
	}

	if templatePtr.IsService() {
		if initErr := s.initializeService(imbpcfg.GetName(), workingPath); initErr != nil {
			err = cd.NewError(cd.Unexpected, initErr.Error())
		}
	}
	if templatePtr.IsWeb() {
		if initErr := s.initializeWeb(imbpcfg.GetName(), workingPath); initErr != nil {
			err = cd.NewError(cd.Unexpected, initErr.Error())
		}
	}

	if re != nil {
		re.Set(nil, err)
	}
}

func (s *Project) refresh(ev event.Event, re event.Result) {

}

func (s *Project) destroy(ev event.Event, re event.Result) {

}

func (s *Project) initializeService(projectName, projectWorkingPath string) (err error) {
	rawPath := path.Join(projectWorkingPath, "application", "main.go")
	dstPath := path.Join(projectWorkingPath, "application", projectName, "cmd")
	if !fp.Exist(dstPath) {
		os.MkdirAll(dstPath, 0755)
	}

	dstPath = path.Join(dstPath, "main.go")
	if err = os.Rename(rawPath, dstPath); err != nil {
		log.Errorf("rename file failed, rawPath:%s, dstPath:%s, err:%s", rawPath, dstPath, err.Error())
		return
	}

	var byteBuffer bytes.Buffer
	// 执行 go mod init github.com/muidea/{projectName} & go mod tidy & go mod vendor 完成go mod初始化
	modInitCmd := exec.Command("go", "mod", "init", "github.com/muidea/"+projectName)
	modInitCmd.Dir = projectWorkingPath
	if err = modInitCmd.Run(); err != nil {
		log.Errorf("go mod init failed, projectName:%s, \nstdErr:%v, \nerr:%s", projectName, byteBuffer.String(), err.Error())
		return
	}

	modTidyCmd := exec.Command("go", "mod", "tidy")
	modTidyCmd.Dir = projectWorkingPath
	if err = modTidyCmd.Run(); err != nil {
		log.Errorf("go mod tidy failed, projectName:%s, \nstdErr:%v, \nerr:%s", projectName, byteBuffer.String(), err.Error())
		return
	}

	modVendorCmd := exec.Command("go", "mod", "vendor")
	modVendorCmd.Dir = projectWorkingPath
	if err = modVendorCmd.Run(); err != nil {
		log.Errorf("go mod vendor failed, projectName:%s, \nstdErr:%v, \nerr:%s", projectName, byteBuffer.String(), err.Error())
		return
	}
	return
}

func (s *Project) initializeWeb(projectName, projectWorkingPath string) (err error) {
	var byteBuffer bytes.Buffer
	yarnInstallCmd := exec.Command("yarn", "install")
	yarnInstallCmd.Dir = projectWorkingPath
	yarnInstallCmd.Stderr = &byteBuffer
	if err = yarnInstallCmd.Run(); err != nil {
		log.Errorf("yarn failed, projectName:%s, \nstdErr:%v, \nerr:%s", projectName, byteBuffer.String(), err.Error())
		return
	}

	return
}
