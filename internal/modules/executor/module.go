package executor

import (
	"os"
	"syscall"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/framework/plugin/module"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicToolkit/internal/modules/executor/biz"
	"github.com/muidea/magicToolkit/internal/modules/executor/pkg/common"
)

func init() {
	module.Register(New())
}

func New() *Executor {
	return &Executor{}
}

type Executor struct {
	bizPtr *biz.Executor
}

func (s *Executor) ID() string {
	return common.ExecutorModule
}

func (s *Executor) Weight() int {
	return 999
}

func (s *Executor) Setup(eventHub event.Hub, backgroundRoutine task.BackgroundRoutine) (err *cd.Error) {

	s.bizPtr = biz.New(eventHub, backgroundRoutine)
	return
}

func (s *Executor) Run() (err *cd.Error) {
	s.bizPtr.Run()

	// 运行完成，主动退出进程
	s.terminates()
	return
}

func (s *Executor) Teardown() {
}

func (s *Executor) terminates() {
	// 获取当前进程
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		return
	}

	// 发送SIGINT信号
	p.Signal(syscall.SIGINT)
}
