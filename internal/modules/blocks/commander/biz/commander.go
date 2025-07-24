package biz

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"

	"github.com/muidea/magicCommon/foundation/log"
)

type BaseCommander struct {
	workspace           string
	debugger            sync.Mutex
	currentCommand      string
	cancelCommanderFunc context.CancelFunc

	// 调试缓冲区
	stdErrBuffer *LogBuffer
	stdOutBuffer *LogBuffer
}

// NewBaseCommander 创建一个新的基础调试器
func NewBaseCommander(workspace string, maxLineSize int) *BaseCommander {
	return &BaseCommander{
		workspace:    workspace,
		stdErrBuffer: NewLogBuffer(maxLineSize),
		stdOutBuffer: NewLogBuffer(maxLineSize),
	}
}

func (s *BaseCommander) Run(workspaceDir, command string, args []string) (err error) {
	var cmdPtr *exec.Cmd
	var cmdContext context.Context
	func() {
		s.debugger.Lock()
		defer s.debugger.Unlock()

		if s.currentCommand == command {
			return
		}

		// 如果已有正在运行的命令，先取消
		if s.cancelCommanderFunc != nil {
			s.cancelCommanderFunc()
		}

		s.currentCommand = command
		log.Infof("working path: %s", workspaceDir)
		log.Infof("run command: '%s'", s.currentCommand)

		s.stdErrBuffer.Reset()
		s.stdOutBuffer.Reset()

		cancelContext, cancelFunc := context.WithCancel(context.Background())
		s.cancelCommanderFunc = cancelFunc
		cmdContext = cancelContext

		cmdPtr = exec.CommandContext(cmdContext, command, args...)
		var sysEnv []string
		for _, env := range os.Environ() {
			if strings.HasPrefix(env, "workspace=") {
				continue
			}
			if strings.HasPrefix(env, "BindPort=") {
				continue
			}
			sysEnv = append(sysEnv, env)
		}
		sysEnv = append(sysEnv, "workspace="+s.workspace, "LOG_LEVEL=4")
		cmdPtr.Env = sysEnv
		cmdPtr.Dir = workspaceDir
		cmdPtr.Stderr = s.stdErrBuffer
		cmdPtr.Stdout = s.stdOutBuffer
		cmdPtr.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	}()

	// 阻塞直到命令完成
	err = cmdPtr.Run()
	if err != nil {
		log.Errorf("run command: '%s' failed, err: %s", s.currentCommand, err.Error())
	} else {
		log.Infof("run command: '%s' done", s.currentCommand)
	}

	s.debugger.Lock()
	defer s.debugger.Unlock()
	s.cancelCommanderFunc = nil
	return
}

func (s *BaseCommander) Terminate() (err error) {
	s.debugger.Lock()
	defer s.debugger.Unlock()

	if s.cancelCommanderFunc != nil {
		s.cancelCommanderFunc()
	}

	s.cancelCommanderFunc = nil
	return
}

func (s *BaseCommander) Status() (statusMsg, errorMsg string) {
	var statusStringBuffer strings.Builder
	s.debugger.Lock()
	defer s.debugger.Unlock()

	if s.cancelCommanderFunc != nil {
		statusStringBuffer.WriteString("'" + s.currentCommand + "' is running")
	} else {
		statusStringBuffer.WriteString("'" + s.currentCommand + "' is done")
	}

	if !s.stdOutBuffer.IsEmpty() {
		statusStringBuffer.WriteString("\nlatest logs:\n")
		for _, val := range s.stdOutBuffer.GetLatestRecords(defaultMaxBufferLinesSize) {
			statusStringBuffer.WriteString("    " + val + "\n")
		}
	}
	statusMsg = statusStringBuffer.String()

	if !s.stdErrBuffer.IsEmpty() {
		errorStringBuffer := strings.Builder{}
		errorStringBuffer.WriteString("\nlatest errors:\n")
		for _, val := range s.stdErrBuffer.GetLatestRecords(defaultMaxBufferLinesSize) {
			errorStringBuffer.WriteString("    " + val + "\n")
		}
		errorMsg = errorStringBuffer.String()
	}

	return
}
