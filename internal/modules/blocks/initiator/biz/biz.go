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

func (s *Initiator) verifyWorkingPath(initiatorPath string) (ret string, err *cd.Error) {
	if initiatorPath == "" {
		err = cd.NewError(cd.IllegalParam, "working path is empty")
		return
	}

	if !filepath.IsAbs(initiatorPath) {
		err = cd.NewError(cd.IllegalParam, "working path must be absolute path")
		return
	}

	normalizedPath := filepath.Clean(initiatorPath)
	if !strings.HasSuffix(normalizedPath, "/initiators") {
		err = cd.NewError(cd.IllegalParam, "initiator path must be under {parentPath}/initiators")
		return
	}
	// 逐级检查其父级目录下是否包含go.mod，如果没有找到，则认为initiator path非法，报错返回
	// 否则截取normalizedPath与parentPath的相对路径，作为initiator path
	// 例如 normalizedPath值为/home/rangh/codespace/magicCodeAgent/internal/initiators
	// go.mod在 /home/rangh/codespace/magicCodeAgent目录下
	// 逐级检查后截取结果为internal/initiators, 将其赋值给ret
	parentPath := normalizedPath[:len(normalizedPath)-len("/initiators")]
	// 逐级检查父级目录是否包含go.mod
	dir := parentPath
	for {
		goModPath := filepath.Join(dir, "go.mod")
		if _, statErr := os.Stat(goModPath); statErr == nil {
			// 找到go.mod，计算相对路径
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
			// 到达根目录，未找到go.mod
			break
		}
		dir = parentDir
	}
	err = cd.NewError(cd.IllegalParam, "no go.mod found in any parent directory")
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

	if vErr := imbicfg.Validate(); vErr != nil {
		if re != nil {
			re.Set(nil, cd.NewError(cd.IllegalParam, vErr.Error()))
		}
		return
	}

	templatePtr := ev.GetData("template").(*models.Template)

	workingPath := config.CurrentWorkingPath()
	initiatorPath := filepath.Join(workingPath, imbicfg.GetPath())
	initiatorPath, err := s.verifyWorkingPath(initiatorPath)
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
