package biz

import (
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/log"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicToolkit/internal/modules/base/biz"
	"github.com/muidea/magicToolkit/internal/modules/module/pkg/common"
	tc "github.com/muidea/magicToolkit/pkg/common"
	"github.com/muidea/magicToolkit/pkg/util"
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

	projectName, prefixPath, err := util.VerifyModuleWorkPath(paramPtr.WorkingPath)
	if err != nil {
		log.Errorf("verify module work path failed, workingPath:%s, err:%s", paramPtr.WorkingPath, err.Error())

		if re != nil {
			re.Set(nil, err)
		}
		return
	}

	const moduleTemplateBasePath = "template/module"
	helper := biz.NewHelper(moduleTemplateBasePath)
	err = helper.Create(projectName, prefixPath, paramPtr)
	if err != nil {
		log.Errorf("create module failed, projectName:%s, prefixPath:%s, err:%s", projectName, prefixPath, err.Error())
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
