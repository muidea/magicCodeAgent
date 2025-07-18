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

	normalizedPath := filepath.Clean(workingPath)
	codeRootPath := normalizedPath
	goModPath := filepath.Join(codeRootPath, "go.mod")
	if _, err := os.Stat(goModPath); err != nil {
		return cd.NewError(cd.IllegalParam, "working path must be a go module")
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

	createrHelper := biz.NewHelper(templatePtr.HostRootFS, imbpc.ProjectTemplateBasePath)
	err = createrHelper.Create(&common.CreateParam{
		WorkingPath:        workingPath,
		ProjectName:        imbpcfg.GetName(),
		ProjectDescription: imbpcfg.GetDescription(),
	})
	if err != nil {
		log.Errorf("create project failed, error: %v", err)
		if re != nil {
			re.Set(nil, err)
		}
		return
	}
}

func (s *Project) refresh(ev event.Event, re event.Result) {

}

func (s *Project) destroy(ev event.Event, re event.Result) {

}

func (s *Project) initializeService(projectName, projectPath string) (err error) {
	rawPath := path.Join(projectPath, "application", "main.go")
	dstPath := path.Join(projectPath, "application", projectName, "cmd")
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
	modInitCmd.Dir = projectPath
	if err = modInitCmd.Run(); err != nil {
		log.Errorf("go mod init failed, projectName:%s, \nstdErr:%v, \nerr:%s", projectName, byteBuffer.String(), err.Error())
		return
	}

	modTidyCmd := exec.Command("go", "mod", "tidy")
	modTidyCmd.Dir = projectPath
	if err = modTidyCmd.Run(); err != nil {
		log.Errorf("go mod tidy failed, projectName:%s, \nstdErr:%v, \nerr:%s", projectName, byteBuffer.String(), err.Error())
		return
	}

	modVendorCmd := exec.Command("go", "mod", "vendor")
	modVendorCmd.Dir = projectPath
	if err = modVendorCmd.Run(); err != nil {
		log.Errorf("go mod vendor failed, projectName:%s, \nstdErr:%v, \nerr:%s", projectName, byteBuffer.String(), err.Error())
		return
	}
	return
}

func (s *Project) initializeWeb(projectName, projectPath string) (err error) {
	var byteBuffer bytes.Buffer
	mvUmiRCCmd := exec.Command("mv", "umirc.ts", ".umirc.ts")
	mvUmiRCCmd.Dir = projectPath
	mvUmiRCCmd.Stderr = &byteBuffer
	if err = mvUmiRCCmd.Run(); err != nil {
		log.Errorf("mv umirc.ts failed, projectName:%s, \nstdErr:%v, \nerr:%s", projectName, byteBuffer.String(), err.Error())
		return
	}

	mvPrettierCmd := exec.Command("mv", "prettierignore", ".prettierignore")
	mvPrettierCmd.Dir = projectPath
	mvPrettierCmd.Stderr = &byteBuffer
	if err = mvPrettierCmd.Run(); err != nil {
		log.Errorf("mv umirc.ts failed, projectName:%s, \nstdErr:%v, \nerr:%s", projectName, byteBuffer.String(), err.Error())
		return
	}

	yarnInstallCmd := exec.Command("yarn", "install")
	yarnInstallCmd.Dir = projectPath
	yarnInstallCmd.Stderr = &byteBuffer
	if err = yarnInstallCmd.Run(); err != nil {
		log.Errorf("yarn failed, projectName:%s, \nstdErr:%v, \nerr:%s", projectName, byteBuffer.String(), err.Error())
		return
	}

	return
}
