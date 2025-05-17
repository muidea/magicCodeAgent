package biz

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

type BaseCommander struct {
	debugger            sync.Mutex
	currentCommand      string
	cancelCommanderFunc context.CancelFunc

	// 调试缓冲区
	stdErrBuffer *CommandBuffer
	stdOutBuffer *CommandBuffer
}

// NewBaseCommander 创建一个新的基础调试器
func NewBaseCommander(maxLineSize int) *BaseCommander {
	return &BaseCommander{
		stdErrBuffer: NewCommandBuffer(maxLineSize),
		stdOutBuffer: NewCommandBuffer(maxLineSize),
	}
}

func (s *BaseCommander) Run(workspaceDir, command string, args []string) (err error) {
	var cmdPtr *exec.Cmd
	func() {
		s.debugger.Lock()
		defer s.debugger.Unlock()
		currentCmad := fmt.Sprintf("%s %s", command, args)
		if s.currentCommand == currentCmad {
			return
		}
		if s.cancelCommanderFunc != nil {
			s.cancelCommanderFunc()
		}

		s.stdErrBuffer.Reset()
		s.stdOutBuffer.Reset()
		s.currentCommand = currentCmad
		cancelContext, cancelFunc := context.WithCancel(context.Background())
		cmdPtr = exec.CommandContext(cancelContext, command, args...)
		s.cancelCommanderFunc = cancelFunc
		cmdPtr.Dir = workspaceDir
		cmdPtr.Stderr = s.stdErrBuffer
		cmdPtr.Stdout = s.stdOutBuffer
	}()

	err = cmdPtr.Run()
	s.debugger.Lock()
	defer s.debugger.Unlock()
	s.cancelCommanderFunc = nil
	s.currentCommand = ""
	return
}

func (s *BaseCommander) Terminate() (err error) {
	s.debugger.Lock()
	defer s.debugger.Unlock()

	if s.cancelCommanderFunc != nil {
		s.cancelCommanderFunc()
		s.stdErrBuffer.Reset()
		s.stdOutBuffer.Reset()
	}

	s.cancelCommanderFunc = nil
	s.currentCommand = ""
	return
}

func (s *BaseCommander) Status() (statusMsg, errorMsg string) {
	var statusStringBuffer strings.Builder
	s.debugger.Lock()
	defer s.debugger.Unlock()

	if s.cancelCommanderFunc != nil {
		statusStringBuffer.WriteString("running: " + s.currentCommand + "\n")
	} else {
		statusStringBuffer.WriteString("stopped: " + s.currentCommand + "\n")
	}

	if !s.stdOutBuffer.IsEmpty() {
		statusStringBuffer.WriteString("latest logs: \n")
		for _, val := range s.stdOutBuffer.GetLatestRecords(defaultMaxBufferLinesSize) {
			statusStringBuffer.WriteString(val + "\n")
		}
	}
	statusMsg = statusStringBuffer.String()

	if !s.stdErrBuffer.IsEmpty() {
		errorStringBuffer := strings.Builder{}
		errorStringBuffer.WriteString("latest errors: \n")
		for _, val := range s.stdErrBuffer.GetLatestRecords(defaultMaxBufferLinesSize) {
			errorStringBuffer.WriteString(val + "\n")
		}
		errorMsg = errorStringBuffer.String()
	}

	return
}
