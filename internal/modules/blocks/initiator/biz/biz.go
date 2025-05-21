package biz

import (
	"os"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/log"
	fp "github.com/muidea/magicCommon/foundation/path"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicCodeAgent/internal/config"
	"github.com/muidea/magicCodeAgent/internal/modules/base/biz"
	"github.com/muidea/magicCodeAgent/internal/modules/blocks/initiator/pkg/common"
	tc "github.com/muidea/magicCodeAgent/pkg/common"
)

type Initiator struct {
	biz.Base
}

func New(
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine,
) *Initiator {
	ptr := &Initiator{
		Base: biz.New(common.InitiatorModule, eventHub, backgroundRoutine),
	}

	ptr.SubscribeFunc(common.CreateNewInitiator, ptr.create)
	ptr.SubscribeFunc(common.RefreshInitiator, ptr.refresh)
	ptr.SubscribeFunc(common.DestroyInitiator, ptr.destroy)

	return ptr
}

func (s *Initiator) create(ev event.Event, re event.Result) {
	paramPtr, paramOK := ev.Data().(*tc.CreateParam)
	if !paramOK {
		log.Errorf("create initiator failed, illegal create initiator param")
		if re != nil {
			re.Set(nil, cd.NewError(cd.Unexpected, "illegal create initiator param"))
		}
		return
	}
	if !fp.Exist(paramPtr.WorkingPath) {
		if err := os.MkdirAll(paramPtr.WorkingPath, 0755); err != nil {
			log.Errorf("create initiator failed, create initiator working path failed, err: %v", err)
			if re != nil {
				re.Set(nil, cd.NewError(cd.Unexpected, "create initiator working path failed"))
			}
			return
		}
	}

	projectName, prefixPath, err := config.VerifyInitiatorWorkPath(paramPtr.WorkingPath)
	if err != nil {
		log.Errorf("create initiator failed, verify initiator working path failed, err: %v", err)
		if re != nil {
			re.Set(nil, err)
		}
		return
	}

	paramPtr.ProjectName = projectName
	helper := biz.NewHelper(common.InitiatorTemplateBasePath)
	err = helper.Create(prefixPath, paramPtr)
	if err != nil {
		log.Errorf("create initiator failed, create initiator failed, err: %v", err)
		if re != nil {
			re.Set(nil, err)
		}
		return
	}
}

func (s *Initiator) refresh(ev event.Event, re event.Result) {

}

func (s *Initiator) destroy(ev event.Event, re event.Result) {

}
