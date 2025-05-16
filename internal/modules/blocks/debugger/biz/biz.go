package biz

import (
	"sync"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicCodeAgent/internal/modules/base/biz"
	imbdc "github.com/muidea/magicCodeAgent/internal/modules/blocks/debugger/pkg/common"
)

type Debugger struct {
	biz.Base

	workspace2Debugger sync.Map
}

func New(
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine,
) *Debugger {
	ptr := &Debugger{
		Base: biz.New(imbdc.DebuggerModule, eventHub, backgroundRoutine),
	}

	ptr.SubscribeFunc(imbdc.DebugCommandRun, ptr.runDebug)
	ptr.SubscribeFunc(imbdc.DebugCommandTerminate, ptr.terminateDebug)
	ptr.SubscribeFunc(imbdc.DebugCommandStatus, ptr.statusDebug)

	return ptr
}

func (s *Debugger) Run() (err *cd.Error) {
	return
}

func (s *Debugger) runDebug(ev event.Event, re event.Result) {
	workspace := ev.Data().(string)
	commandPtr := ev.GetData("command").(*imbdc.DebugCommand)
	result, err := s.RunDebug(workspace, *commandPtr)
	if re != nil {
		re.Set(result, err)
	}
}

func (s *Debugger) RunDebug(workspace string, commandPtr imbdc.DebugCommand) (ret *imbdc.DebugResult, err *cd.Error) {
	var debuggerPtr *BaseDebugger
	debuggerVal, debuggerOK := s.workspace2Debugger.Load(workspace)
	if !debuggerOK {
		debuggerPtr = NewBaseDebugger(defaultMaxBufferLinesSize)
	} else {
		debuggerPtr = debuggerVal.(*BaseDebugger)
	}

	runErr := debuggerPtr.Run(commandPtr.WorkspaceDir, commandPtr.Command, commandPtr.Args)
	if runErr != nil {
		err = cd.NewError(cd.Unexpected, runErr.Error())
		return
	}

	ret = &imbdc.DebugResult{
		Workspace: workspace,
		Status:    debuggerPtr.Status(),
	}

	return
}

func (s *Debugger) terminateDebug(ev event.Event, re event.Result) {
	workspace := ev.Data().(string)
	result, err := s.TerminateDebug(workspace)
	if re != nil {
		re.Set(result, err)
	}
}

func (s *Debugger) TerminateDebug(workspace string) (ret *imbdc.DebugResult, err *cd.Error) {
	var terminateErr error
	debuggerVal, debuggerOK := s.workspace2Debugger.Load(workspace)
	if debuggerOK {
		debuggerPtr := debuggerVal.(*BaseDebugger)
		terminateErr = debuggerPtr.Terminate()
	}
	if terminateErr != nil {
		err = cd.NewError(cd.Unexpected, terminateErr.Error())
		return
	}
	ret = &imbdc.DebugResult{
		Workspace: workspace,
		Status:    "terminated ok!",
	}
	return
}

func (s *Debugger) statusDebug(ev event.Event, re event.Result) {
	workspace := ev.Data().(string)
	result, err := s.StatusDebug(workspace)
	if re != nil {
		re.Set(result, err)
	}
}

func (s *Debugger) StatusDebug(workspace string) (ret *imbdc.DebugResult, err *cd.Error) {
	var statusValue string = "workspace not found"
	debuggerVal, debuggerOK := s.workspace2Debugger.Load(workspace)
	if debuggerOK {
		debuggerPtr := debuggerVal.(*BaseDebugger)
		statusValue = debuggerPtr.Status()
	}
	ret = &imbdc.DebugResult{
		Workspace: workspace,
		Status:    statusValue,
	}
	return
}
