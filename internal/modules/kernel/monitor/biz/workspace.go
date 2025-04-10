package biz

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	fu "github.com/muidea/magicCommon/foundation/util"
	"github.com/muidea/magicToolkit/internal/config"
)

type FileInfo struct {
	Name     string
	Path     string
	Size     int64
	Modified int64
}

type Workspace struct {
	workspacePtr *config.Workspace
	fileInfo     sync.Map
}

func NewWorkspace(workspacePtr *config.Workspace) *Workspace {
	ret := &Workspace{
		workspacePtr: workspacePtr,
	}
	return ret
}

func (s *Workspace) AddFile(filePath string) bool {
	if !strings.HasPrefix(filePath, s.workspacePtr.Path) {
		return false
	}

	s.fileInfo.Store(filePath, nil)
	return true
}

func (s *Workspace) AddDir(dirPath string) bool {
	if !strings.HasPrefix(dirPath, s.workspacePtr.Path) {
		return false
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

func (s *Workspace) RemoveFile(filePath string) bool {
	if !strings.HasPrefix(filePath, s.workspacePtr.Path) {
		return false
	}

	s.fileInfo.Delete(filePath)
	return true
}

func (s *Workspace) ClearDir(dirPath string) bool {
	if !strings.HasPrefix(dirPath, s.workspacePtr.Path) {
		return false
	}

	s.fileInfo.Range(func(key, value interface{}) bool {
		if strings.HasPrefix(key.(string), dirPath) {
			s.fileInfo.Delete(key)
		}
		return true
	})

	return true
}

func (s *Workspace) Filter(filter fu.Filter) (ret []string) {
	s.fileInfo.Range(func(key, value interface{}) bool {
		if filter.Filter(key) {
			ret = append(ret, key.(string))
		}
		return true
	})
	return
}
