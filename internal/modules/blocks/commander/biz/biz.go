package biz

import (
	"fmt"
	"strings"
	"sync"
	"time"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/log"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicCodeAgent/internal/modules/base/biz"
	imbcc "github.com/muidea/magicCodeAgent/internal/modules/blocks/commander/pkg/common"
)

var defaultCommandTimeOut = 2
var defaultConsoleLinesSize = 100

type Commander struct {
	biz.Base

	workspace2Commander sync.Map

	// TODO 这里以后需要考虑多个workspace的情况
	consoleBuffer *LogBuffer
}

func New(
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine,
) *Commander {
	ptr := &Commander{
		Base:          biz.New(imbcc.CommanderModule, eventHub, backgroundRoutine),
		consoleBuffer: NewLogBuffer(defaultConsoleLinesSize),
	}

	ptr.SubscribeFunc(imbcc.CommandRun, ptr.runCommand)
	ptr.SubscribeFunc(imbcc.CommandTerminate, ptr.terminateCommand)
	ptr.SubscribeFunc(imbcc.CommandStatus, ptr.statusCommand)

	return ptr
}

func (s *Commander) Run() (err *cd.Error) {
	return
}

func (s *Commander) runCommand(ev event.Event, re event.Result) {
	workspace := ev.Data().(string)
	commandPtr := ev.GetData("command").(*imbcc.CommanderParam)
	result, err := s.RunCommand(workspace, commandPtr)
	if re != nil {
		re.Set(result, err)
	}
}

func (s *Commander) RunCommand(workspace string, paramPtr *imbcc.CommanderParam) (ret *imbcc.CommanderResult, err *cd.Error) {
	var commanderPtr *BaseCommander
	commanderVal, commanderOK := s.workspace2Commander.Load(workspace)
	if !commanderOK {
		commanderPtr = NewBaseCommander(workspace, defaultMaxBufferLinesSize)
	} else {
		commanderPtr = commanderVal.(*BaseCommander)
	}

	s.consoleBuffer.Reset()

	preTime := time.Now()
	timeoutVal := time.Minute * time.Duration(defaultCommandTimeOut)
	if paramPtr.Timeout > 0 {
		timeoutVal = time.Duration(paramPtr.Timeout) * time.Second
	}
	s.Base.SyncTaskWithTimeOut(func() {
		_ = commanderPtr.Run(paramPtr.WorkspaceDir, paramPtr.Command, paramPtr.Args)
	}, timeoutVal)

	commanderPtr.Terminate()

	log.Warnf("run command elapsed: %v", time.Since(preTime))
	ret = &imbcc.CommanderResult{
		Workspace: workspace,
	}
	ret.StatusMessage, ret.ErrorMessage = commanderPtr.Status()
	if s.consoleBuffer != nil {
		records := s.consoleBuffer.GetAllRecords()
		if len(records) > 0 {
			ret.ConsoleMessage = strings.Join(records, "\n")
			log.Infof("console message: %s", ret.ConsoleMessage)
		}
	}
	//log.Infof("run command result: \n%s, \n%s", ret.StatusMessage, ret.ErrorMessage)

	return
}

func (s *Commander) terminateCommand(ev event.Event, re event.Result) {
	workspace := ev.Data().(string)
	result, err := s.TerminateCommand(workspace)
	if re != nil {
		re.Set(result, err)
	}
}

func (s *Commander) TerminateCommand(workspace string) (ret *imbcc.CommanderResult, err *cd.Error) {
	var terminateErr error
	commanderVal, commanderOK := s.workspace2Commander.Load(workspace)
	if commanderOK {
		commanderPtr := commanderVal.(*BaseCommander)
		terminateErr = commanderPtr.Terminate()
	}
	if terminateErr != nil {
		err = cd.NewError(cd.Unexpected, terminateErr.Error())
		return
	}
	ret = &imbcc.CommanderResult{
		Workspace:     workspace,
		StatusMessage: "terminated ok!",
	}
	return
}

func (s *Commander) statusCommand(ev event.Event, re event.Result) {
	workspace := ev.Data().(string)
	result, err := s.StatusCommand(workspace)
	if re != nil {
		re.Set(result, err)
	}
}

func (s *Commander) StatusCommand(workspace string) (ret *imbcc.CommanderResult, err *cd.Error) {
	var statusMessage string
	var errorMessage string = fmt.Sprintf("error: workspace%s not found", workspace)
	commanderVal, commanderOK := s.workspace2Commander.Load(workspace)
	if commanderOK {
		commanderPtr := commanderVal.(*BaseCommander)
		statusMessage, errorMessage = commanderPtr.Status()
	}
	ret = &imbcc.CommanderResult{
		Workspace:     workspace,
		StatusMessage: statusMessage,
		ErrorMessage:  errorMessage,
	}
	return
}

func (s *Commander) AppendConsoleLogs(workspace string, consoleLog string) {
	if s.consoleBuffer == nil {
		s.consoleBuffer = NewLogBuffer(defaultConsoleLinesSize)
	}

	s.consoleBuffer.Write([]byte(consoleLog))
}
