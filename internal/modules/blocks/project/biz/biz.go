package biz

import (
	"bytes"
	"os"
	"os/exec"
	"path"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/log"
	fp "github.com/muidea/magicCommon/foundation/path"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicCodeAgent/internal/config"
	"github.com/muidea/magicCodeAgent/internal/modules/base/biz"
	imbpc "github.com/muidea/magicCodeAgent/internal/modules/blocks/project/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/common"
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

func (s *Project) create(ev event.Event, re event.Result) {
	paramPtr, paramOK := ev.Data().(*common.CreateParam)
	if !paramOK {
		log.Errorf("create project failed, illegal create project param")
		if re != nil {
			re.Set(nil, cd.NewError(cd.Unexpected, "illegal create project param"))
		}
		return
	}

	err := config.VerifyProjectWorkPath(paramPtr.ProjectType, paramPtr.WorkingPath)
	if err != nil {
		log.Errorf("verify project work path failed, workingPath:%s, err:%s", paramPtr.WorkingPath, err.Error())

		if re != nil {
			re.Set(nil, err)
		}
		return
	}
	if paramPtr.ProjectType != common.Service && paramPtr.ProjectType != common.Web {
		log.Errorf("invalid project name, projectName:%s", paramPtr.ProjectName)
		if re != nil {
			re.Set(nil, cd.NewError(cd.InvalidParameter, "invalid project name"))
		}
		return
	}

	helper := biz.NewHelper(path.Join(imbpc.ProjectTemplateBasePath, paramPtr.ProjectType))

	// 这里主动修正一下ProjectName，因为在传入参数时的ProprojectName里指定了类型，所以这里需要修正一下
	//paramPtr.ProjectName = projectName
	err = helper.Create("", paramPtr)
	if err != nil {
		log.Errorf("create project failed, projectName:%s, err:%s", paramPtr.ProjectName, err.Error())
		if re != nil {
			re.Set(nil, err)
		}
		return
	}

	if paramPtr.ProjectType == common.Service {
		if err := s.initializeService(paramPtr.ProjectName, paramPtr.WorkingPath); err != nil {
			log.Errorf("initialize service failed, err:%s", err.Error())
			if re != nil {
				re.Set(nil, cd.NewError(cd.Unexpected, err.Error()))
			}
			return
		}
	}

	if paramPtr.ProjectType == common.Web {
		if err := s.initializeWeb(paramPtr.ProjectName, paramPtr.WorkingPath); err != nil {
			log.Errorf("initialize web failed, err:%s", err.Error())
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

	yarnInstallCmd := exec.Command("yarn", "install")
	yarnInstallCmd.Dir = projectPath
	yarnInstallCmd.Stderr = &byteBuffer
	if err = yarnInstallCmd.Run(); err != nil {
		log.Errorf("yarn failed, projectName:%s, \nstdErr:%v, \nerr:%s", projectName, byteBuffer.String(), err.Error())
		return
	}

	return
}
