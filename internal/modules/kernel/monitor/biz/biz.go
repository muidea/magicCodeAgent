package biz

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/log"
	fu "github.com/muidea/magicCommon/foundation/util"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicToolkit/internal/config"
	"github.com/muidea/magicToolkit/internal/modules/base/biz"
	"github.com/muidea/magicToolkit/internal/modules/kernel/monitor/pkg/common"
	"github.com/muidea/magicToolkit/pkg/util"
)

type Monitor struct {
	biz.Base

	fsWatcher *fsnotify.Watcher
	fileCache sync.Map
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

	ptr.SubscribeFunc(common.AddWorkSpace, ptr.addWorkspace)
	ptr.SubscribeFunc(common.RemoveWorkSpace, ptr.removeWorkspace)
	ptr.SubscribeFunc(common.RefreshWorkSpace, ptr.refreshWorkspace)

	return ptr
}

func (s *Monitor) Run() (err *cd.Error) {
	s.Base.AsyncTask(s.monitorTask)

	workspace := config.Workspace()
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
				if fsnotify.Create == event.Op {
					s.fileCache.Store(event.Name, event.Name)
				} else {
					s.fileCache.Delete(event.Name)
				}
			})
		case err := <-s.fsWatcher.Errors:
			log.Errorf("Monitor error: %v", err)
		}
	}
}

func (s *Monitor) addWorkspace(ev event.Event, re event.Result) {
	workspaceVal, workspaceOK := ev.Data().(string)
	if !workspaceOK {
		log.Errorf("add workspace failed, invalid workspace:%v", ev.Data())
		return
	}
	err := s.addWorkspaceInternal(workspaceVal)
	if re != nil {
		re.Set(nil, err)
	}
}

func (s *Monitor) addWorkspaceInternal(workspace string) (err *cd.Error) {
	if !util.Exists(workspace) {
		log.Errorf("add workspace failed, workspace not exists:%s", workspace)
		return cd.NewError(cd.Unexpected, "workspace not exists")
	}

	filepath.Walk(workspace, func(path string, info os.FileInfo, err error) error {
		if isIgnored(path) {
			return nil
		}

		if info.IsDir() {
			s.fsWatcher.Add(path)
			return nil
		}

		s.fileCache.Store(path, path)
		return nil
	})
	return
}

func (s *Monitor) removeWorkspace(ev event.Event, re event.Result) {
	workspaceVal, workspaceOK := ev.Data().(string)
	if !workspaceOK {
		log.Errorf("remove workspace failed, invalid workspace:%v", ev.Data())
		return
	}
	if !util.Exists(workspaceVal) {
		log.Errorf("remove workspace failed, workspace not exists:%s", workspaceVal)
		if re != nil {
			re.Set(nil, cd.NewError(cd.Unexpected, "workspace not exists"))
		}
		return
	}

	filepath.Walk(workspaceVal, func(path string, info os.FileInfo, err error) error {
		if isIgnored(path) {
			return nil
		}

		if info.IsDir() {
			s.fsWatcher.Remove(path)
			return nil
		}

		s.fileCache.Delete(path)
		return nil
	})
	if re != nil {
		re.Set(nil, nil)
	}
}

func (s *Monitor) refreshWorkspace(ev event.Event, re event.Result) {
}

func (s *Monitor) FilterWorkSpaceFileList(filter fu.Filter) (ret []string, err *cd.Error) {
	s.fileCache.Range(func(key, value interface{}) bool {
		if filter.Filter(key) {
			ret = append(ret, key.(string))
		}
		return true
	})
	return
}
