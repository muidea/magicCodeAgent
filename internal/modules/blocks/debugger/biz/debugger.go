package biz

import (
	"context"
	"os/exec"
)

type BaseDebugger struct {
	cancelDebuggerFunc context.CancelFunc
}

func (s *BaseDebugger) Run(workspaceDir, command string, args []string) (err error) {
	cancelContext, cancelFunc := context.WithCancel(context.Background())
	cmdPtr := exec.CommandContext(cancelContext, command, args...)
	s.cancelDebuggerFunc = cancelFunc
	cmdPtr.Dir = workspaceDir
	cmdPtr.Stderr = s

	err = cmdPtr.Run()
	return
}

func (s *BaseDebugger) Terminate() (err error) {
	if s.cancelDebuggerFunc != nil {
		s.cancelDebuggerFunc()
	}
	s.cancelDebuggerFunc = nil
	return
}

func (s *BaseDebugger) Write(p []byte) (n int, err error) {
	return
}
