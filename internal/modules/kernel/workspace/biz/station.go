package biz

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/muidea/magicCodeAgent/internal/config"
	fp "github.com/muidea/magicCommon/foundation/path"
	fu "github.com/muidea/magicCommon/foundation/util"
)

type Observer interface {
	ID() string
	OnClose(id string)
	OnSend(id string, data []byte)
	OnRecv(id string, data []byte)
}

type Station struct {
	configPtr *config.Workspace
	fileInfo  sync.Map

	stationMutex sync.RWMutex
	observerList []Observer
}

func NewWorkspace(workspacePtr *config.Workspace) *Station {
	ret := &Station{
		configPtr: workspacePtr,
	}
	return ret
}

func (s *Station) Available() bool {
	return s.configPtr.Available()
}

func (s *Station) checkIgnore(fileFullPath string) bool {
	filePath, fileName := path.Split(fileFullPath)
	for _, ignore := range s.configPtr.Ignores {
		ignorePath := filepath.Join(s.configPtr.Path, ignore)
		if fp.IsDir(ignorePath) {
			if strings.HasPrefix(filePath, ignorePath) {
				return true
			}

			continue
		}

		if fileName == ignore {
			return true
		}
	}
	return false
}

func (s *Station) AddFile(filePath string) bool {
	if !strings.HasPrefix(filePath, s.configPtr.Path) {
		return false
	}

	if s.checkIgnore(filePath) {
		return true
	}

	s.fileInfo.Store(filePath, nil)
	return true
}

func (s *Station) AddDir(dirPath string) bool {
	if !strings.HasPrefix(dirPath, s.configPtr.Path) {
		return false
	}

	if s.checkIgnore(dirPath) {
		return true
	}

	filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if isIgnored(path) {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		s.fileInfo.Store(path, nil)
		return nil
	})

	return true
}

func (s *Station) RemoveFile(filePath string) bool {
	if !strings.HasPrefix(filePath, s.configPtr.Path) {
		return false
	}

	if s.checkIgnore(filePath) {
		return true
	}

	s.fileInfo.Delete(filePath)
	return true
}

func (s *Station) ClearDir(dirPath string) bool {
	if !strings.HasPrefix(dirPath, s.configPtr.Path) {
		return false
	}

	if s.checkIgnore(dirPath) {
		return true
	}

	s.fileInfo.Range(func(key, value interface{}) bool {
		if strings.HasPrefix(key.(string), dirPath) {
			s.fileInfo.Delete(key)
		}
		return true
	})

	return true
}

func (s *Station) Filter(filter fu.Filter) (ret []string) {
	s.fileInfo.Range(func(key, value interface{}) bool {
		if filter != nil && filter.Filter(key) {
			ret = append(ret, key.(string))
		} else if filter == nil {
			ret = append(ret, key.(string))
		}

		return true
	})
	return
}

func (s *Station) BindObserver(observer Observer) {
	s.stationMutex.Lock()
	defer s.stationMutex.Unlock()

	s.observerList = append(s.observerList, observer)
}

func (s *Station) UnbindObserver(observer Observer) {
	s.stationMutex.Lock()
	defer s.stationMutex.Unlock()

	for index, val := range s.observerList {
		if val.ID() == observer.ID() {
			s.observerList = append(s.observerList[:index], s.observerList[index+1:]...)
			break
		}
	}
}

// 处理agent的sse数据流
func (s *Station) OnClose() {

}

func (s *Station) OnRecv(event string, data []byte) error {
	return nil
}
