package cli

import (
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/framework/plugin/module"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicCodeAgent/internal/modules/kernel/cli/biz"
	"github.com/muidea/magicCodeAgent/internal/modules/kernel/cli/pkg/common"
)

func init() {
	module.Register(New())
}

func New() *Cli {
	return &Cli{}
}

type Cli struct {
	bizPtr *biz.Cli
}

func (s *Cli) ID() string {
	return common.CliModule
}

func (s *Cli) Weight() int {
	return 999
}

func (s *Cli) Setup(eventHub event.Hub, backgroundRoutine task.BackgroundRoutine) (err *cd.Error) {

	s.bizPtr = biz.New(eventHub, backgroundRoutine)
	return
}

func (s *Cli) Run() (err *cd.Error) {
	s.bizPtr.Run()

	return
}

func (s *Cli) Teardown() {
}
