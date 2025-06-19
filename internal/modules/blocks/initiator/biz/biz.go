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
	imbic "github.com/muidea/magicCodeAgent/internal/modules/blocks/initiator/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/common"
)

type Initiator struct {
	biz.Base
}

func New(
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine,
) *Initiator {
	ptr := &Initiator{
		Base: biz.New(imbic.InitiatorModule, eventHub, backgroundRoutine),
	}

	ptr.SubscribeFunc(imbic.CreateNewInitiator, ptr.create)
	ptr.SubscribeFunc(imbic.RefreshInitiator, ptr.refresh)
	ptr.SubscribeFunc(imbic.DestroyInitiator, ptr.destroy)

	return ptr
}

func (s *Initiator) create(ev event.Event, re event.Result) {
	paramPtr, paramOK := ev.Data().(*common.CreateParam)
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

	prefixPath, err := config.VerifyInitiatorWorkPath(paramPtr.WorkingPath)
	if err != nil {
		log.Errorf("create initiator failed, verify initiator working path failed, err: %v", err)
		if re != nil {
			re.Set(nil, err)
		}
		return
	}

	//paramPtr.ProjectName = projectName
	helper := biz.NewHelper(path.Join(imbic.InitiatorTemplateBasePath, common.Service))
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
