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
	"github.com/muidea/magicCodeAgent/pkg/util"
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

func (s *Module) verifyWorkingPath(workingPath string, templatePtr *models.Template) (ret string, err *cd.Error) {
	if workingPath == "" {
		err = cd.NewError(cd.IllegalParam, "working path is empty")
		return
	}

	if !filepath.IsAbs(workingPath) {
		err = cd.NewError(cd.IllegalParam, "working path must be absolute path")
		return
	}

	if templatePtr.IsWeb() {
		ret = "src/pages"
		return
	}

	normalizedPath := filepath.Clean(workingPath)
	// 检查是否以 /modules/kernel 或 /modules/blocks 结尾
	if strings.HasSuffix(normalizedPath, "/modules/kernel") || strings.HasSuffix(normalizedPath, "/modules/blocks") {
		// 逐级向上查找 go.mod
		dir := normalizedPath
		for {
			goModPath := filepath.Join(dir, "go.mod")
			if _, statErr := os.Stat(goModPath); statErr == nil {
				// 找到 go.mod，计算相对路径
				relPath, relErr := filepath.Rel(dir, normalizedPath)
				if relErr != nil {
					err = cd.NewError(cd.Unexpected, "failed to compute relative path")
					return
				}
				ret = relPath
				return
			}
			// 向上移动一级
			parentDir := filepath.Dir(dir)
			if parentDir == dir {
				// 到达根目录，未找到 go.mod
				break
			}
			dir = parentDir
		}
		err = cd.NewError(cd.IllegalParam, "working path must be a go module")
		return
	}

	err = cd.NewError(cd.IllegalParam, "working path must be under {parentPath}/modules/kernel or {parentPath}/modules/blocks")
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
	if vErr := imbmcfg.Validate(); vErr != nil {
		if re != nil {
			re.Set(nil, cd.NewError(cd.IllegalParam, vErr.Error()))
		}
		return
	}

	templatePtr := ev.GetData("template").(*models.Template)

	workingPath := config.CurrentWorkingPath()
	workingPath = filepath.Join(workingPath, imbmcfg.GetPath())
	modulePath, err := s.verifyWorkingPath(workingPath, templatePtr)
	if err != nil {
		log.Errorf("create module failed, verify module working path failed, err: %v", err)
		if re != nil {
			re.Set(nil, err)
		}
		return
	}

	moduleWorkingPath := path.Join(workingPath, util.ToLower(imbmcfg.GetName()))
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
	}

	if re != nil {
		re.Set(nil, err)
	}
}

func (s *Module) refresh(ev event.Event, re event.Result) {

}

func (s *Module) destroy(ev event.Event, re event.Result) {

}
