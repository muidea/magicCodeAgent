package biz

import (
	"os"
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

func (s *Initiator) verifyWorkingPath(workingPath string) *cd.Error {
	if workingPath == "" {
		return cd.NewError(cd.IllegalParam, "working path is empty")
	}

	if !filepath.IsAbs(workingPath) {
		return cd.NewError(cd.IllegalParam, "working path must be absolute path")
	}

	normalizedPath := filepath.Clean(workingPath)
	if !strings.HasSuffix(normalizedPath, "/internal/initiator") {
		return cd.NewError(cd.IllegalParam, "working path must be internal/initiator")
	}
	codeRootPath := normalizedPath[:len(normalizedPath)-len("/internal/initiator")]
	goModPath := filepath.Join(codeRootPath, "go.mod")
	if _, err := os.Stat(goModPath); err != nil {
		return cd.NewError(cd.IllegalParam, "working path must be a go module")
	}
	return nil
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
	err := s.verifyWorkingPath(workingPath)
	if err != nil {
		log.Errorf("create initiator failed, verify initiator working path failed, error: %v", err)
		if re != nil {
			re.Set(nil, err)
		}
		return
	}

	createrHelper := biz.NewHelper(templatePtr.HostRootFS, imbic.InitiatorTemplateBasePath)
	err = createrHelper.Create(&common.CreateParam{
		WorkingPath:   workingPath,
		ProjectName:   workspacePtr.Name,
		InitiatorName: imbicfg.GetName(),
		InitiatorPath: imbicfg.GetPath(),
	})
	if err != nil {
		log.Errorf("create initiator failed, error: %v", err)
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
