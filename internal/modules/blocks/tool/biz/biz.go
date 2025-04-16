package biz

import (
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicCodeAgent/internal/modules/base/biz"
	"github.com/muidea/magicCodeAgent/internal/modules/blocks/tool/pkg/common"
)

const (
	WorkspaceDir = "/home/rangh/webspace"
)

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
