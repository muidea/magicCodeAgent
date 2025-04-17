package biz

import (
	"os"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicCodeAgent/internal/modules/base/biz"
	"github.com/muidea/magicCodeAgent/internal/modules/blocks/tool/pkg/common"
)

var WorkspaceDir = "/home/rangh/codespace/magicCodeAgent/internal/modules/blocks/tool/biz/temp"

func init() {
	WorkspaceDir, _ = os.Getwd()
}

type Tool struct {
	biz.Base
}

func New(
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine,
) *Tool {
	ptr := &Tool{
		Base: biz.New(common.ToolModule, eventHub, backgroundRoutine),
	}

	return ptr
}

func (s *Tool) Run() (err *cd.Error) {
	return
}
