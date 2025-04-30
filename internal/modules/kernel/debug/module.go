package debug

import (
	"os"
	"syscall"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/framework/plugin/module"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicCodeAgent/internal/modules/kernel/debug/biz"
	"github.com/muidea/magicCodeAgent/internal/modules/kernel/debug/pkg/common"
)

func init() {
	module.Register(New())
}

func New() *Debug {
	return &Debug{}
}

type Debug struct {
	bizPtr *biz.Debug
}

func (s *Debug) ID() string {
	return common.DebugModule
}

func (s *Debug) Weight() int {
	return 999
}

func (s *Debug) Setup(eventHub event.Hub, backgroundRoutine task.BackgroundRoutine) (err *cd.Error) {

	s.bizPtr = biz.New(eventHub, backgroundRoutine)
	return
}

func (s *Debug) Run() (err *cd.Error) {
	s.bizPtr.Run()

	// 运行完成，主动退出进程
	s.terminates()
	return
}

func (s *Debug) Teardown() {
}

func (s *Debug) terminates() {
	// 获取当前进程
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		return
	}

	// 发送SIGINT信号
	p.Signal(syscall.SIGINT)
}
