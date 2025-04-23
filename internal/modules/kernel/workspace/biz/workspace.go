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

type WorkspaceStation struct {
	workspacePtr *config.Workspace
	fileInfo     sync.Map
}

func NewWorkspace(workspacePtr *config.Workspace) *WorkspaceStation {
	ret := &WorkspaceStation{
		workspacePtr: workspacePtr,
	}
	return ret
}

func (s *WorkspaceStation) Available() bool {
	return s.workspacePtr.Available()
}

func (s *WorkspaceStation) checkIgnore(fileFullPath string) bool {
	filePath, fileName := path.Split(fileFullPath)
	for _, ignore := range s.workspacePtr.Ignores {
		ignorePath := filepath.Join(s.workspacePtr.Path, ignore)
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

func (s *WorkspaceStation) AddFile(filePath string) bool {
	if !strings.HasPrefix(filePath, s.workspacePtr.Path) {
		return false
	}

	if s.checkIgnore(filePath) {
		return true
	}

	s.fileInfo.Store(filePath, nil)
	return true
}

func (s *WorkspaceStation) AddDir(dirPath string) bool {
	if !strings.HasPrefix(dirPath, s.workspacePtr.Path) {
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

func (s *WorkspaceStation) RemoveFile(filePath string) bool {
	if !strings.HasPrefix(filePath, s.workspacePtr.Path) {
		return false
	}

	if s.checkIgnore(filePath) {
		return true
	}

	s.fileInfo.Delete(filePath)
	return true
}

func (s *WorkspaceStation) ClearDir(dirPath string) bool {
	if !strings.HasPrefix(dirPath, s.workspacePtr.Path) {
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

func (s *WorkspaceStation) Filter(filter fu.Filter) (ret []string) {
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
