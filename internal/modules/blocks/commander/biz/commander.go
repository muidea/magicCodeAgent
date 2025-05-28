package biz

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/muidea/magicCommon/foundation/log"
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
	currentCmd := fmt.Sprintf("%s %s", command, args)
	func() {
		s.debugger.Lock()
		defer s.debugger.Unlock()
		if s.currentCommand == currentCmd {
			return
		}
		if s.cancelCommanderFunc != nil {
			s.cancelCommanderFunc()
		}

		log.Infof("run command: %s", currentCmd)
		s.stdErrBuffer.Reset()
		s.stdOutBuffer.Reset()
		s.currentCommand = currentCmd
		cancelContext, cancelFunc := context.WithCancel(context.Background())
		// 通过/bin/bash -c 命令执行command args
		paramStr := command + " " + strings.Join(args, " ")
		cmdPtr = exec.CommandContext(cancelContext, "/bin/bash", "-c", paramStr)
		s.cancelCommanderFunc = cancelFunc
		cmdPtr.Env = os.Environ()
		cmdPtr.Dir = workspaceDir
		cmdPtr.Stderr = s.stdErrBuffer
		cmdPtr.Stdout = s.stdOutBuffer
	}()

	err = cmdPtr.Run()
	s.debugger.Lock()
	defer s.debugger.Unlock()
	s.cancelCommanderFunc = nil
	//s.currentCommand = ""
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
		statusStringBuffer.WriteString("'" + s.currentCommand + "' is  running")
	} else {
		statusStringBuffer.WriteString("'" + s.currentCommand + "' is  done")
	}

	if !s.stdOutBuffer.IsEmpty() {
		statusStringBuffer.WriteString("\nlatest logs:")
		for _, val := range s.stdOutBuffer.GetLatestRecords(defaultMaxBufferLinesSize) {
			statusStringBuffer.WriteString("  " + val)
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
