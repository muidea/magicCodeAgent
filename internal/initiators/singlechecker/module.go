package singlechecker

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/log"
	"github.com/muidea/magicCommon/framework/plugin/initiator"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicCodeAgent/internal/config"
	iiscfg "github.com/muidea/magicCodeAgent/internal/initiators/singlechecker/config"
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
		log.Errorf("open single check file failed, error:%s", fileErr.Error())
		os.Exit(1)
	}
	sysErr := syscall.Flock(int(fileHandle.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if sysErr != nil {
		if iiscfg.IsForceRestart() {
			// 读取锁文件内容获取PID
			pidBytes, readErr := os.ReadFile(singleCheckFile)
			if readErr != nil {
				log.Errorf("read single check file failed, error:%s", readErr.Error())
				os.Exit(1)
			}
			pidStr := string(pidBytes)
			pid, parseErr := strconv.Atoi(pidStr)
			if parseErr == nil {
				if process, err := os.FindProcess(pid); err == nil {
					process.Signal(syscall.SIGKILL)
				}

				for i := 0; i < 10; i++ {
					time.Sleep(time.Second)
					if _, err := os.Stat(fmt.Sprintf("/proc/%d", pid)); os.IsNotExist(err) {
						break
					}
				}
			}
			// 重新尝试获取锁
			fileHandle.Close()
			fileHandle, fileErr = os.OpenFile(singleCheckFile, os.O_CREATE|os.O_RDWR, 0666)
			if fileErr != nil {
				log.Errorf("open single_check.lock failed, error:%s", fileErr.Error())
				os.Exit(1)
			}
			sysErr = syscall.Flock(int(fileHandle.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
		}

		if sysErr != nil {
			log.Errorf("lock single_check.lock failed, error:%s", sysErr.Error())
			os.Exit(1)
		}
	}

	fileErr = fileHandle.Truncate(0)
	if fileErr == nil {
		_, fileErr = fileHandle.WriteAt([]byte(strconv.Itoa(os.Getpid())), 0)
		if fileErr != nil {
			log.Errorf("write single_check.lock failed, error:%s", fileErr.Error())
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
