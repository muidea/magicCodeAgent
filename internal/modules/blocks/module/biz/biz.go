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

func (s *Module) verifyWorkingPath(workingPath string) (ret string, err *cd.Error) {
	if workingPath == "" {
		err = cd.NewError(cd.IllegalParam, "working path is empty")
		return
	}

	if !filepath.IsAbs(workingPath) {
		err = cd.NewError(cd.IllegalParam, "working path must be absolute path")
		return
	}

	if vErr := imbmcfg.Validate(); vErr != nil {
		err = cd.NewError(cd.IllegalParam, vErr.Error())
		return
	}

	normalizedPath := filepath.Clean(workingPath)
	if strings.HasSuffix(normalizedPath, "/internal/modules/kernel") {
		codeRootPath := normalizedPath[:len(normalizedPath)-len("/internal/modules/kernel")]
		goModPath := filepath.Join(codeRootPath, "go.mod")
		if _, oErr := os.Stat(goModPath); oErr != nil {
			err = cd.NewError(cd.IllegalParam, "working path must be a go module")
			return
		}
		ret = "internal/modules/kernel"
		return
	}

	if strings.HasSuffix(normalizedPath, "/internal/modules/blocks") {
		codeRootPath := normalizedPath[:len(normalizedPath)-len("/internal/modules/blocks")]
		goModPath := filepath.Join(codeRootPath, "go.mod")
		if _, oErr := os.Stat(goModPath); oErr != nil {
			err = cd.NewError(cd.IllegalParam, "working path must be a go module")
			return
		}
		ret = "internal/modules/blocks"
		return
	}

	err = cd.NewError(cd.IllegalParam, "working path must be in /internal/modules/kernel or /internal/modules/blocks")
	return
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
	modulePath, err := s.verifyWorkingPath(workingPath)
	if err != nil {
		log.Errorf("create module failed, verify module working path failed, err: %v", err)
		if re != nil {
			re.Set(nil, err)
		}
		return
	}

	moduleWorkingPath := workingPath //filepath.Join(workingPath, imbmcfg.GetName())
	if oErr := os.MkdirAll(moduleWorkingPath, 0755); oErr != nil {
		log.Errorf("create module failed, create module working path failed, err: %v", oErr)
		if re != nil {
			re.Set(nil, cd.NewError(cd.Unexpected, "create module working path failed"))
		}
		return
	}
	log.Infof("create module working path: %v", moduleWorkingPath)
	moduleTemplateBasePath := path.Join(models.TemplatePath, templatePtr.Category, templatePtr.Name, "template", "module")
	createrHelper := biz.NewHelper(templatePtr.HostRootFS, moduleTemplateBasePath)
	err = createrHelper.Create(&common.CreateParam{
		WorkingPath: moduleWorkingPath,
		ProjectName: workspacePtr.Name,
		ModuleName:  imbmcfg.GetName(),
		ModulePath:  modulePath,
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
