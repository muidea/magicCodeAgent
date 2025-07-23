package singlechecker

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"syscall"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/framework/plugin/initiator"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicCodeAgent/internal/config"
	"github.com/muidea/magicCodeAgent/internal/initiators/singlechecker/pkg/common"
)

func init() {
	initiator.Register(New())
}

type Singlechecker struct {
	checkFile *os.File
}

func New() *Singlechecker {
	return &Singlechecker{}
}

func (s *Singlechecker) ID() string {
	return common.SinglecheckerInitiator
}

func (s *Singlechecker) Weight() int {
	return 0
}

func (s *Singlechecker) Setup(eventHub event.Hub, backgroundRoutine task.BackgroundRoutine) (err *cd.Error) {
	singleCheckPath := filepath.Join(config.RootConfigDir(), ".temp")
	if _, oErr := os.Stat(singleCheckPath); oErr != nil {
		if !os.IsNotExist(oErr) {
			return cd.NewError(cd.Unexpected, "create singleton check path failed")
		}
		if oErr := os.MkdirAll(singleCheckPath, 0755); oErr != nil {
			return cd.NewError(cd.Unexpected, "create singleton check path failed")
		}
	}

	singleCheckFile := filepath.Join(singleCheckPath, fmt.Sprintf("%s.lock", config.EndpointName()))
	fileHandle, fileErr := os.OpenFile(singleCheckFile, os.O_CREATE|os.O_RDWR, 0666)
	if fileErr != nil {
		os.Exit(1)
	}
	sysErr := syscall.Flock(int(fileHandle.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if sysErr != nil {
		os.Exit(1)
	}

	fileErr = fileHandle.Truncate(0)
	if fileErr == nil {
		_, fileErr = fileHandle.WriteAt([]byte(strconv.Itoa(os.Getpid())), 0)
		if fileErr != nil {
			os.Exit(1)
		}
	}
	s.checkFile = fileHandle
	return nil
}

func (s *Singlechecker) Run() (err *cd.Error) {
	return
}

func (s *Singlechecker) Teardown() {
	if s.checkFile == nil {
		return
	}
	syscall.Flock(int(s.checkFile.Fd()), syscall.LOCK_UN)
	s.checkFile.Close()
}
