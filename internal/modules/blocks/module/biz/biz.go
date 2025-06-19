package biz

import (
	"os"
	"path"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/log"
	fp "github.com/muidea/magicCommon/foundation/path"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicCodeAgent/internal/config"
	"github.com/muidea/magicCodeAgent/internal/modules/base/biz"
	"github.com/muidea/magicCodeAgent/internal/modules/blocks/module/pkg/common"
	tc "github.com/muidea/magicCodeAgent/pkg/common"
)

type Module struct {
	biz.Base
}

func New(
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine,
) *Module {
	ptr := &Module{
		Base: biz.New(common.ModuleModule, eventHub, backgroundRoutine),
	}

	ptr.SubscribeFunc(common.CreateNewModule, ptr.create)
	ptr.SubscribeFunc(common.RefreshModule, ptr.refresh)
	ptr.SubscribeFunc(common.DestroyModule, ptr.destroy)

	return ptr
}

func (s *Module) create(ev event.Event, re event.Result) {
	paramPtr, paramOK := ev.Data().(*tc.CreateParam)
	if !paramOK {
		log.Errorf("create module failed, illegal create module param")
		if re != nil {
			re.Set(nil, cd.NewError(cd.Unexpected, "illegal create module param"))
		}
		return
	}
	if !fp.Exist(paramPtr.WorkingPath) {
		if err := os.MkdirAll(paramPtr.WorkingPath, 0755); err != nil {
			log.Errorf("create module failed, create module working path failed, err: %v", err)
			if re != nil {
				re.Set(nil, cd.NewError(cd.Unexpected, "create module working path failed"))
			}
			return
		}
	}

	prefixPath, err := config.VerifyModuleWorkPath(paramPtr.ProjectType, paramPtr.WorkingPath)
	if err != nil {
		log.Errorf("verify module work path failed, projectType:%s, workingPath:%s, err:%s", paramPtr.ProjectType, paramPtr.WorkingPath, err.Error())

		if re != nil {
			re.Set(nil, err)
		}
		return
	}

	//paramPtr.ProjectName = paramPtr.ProjectName
	helper := biz.NewHelper(path.Join(common.ModuleTemplateBasePath, paramPtr.ProjectType))
	err = helper.Create(prefixPath, paramPtr)
	if err != nil {
		log.Errorf("create module failed, projectName:%s, prefixPath:%s, err:%s", paramPtr.ProjectName, prefixPath, err.Error())
		if re != nil {
			re.Set(nil, err)
		}
		return
	}
}

func (s *Module) refresh(ev event.Event, re event.Result) {

}

func (s *Module) destroy(ev event.Event, re event.Result) {

}
