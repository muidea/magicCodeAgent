package biz

import (
	"os"
	"os/exec"
	"path"
	"strings"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/log"
	fp "github.com/muidea/magicCommon/foundation/path"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicCodeAgent/internal/modules/base/biz"
	"github.com/muidea/magicCodeAgent/internal/modules/blocks/project/pkg/common"
	tc "github.com/muidea/magicCodeAgent/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/util"
)

type Project struct {
	biz.Base
}

func New(
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine,
) *Project {
	ret := &Project{
		Base: biz.New(common.ProjectModule, eventHub, backgroundRoutine),
	}

	ret.SubscribeFunc(common.CreateNewProject, ret.create)
	ret.SubscribeFunc(common.RefreshProject, ret.refresh)
	ret.SubscribeFunc(common.DestroyProject, ret.destroy)

	return ret
}

func (s *Project) create(ev event.Event, re event.Result) {
	paramPtr, paramOK := ev.Data().(*tc.CreateParam)
	if !paramOK {
		log.Errorf("create project failed, illegal create project param")
		if re != nil {
			re.Set(nil, cd.NewError(cd.Unexpected, "illegal create project param"))
		}
		return
	}

	err := util.VerifyProjectWorkPath(paramPtr.WorkingPath)
	if err != nil {
		log.Errorf("verify project work path failed, workingPath:%s, err:%s", paramPtr.WorkingPath, err.Error())

		if re != nil {
			re.Set(nil, err)
		}
		return
	}
	items := strings.Split(paramPtr.ProjectName, "/")
	projectType := items[0]
	if len(items) != 2 || (projectType != "service" && projectType != "web") {
		log.Errorf("invalid project name, projectName:%s", paramPtr.ProjectName)
		if re != nil {
			re.Set(nil, cd.NewError(cd.InvalidParameter, "invalid project name"))
		}
	}

	projectName := items[1]
	helper := biz.NewHelper(path.Join(common.ProjectTemplateBasePath, projectType))

	// 这里主动修正一下ProjectName，因为在传入参数时的ProprojectName里指定了类型，所以这里需要修正一下
	paramPtr.ProjectName = projectName
	err = helper.Create("", paramPtr)
	if err != nil {
		log.Errorf("create project failed, projectName:%s, err:%s", projectName, err.Error())
		if re != nil {
			re.Set(nil, err)
		}
		return
	}

	if projectType == "service" {
		if err := s.initializeService(projectName, paramPtr.WorkingPath); err != nil {
			log.Errorf("initialize service failed, err:%s", err.Error())
			if re != nil {
				re.Set(nil, cd.NewError(cd.Unexpected, err.Error()))
			}
			return
		}
	}

	if re != nil {
		re.Set(nil, nil)
	}
}

func (s *Project) refresh(ev event.Event, re event.Result) {

}

func (s *Project) destroy(ev event.Event, re event.Result) {

}

func (s *Project) initializeService(projectName, projectPath string) (err error) {
	rawPath := path.Join(projectPath, "application", "main.go")
	dstPath := path.Join(projectPath, "application", projectName)
	if !fp.Exist(dstPath) {
		os.MkdirAll(dstPath, 0755)
	}

	dstPath = path.Join(dstPath, "main.go")
	if err = os.Rename(rawPath, dstPath); err != nil {
		log.Errorf("rename file failed, rawPath:%s, dstPath:%s, err:%s", rawPath, dstPath, err.Error())
		return
	}

	// 执行 go mod init github.com/muidea/{projectName} & go mod tidy & go mod vendor 完成go mod初始化
	modInitCmd := exec.Command("go", "mod", "init", "github.com/muidea/"+projectName)
	modInitCmd.Dir = projectPath
	if err = modInitCmd.Run(); err != nil {
		log.Errorf("go mod init failed, projectName:%s, err:%s", projectName, err.Error())
		return
	}

	modTidyCmd := exec.Command("go", "mod", "tidy")
	modTidyCmd.Dir = projectPath
	if err = modTidyCmd.Run(); err != nil {
		log.Errorf("go mod tidy failed, projectName:%s, err:%s", projectName, err.Error())
		return
	}

	modVendorCmd := exec.Command("go", "mod", "vendor")
	modVendorCmd.Dir = projectPath
	if err = modVendorCmd.Run(); err != nil {
		log.Errorf("go mod vendor failed, projectName:%s, err:%s", projectName, err.Error())
		return
	}
	return
}
