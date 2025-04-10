package biz

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/log"
	fp "github.com/muidea/magicCommon/foundation/path"
	fu "github.com/muidea/magicCommon/foundation/util"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicToolkit/internal/config"
	"github.com/muidea/magicToolkit/internal/modules/base/biz"
	"github.com/muidea/magicToolkit/internal/modules/kernel/monitor/pkg/common"
	"github.com/muidea/magicToolkit/pkg/util"
)

type Monitor struct {
	biz.Base

	fsWatcher      *fsnotify.Watcher
	workspaceCache sync.Map
}

func New(
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine,
) *Monitor {
	monitorPtr, monitorErr := fsnotify.NewWatcher()
	if monitorErr != nil {
		log.Errorf("new monitor failed, err:%s", monitorErr.Error())
		return nil
	}

	ptr := &Monitor{
		Base:      biz.New(common.MonitorModule, eventHub, backgroundRoutine),
		fsWatcher: monitorPtr,
	}

	ptr.SubscribeFunc(common.AddWorkspace, ptr.addWorkspace)
	ptr.SubscribeFunc(common.RemoveWorkspace, ptr.removeWorkspace)
	ptr.SubscribeFunc(common.RefreshWorkspace, ptr.refreshWorkspace)
	ptr.SubscribeFunc(common.FilterWorkspace, ptr.filterWorkspace)

	return ptr
}

func (s *Monitor) Run() (err *cd.Error) {
	s.Base.AsyncTask(s.monitorTask)

	workspace := config.Workspaces()
	for _, v := range workspace {
		s.addWorkspaceInternal(v)
	}

	return
}

func (s *Monitor) monitorTask() {
	for {
		select {
		case event, ok := <-s.fsWatcher.Events:
			if !ok {
				return
			}
			if isIgnored(event.Name) {
				continue
			}

			// 去抖动处理（500ms窗口）
			time.AfterFunc(500*time.Millisecond, func() {
				switch event.Op {
				case fsnotify.Create:
					s.disposeCreate(event)
				case fsnotify.Remove:
					s.disposeRemove(event)
				}
			})
		case err := <-s.fsWatcher.Errors:
			log.Errorf("Monitor error: %v", err)
		}
	}
}

func (s *Monitor) disposeCreate(event fsnotify.Event) {
	if fp.IsDir(event.Name) {
		s.fsWatcher.Add(event.Name)
		s.workspaceCache.Range(func(key, value interface{}) bool {
			workspaceVal, workspaceOK := value.(*Workspace)
			if !workspaceOK {
				return true
			}
			if workspaceVal.AddDir(event.Name) {
				return false
			}
			return true
		})
	} else {
		s.workspaceCache.Range(func(key, value interface{}) bool {
			workspaceVal, workspaceOK := value.(*Workspace)
			if !workspaceOK {
				return true
			}
			if workspaceVal.AddFile(event.Name) {
				return false
			}

			return true
		})
	}
}

func (s *Monitor) disposeRemove(event fsnotify.Event) {
	if fp.IsDir(event.Name) {
		s.fsWatcher.Remove(event.Name)
		s.workspaceCache.Range(func(key, value interface{}) bool {
			workspaceVal, workspaceOK := value.(*Workspace)
			if !workspaceOK {
				return true
			}

			if workspaceVal.ClearDir(event.Name) {
				return false
			}
			return true
		})
	} else {
		s.workspaceCache.Range(func(key, value interface{}) bool {
			workspaceVal, workspaceOK := value.(*Workspace)
			if !workspaceOK {
				return true
			}
			if workspaceVal.RemoveFile(event.Name) {
				return false
			}
			return true
		})
	}
}

func (s *Monitor) addWorkspace(ev event.Event, re event.Result) {
	workspaceVal, workspaceOK := ev.Data().(*config.Workspace)
	if !workspaceOK {
		log.Errorf("add workspace failed, invalid workspace:%v", ev.Data())
		return
	}
	err := s.addWorkspaceInternal(workspaceVal)
	if re != nil {
		re.Set(nil, err)
	}
}

