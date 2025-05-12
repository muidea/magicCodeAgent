package biz

import (
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/log"
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

	helper := biz.NewHelper(common.ProjectTemplateBasePath)
	err = helper.Create(paramPtr.Name, "", paramPtr)
	if err != nil {
		log.Errorf("create project failed, projectName:%s, err:%s", paramPtr.Name, err.Error())
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
