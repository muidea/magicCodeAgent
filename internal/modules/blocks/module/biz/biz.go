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
	imbmcfg "github.com/muidea/magicCodeAgent/internal/modules/blocks/module/config"
	imbmc "github.com/muidea/magicCodeAgent/internal/modules/blocks/module/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
)

type Module struct {
	biz.Base
}

func New(
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine,
) *Module {
	ptr := &Module{
		Base: biz.New(imbmc.ModuleModule, eventHub, backgroundRoutine),
	}

	ptr.SubscribeFunc(imbmc.CreateNewModule, ptr.create)
	ptr.SubscribeFunc(imbmc.RefreshModule, ptr.refresh)
	ptr.SubscribeFunc(imbmc.DestroyModule, ptr.destroy)

	return ptr
}

func (s *Module) verifyWorkingPath(workingPath string) *cd.Error {
	if workingPath == "" {
		return cd.NewError(cd.IllegalParam, "working path is empty")
	}

	if !filepath.IsAbs(workingPath) {
		return cd.NewError(cd.IllegalParam, "working path must be absolute path")
	}

	normalizedPath := filepath.Clean(workingPath)
	if strings.HasSuffix(normalizedPath, "/internal/modules/kernel") {
		codeRootPath := normalizedPath[:len(normalizedPath)-len("/internal/modules/kernel")]
		goModPath := filepath.Join(codeRootPath, "go.mod")
		if _, err := os.Stat(goModPath); err != nil {
			return cd.NewError(cd.IllegalParam, "working path must be a go module")
		}
		return nil
	}

	if strings.HasSuffix(normalizedPath, "/internal/modules/blocks") {
		codeRootPath := normalizedPath[:len(normalizedPath)-len("/internal/modules/blocks")]
		goModPath := filepath.Join(codeRootPath, "go.mod")
		if _, err := os.Stat(goModPath); err != nil {
			return cd.NewError(cd.IllegalParam, "working path must be a go module")
		}
		return nil
	}

	return cd.NewError(cd.IllegalParam, "working path must be /internal/modules/kernel or /internal/modules/blocks")
}

func (s *Module) create(ev event.Event, re event.Result) {
	workspacePtr, workspaceOK := ev.Data().(*models.Workspace)
	if !workspaceOK {
		log.Errorf("create module failed, illegal create module param")
		if re != nil {
			re.Set(nil, cd.NewError(cd.Unexpected, "illegal create module param"))
		}
		return
	}
	templatePtr := ev.GetData("template").(*models.Template)

	workingPath := config.CurrentWorkingPath()
	workingPath = filepath.Join(workingPath, imbmcfg.GetPath())
	err := s.verifyWorkingPath(workingPath)
	if err != nil {
		log.Errorf("create module failed, verify module working path failed, err: %v", err)
		if re != nil {
			re.Set(nil, err)
		}
		return
	}

	createrHelper := biz.NewHelper(templatePtr.HostRootFS, imbmc.ModuleTemplateBasePath)
	err = createrHelper.Create(&common.CreateParam{
		WorkingPath: workingPath,
		ProjectName: workspacePtr.Name,
		ModuleName:  imbmcfg.GetName(),
		ModulePath:  imbmcfg.GetPath(),
	})
	if err != nil {
		log.Errorf("create module failed, err: %v", err)
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
