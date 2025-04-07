package cli

import (
	"os"
	"syscall"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/framework/plugin/module"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicToolkit/internal/modules/kernel/cli/biz"
	"github.com/muidea/magicToolkit/internal/modules/kernel/cli/pkg/common"
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

	// 运行完成，主动退出进程
	s.terminates()
	return
}

func (s *Cli) Teardown() {
}

func (s *Cli) terminates() {
	// 获取当前进程
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		return
	}

	// 发送SIGINT信号
	p.Signal(syscall.SIGINT)
}
