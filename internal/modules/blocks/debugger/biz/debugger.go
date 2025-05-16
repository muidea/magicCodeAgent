package biz

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

type BaseDebugger struct {
	debugger           sync.Mutex
	currentCommand     string
	cancelDebuggerFunc context.CancelFunc

	// 调试缓冲区
	buffer *DebugBuffer
}

// NewBaseDebugger 创建一个新的基础调试器
func NewBaseDebugger(maxLineSize int) *BaseDebugger {
	return &BaseDebugger{
		buffer: NewDebugBuffer(maxLineSize),
	}
}

func (s *BaseDebugger) Run(workspaceDir, command string, args []string) (err error) {
	var cmdPtr *exec.Cmd
	func() {
		s.debugger.Lock()
		defer s.debugger.Unlock()
		currentCmad := fmt.Sprintf("%s %s", command, args)
		if s.currentCommand == currentCmad {
			return
		}
		if s.cancelDebuggerFunc != nil {
			s.cancelDebuggerFunc()
			s.buffer.Reset()
		}
		s.currentCommand = currentCmad
		cancelContext, cancelFunc := context.WithCancel(context.Background())
		cmdPtr = exec.CommandContext(cancelContext, command, args...)
		s.cancelDebuggerFunc = cancelFunc
		cmdPtr.Dir = workspaceDir
		cmdPtr.Stderr = s
	}()

	err = cmdPtr.Run()
	s.debugger.Lock()
	defer s.debugger.Unlock()
	s.cancelDebuggerFunc = nil
	s.currentCommand = ""
	return
}

func (s *BaseDebugger) Terminate() (err error) {
	s.debugger.Lock()
	defer s.debugger.Unlock()

	if s.cancelDebuggerFunc != nil {
		s.cancelDebuggerFunc()
		s.buffer.Reset()
	}

	s.cancelDebuggerFunc = nil
	s.currentCommand = ""
	return
}

func (s *BaseDebugger) Status() (status string) {
	var stringBuffer strings.Builder
	s.debugger.Lock()
	defer s.debugger.Unlock()

	if s.cancelDebuggerFunc != nil {
		stringBuffer.WriteString("running: " + s.currentCommand + "\n")
	} else {
		stringBuffer.WriteString("stopped: " + s.currentCommand + "\n")
	}

	stringBuffer.WriteString("latest logs: \n")
	for _, val := range s.buffer.GetLatestRecords(defaultMaxBufferLinesSize) {
		stringBuffer.WriteString(val + "\n")
	}

	status = stringBuffer.String()
	return
}

func (s *BaseDebugger) Write(p []byte) (n int, err error) {
	// 如果缓冲区未初始化，则创建一个默认的缓冲区
	if s.buffer == nil {
		s.buffer = NewDebugBuffer(100) // 默认缓冲100行记录
	}

	// 将数据写入缓冲区
	return s.buffer.Write(p)
}

// SetBufferSize 设置调试器缓冲区大小
func (s *BaseDebugger) SetBufferSize(size int) {
	// 如果缓冲区未初始化，则创建一个新的缓冲区
	if s.buffer == nil {
		s.buffer = NewDebugBuffer(size)
		return
	}

	// 设置现有缓冲区的大小
	s.buffer.SetBufferSize(size)
}

// GetFunctionRecords 查询包含指定函数名的记录
func (s *BaseDebugger) GetFunctionRecords(functionName string) []string {
	if s.buffer == nil {
		return []string{}
	}

	return s.buffer.GetFunctionRecords(functionName)
}

// GetAllRecords 获取所有缓冲的记录
func (s *BaseDebugger) GetAllRecords() []string {
	if s.buffer == nil {
		return []string{}
	}

	return s.buffer.GetAllRecords()
}

// GetLatestRecords 获取最新的N行记录
func (s *BaseDebugger) GetLatestRecords(count int) []string {
	if s.buffer == nil {
		return []string{}
	}

	return s.buffer.GetLatestRecords(count)
}
