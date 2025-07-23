package biz

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/log"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicCodeAgent/internal/config"
	"github.com/muidea/magicCodeAgent/internal/modules/base/biz"
	imbicfg "github.com/muidea/magicCodeAgent/internal/modules/blocks/initiator/config"
	imbic "github.com/muidea/magicCodeAgent/internal/modules/blocks/initiator/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
	"github.com/muidea/magicCodeAgent/pkg/util"
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

func (s *Initiator) verifyWorkingPath(workingPath string) (ret string, err *cd.Error) {
	if workingPath == "" {
		err = cd.NewError(cd.IllegalParam, "working path is empty")
		return
	}

	if !filepath.IsAbs(workingPath) {
		err = cd.NewError(cd.IllegalParam, "working path must be absolute path")
		return
	}

	if vErr := imbicfg.Validate(); vErr != nil {
		err = cd.NewError(cd.IllegalParam, vErr.Error())
		return
	}

	normalizedPath := filepath.Clean(workingPath)
	if !strings.HasSuffix(normalizedPath, "/internal/initiators") {
		err = cd.NewError(cd.IllegalParam, "working path must be internal/initiators")
		return
	}
	codeRootPath := normalizedPath[:len(normalizedPath)-len("/internal/initiators")]
	goModPath := filepath.Join(codeRootPath, "go.mod")
	if _, oErr := os.Stat(goModPath); oErr != nil {
		err = cd.NewError(cd.IllegalParam, "working path must be a go module")
		return
	}

	ret = "internal/initiators"
	return
}

func (s *Initiator) create(ev event.Event, re event.Result) {
	workspacePtr, workspaceOK := ev.Data().(*models.Workspace)
	if !workspaceOK {
		log.Errorf("create initiator failed, illegal create initiator param")
		if re != nil {
			re.Set(nil, cd.NewError(cd.Unexpected, "illegal create initiator param"))
		}
		return
	}
	templatePtr := ev.GetData("template").(*models.Template)

	workingPath := config.CurrentWorkingPath()
	workingPath = filepath.Join(workingPath, imbicfg.GetPath())
	initiatorPath, err := s.verifyWorkingPath(workingPath)
	if err != nil {
		log.Errorf("create initiator failed, verify initiator working path failed, error: %v", err)
		if re != nil {
			re.Set(nil, err)
		}
		return
	}

	initiatorWorkingPath := path.Join(workingPath, util.ToLower(imbicfg.GetName()))
	if oErr := os.MkdirAll(initiatorWorkingPath, 0755); oErr != nil {
		log.Errorf("create initiator failed, create initiator working path failed, error: %v", oErr)
		if re != nil {
			re.Set(nil, cd.NewError(cd.Unexpected, "create initiator working path failed"))
		}
		return
	}
	log.Infof("create initiator working path: %v", initiatorWorkingPath)
	initiatorTemplateBasePath := path.Join(models.TemplatePath, templatePtr.Category, templatePtr.Name, "template", "initiator")
	createrHelper := biz.NewHelper(templatePtr.HostRootFS, initiatorTemplateBasePath)
	err = createrHelper.Create(&common.CreateParam{
		WorkingPath:   initiatorWorkingPath,
		ProjectName:   workspacePtr.Name,
		InitiatorName: imbicfg.GetName(),
		InitiatorPath: initiatorPath,
	})
	if err != nil {
		log.Errorf("create initiator failed, error: %v", err)
	}
	if re != nil {
		re.Set(nil, err)
	}
}

func (s *Initiator) refresh(ev event.Event, re event.Result) {

}

func (s *Initiator) destroy(ev event.Event, re event.Result) {

}