func (s *Monitor) addWorkspaceInternal(workspacePtr *config.Workspace) (err *cd.Error) {
	if !util.Exists(workspacePtr.Path) {
		log.Errorf("add workspace failed, workspace not exists:%s", workspacePtr.Path)
		return cd.NewError(cd.Unexpected, "workspace not exists")
	}

	workspace := NewWorkspace(workspacePtr)
	s.workspaceCache.Store(workspacePtr.Path, workspace)

	filepath.Walk(workspacePtr.Path, func(path string, info os.FileInfo, err error) error {
		if isIgnored(path) {
			return nil
		}

		if info.IsDir() {
			s.fsWatcher.Add(path)
			return nil
		}

		workspace.AddFile(path)
		return nil
	})
	return
}

func (s *Monitor) removeWorkspace(ev event.Event, re event.Result) {
	workspaceVal, workspaceOK := ev.Data().(*config.Workspace)
	if !workspaceOK {
		log.Errorf("remove workspace failed, invalid workspace:%v", ev.Data())
		return
	}

	err := s.removeWorkspaceInternal(workspaceVal)
	if re != nil {
		re.Set(nil, err)
	}
}

func (s *Monitor) removeWorkspaceInternal(workspacePtr *config.Workspace) (err *cd.Error) {
	if !util.Exists(workspacePtr.Path) {
		log.Warnf("remove workspace, workspace not exists:%s", workspacePtr.Path)
		s.workspaceCache.Delete(workspacePtr.Path)
		return
	}

	curWorkspace, ok := s.workspaceCache.Load(workspacePtr.Path)
	filepath.Walk(workspacePtr.Path, func(path string, info os.FileInfo, err error) error {
		if isIgnored(path) {
			return nil
		}

		if info.IsDir() {
			s.fsWatcher.Remove(path)
			return nil
		}
		if ok {
			curWorkspace.(*Workspace).RemoveFile(path)
		}
		return nil
	})

	s.workspaceCache.Delete(workspacePtr.Path)
	return
}

func (s *Monitor) refreshWorkspace(ev event.Event, re event.Result) {
	workspaceVal, workspaceOK := ev.Data().(*config.Workspace)
	if !workspaceOK {
		log.Errorf("refresh workspace failed, invalid workspace:%v", ev.Data())
		return
	}

	_, ok := s.workspaceCache.Load(workspaceVal.Path)
	if ok {
		err := s.removeWorkspaceInternal(workspaceVal)
		if err != nil {
			log.Errorf("refresh workspace failed, remove workspace failed:%s", workspaceVal.Path)
			return
		}

		err = s.addWorkspaceInternal(workspaceVal)
		if err != nil {
			log.Errorf("refresh workspace failed, add workspace failed:%s", workspaceVal.Path)
			return
		}
	}
}

func (s *Monitor) filterWorkspace(ev event.Event, re event.Result) {
	var filter fu.Filter
	val := ev.Data()
	if val != nil {
		filter = val.(fu.Filter)
	}

	workspaceList, workspaceErr := s.FilterWorkspace(filter)
	if re != nil {
		re.Set(workspaceList, workspaceErr)
	}
}

func (s *Monitor) FilterWorkspace(filter fu.Filter) (ret []*config.Workspace, err *cd.Error) {
	s.workspaceCache.Range(func(key, value interface{}) bool {
		workspace := value.(*config.Workspace)
		if filter == nil {
			ret = append(ret, workspace)
			return true
		}

		if filter.Filter(workspace) {
			ret = append(ret, workspace)
		}
		return true
	})
	return
}

func (s *Monitor) FilterWorkspaceFileList(workspace string, filter fu.Filter) (ret []string, err *cd.Error) {
	curWorkspace, ok := s.workspaceCache.Load(workspace)
	if !ok {
		err = cd.NewError(cd.NotFound, fmt.Sprintf("workspace:%s not found", workspace))
		return
	}

	if !fp.Exist(curWorkspace.(*config.Workspace).Path) {
		err = cd.NewError(cd.NotFound, fmt.Sprintf("workspace:%s, path:%s not exists", workspace, curWorkspace.(*config.Workspace).Path))
		return
	}

	ret = curWorkspace.(*Workspace).Filter(filter)
	return
}
