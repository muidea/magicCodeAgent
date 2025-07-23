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
		log.Infof("workspace: %s", workspaceDir)
		log.Infof("run command: '%s'", s.currentCommand)

		s.stdErrBuffer.Reset()
		s.stdOutBuffer.Reset()

		cancelContext, cancelFunc := context.WithCancel(context.Background())
		s.cancelCommanderFunc = cancelFunc
		cmdContext = cancelContext

		cmdPtr = exec.CommandContext(cmdContext, command, args...)
		sysEnv := append(os.Environ(), "workspace="+s.workspace, "LOG_LEVEL=4")
		cmdPtr.Env = sysEnv
		cmdPtr.Dir = workspaceDir
		cmdPtr.Stderr = s.stdErrBuffer
		cmdPtr.Stdout = s.stdOutBuffer
		cmdPtr.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	}()

	// 后台监听上下文取消，确保杀死整组进程
	go func(process *os.Process, ctx context.Context) {
		if ctx == nil || process == nil {
			return
		}
		<-ctx.Done()

		// 使用负 PID 杀掉整个 session（setsid 创建的）
		log.Warnf("cancel context received, killing process group for PID %d", process.Pid)
		_ = syscall.Kill(-process.Pid, syscall.SIGKILL)
	}(cmdPtr.Process, cmdContext)

	// 阻塞直到命令完成
	err = cmdPtr.Run()
	if err != nil {
		log.Errorf("run command: '%s' failed, err: %s", s.currentCommand, err.Error())
	}

	s.debugger.Lock()
	defer s.debugger.Unlock()
	s.cancelCommanderFunc = nil
	return
}

func (s *BaseCommander) Terminate() (err error) {
	s.debugger.Lock()
	defer s.debugger.Unlock()

	log.Warnf("terminate command: '%s'", s.currentCommand)
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
		statusStringBuffer.WriteString("'" + s.currentCommand + "' is running")
	} else {
		statusStringBuffer.WriteString("'" + s.currentCommand + "' is done")
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
