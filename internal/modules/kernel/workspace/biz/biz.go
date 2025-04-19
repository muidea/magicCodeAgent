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

	"github.com/muidea/magicCodeAgent/internal/config"
	"github.com/muidea/magicCodeAgent/internal/modules/base/biz"
	"github.com/muidea/magicCodeAgent/internal/modules/kernel/workspace/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/util"
)

type Workspace struct {
	biz.Base

	fsWatcher      *fsnotify.Watcher
	workspaceCache sync.Map
}

func New(
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine,
) *Workspace {
	workspacePtr, workspaceErr := fsnotify.NewWatcher()
	if workspaceErr != nil {
		log.Errorf("new workspace failed, err:%s", workspaceErr.Error())
		return nil
	}

	ptr := &Workspace{
		Base:      biz.New(common.WorkspaceModule, eventHub, backgroundRoutine),
		fsWatcher: workspacePtr,
	}

	ptr.SubscribeFunc(common.AddWorkspace, ptr.addWorkspace)
	ptr.SubscribeFunc(common.UpdateWorkspace, ptr.updateWorkspace)
	ptr.SubscribeFunc(common.RemoveWorkspace, ptr.removeWorkspace)
	ptr.SubscribeFunc(common.RefreshWorkspace, ptr.refreshWorkspace)
	ptr.SubscribeFunc(common.FilterWorkspace, ptr.filterWorkspace)
	ptr.SubscribeFunc(common.FilterWorkspaceFile, ptr.filterWorkspaceFile)

	return ptr
}

func (s *Workspace) Run() (err *cd.Error) {
	s.Base.AsyncTask(s.workspaceTask)

	workspace := config.Workspaces()
	for _, v := range workspace {
		s.AddWorkspace(v)
	}

	return
}

func (s *Workspace) workspaceTask() {
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
			log.Errorf("Workspace error: %v", err)
		}
	}
}

func (s *Workspace) disposeCreate(event fsnotify.Event) {
	if fp.IsDir(event.Name) {
		s.fsWatcher.Add(event.Name)
		s.workspaceCache.Range(func(key, value interface{}) bool {
			workspaceVal, workspaceOK := value.(*WorkspaceInfo)
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
			workspaceVal, workspaceOK := value.(*WorkspaceInfo)
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

func (s *Workspace) disposeRemove(event fsnotify.Event) {
	if fp.IsDir(event.Name) {
		s.fsWatcher.Remove(event.Name)
		s.workspaceCache.Range(func(key, value interface{}) bool {
			workspaceVal, workspaceOK := value.(*WorkspaceInfo)
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
			workspaceVal, workspaceOK := value.(*WorkspaceInfo)
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

func (s *Workspace) addWorkspace(ev event.Event, re event.Result) {
	workspaceVal, workspaceOK := ev.Data().(*config.Workspace)
	if !workspaceOK {
		log.Errorf("add workspace failed, invalid workspace:%v", ev.Data())
		return
	}
	workspacePtr, workspaceErr := s.AddWorkspace(workspaceVal)
	if re != nil {
		re.Set(workspacePtr, workspaceErr)
	}
}

func (s *Workspace) AddWorkspace(workspacePtr *config.Workspace) (ret *config.Workspace, err *cd.Error) {
	workspace, ok := s.workspaceCache.Load(workspacePtr.Name)
	if ok {
		ret = workspace.(*WorkspaceInfo).workspacePtr
		return
	}

	if util.Exists(workspacePtr.Path) {
		workspacePtr.Status = config.EnableStatus
	} else {
		workspacePtr.Status = config.OfflineStatus
	}

	workspacePtr.UpdateTime = time.Now().Format(fu.CSTLayout)
	if eVal := config.AddWorkspace(workspacePtr); eVal != nil {
		log.Errorf("add workspace failed, add workspace failed:%s", workspacePtr.Path)
		err = cd.NewError(cd.Unexpected, eVal.Error())
		return
	}

	if workspacePtr.Status == config.EnableStatus {
		err = s.enableWorkspaceInternal(workspacePtr)
	}
	if err != nil {
		return
	}

	ret = workspacePtr
	return
}

func (s *Workspace) updateWorkspace(ev event.Event, re event.Result) {
	workspaceVal, workspaceOK := ev.Data().(*config.Workspace)
	if !workspaceOK {
		log.Errorf("update workspace failed, invalid workspace:%v", ev.Data())
		return
	}

	err := s.UpdateWorkspace(workspaceVal)
	if re != nil {
		re.Set(nil, err)
	}
}

func (s *Workspace) UpdateWorkspace(workspacePtr *config.Workspace) (err *cd.Error) {
	if eVal := config.UpdateWorkspace(workspacePtr); eVal != nil {
		log.Errorf("update workspace failed, update workspace failed:%s", workspacePtr.Path)
		err = cd.NewError(cd.Unexpected, eVal.Error())
		return
	}

	err = s.RefreshWorkspace(workspacePtr.Name)
	return
}

func (s *Workspace) removeWorkspace(ev event.Event, re event.Result) {
	workspaceName, workspaceOK := ev.Data().(string)
	if !workspaceOK {
		log.Errorf("remove workspace failed, invalid workspace:%v", ev.Data())
		return
	}

	err := s.RemoveWorkspace(workspaceName)
	if re != nil {
		re.Set(nil, err)
	}
}

func (s *Workspace) RemoveWorkspace(workspaceName string) (err *cd.Error) {
	curWorkspace, ok := s.workspaceCache.Load(workspaceName)
	if ok {
		workspacePtr := curWorkspace.(*WorkspaceInfo).workspacePtr
		err = s.disableWorkspaceInternal(workspacePtr)
		if err != nil {
			return
		}
	}

	if eVal := config.DelWorkspace(workspaceName); eVal != nil {
		log.Errorf("remove workspace failed, remove workspace:%s failed", workspaceName)
		return cd.NewError(cd.Unexpected, eVal.Error())
	}
	return
}

func (s *Workspace) refreshWorkspace(ev event.Event, re event.Result) {
	workspaceVal, workspaceOK := ev.Data().(*config.Workspace)
	if !workspaceOK {
		log.Errorf("refresh workspace failed, invalid workspace:%v", ev.Data())
		return
	}

	err := s.RefreshWorkspace(workspaceVal.Name)
	if re != nil {
		re.Set(nil, err)
	}
}

func (s *Workspace) RefreshWorkspace(workspaceName string) (err *cd.Error) {
	wPtr, wErr := config.RefreshWorkspace(workspaceName)
	if wErr == nil {
		log.Warnf("refresh workspace, workspace not exists:%s", workspaceName)
		return
	}

	if wPtr.Status != config.EnableStatus {
		log.Warnf("refresh workspace, workspace not enable:%s", wPtr.Path)
		return
	}

	err = s.disableWorkspaceInternal(wPtr)
	if err != nil {
		return
	}

	err = s.enableWorkspaceInternal(wPtr)
	return
}

func (s *Workspace) EnableWorkspace(workspaceName string) (err *cd.Error) {
	wPtr := config.GetWorkspace(workspaceName)
	if wPtr == nil {
		log.Warnf("enable workspace, workspace not exists:%s", workspaceName)
		return
	}

	err = s.enableWorkspaceInternal(wPtr)
	if err != nil {
		log.Errorf("enable workspace failed, enable workspace failed:%s", wPtr.Path)
		return
	}

	wPtr.Status = config.EnableStatus
	wPtr.UpdateTime = time.Now().Format(fu.CSTLayout)
	if eVal := config.EnableWorkspace(workspaceName); eVal != nil {
		log.Errorf("enable workspace failed, enable workspace failed:%s", wPtr.Path)
		return cd.NewError(cd.Unexpected, eVal.Error())
	}

	return
}

func (s *Workspace) enableWorkspaceInternal(workspacePtr *config.Workspace) (err *cd.Error) {
	if !util.Exists(workspacePtr.Path) {
		log.Warnf("enable workspace, workspace not exists:%s", workspacePtr.Path)
		return
	}

	workspace := NewWorkspace(workspacePtr)
	s.workspaceCache.Store(workspacePtr.Name, workspace)

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

func (s *Workspace) DisableWorkspace(workspaceName string) (err *cd.Error) {
	wPtr := config.GetWorkspace(workspaceName)
	if wPtr == nil {
		log.Warnf("disable workspace, workspace not exists:%s", workspaceName)
		return
	}

	err = s.disableWorkspaceInternal(wPtr)
	if err != nil {
		return
	}

	wPtr.Status = config.DisableStatus
	wPtr.UpdateTime = time.Now().Format(fu.CSTLayout)
	if eVal := config.DisableWorkspace(workspaceName); eVal != nil {
		log.Errorf("disable workspace failed, disable workspace failed:%s", wPtr.Path)
		return cd.NewError(cd.Unexpected, eVal.Error())
	}

	return
}

func (s *Workspace) disableWorkspaceInternal(workspacePtr *config.Workspace) (err *cd.Error) {
	if !util.Exists(workspacePtr.Path) {
		log.Warnf("remove workspace, workspace not exists:%s", workspacePtr.Path)
		s.workspaceCache.Delete(workspacePtr.Name)
		return
	}

	curWorkspace, ok := s.workspaceCache.Load(workspacePtr.Name)
	filepath.Walk(workspacePtr.Path, func(path string, info os.FileInfo, err error) error {
		if isIgnored(path) {
			return nil
		}

		if info.IsDir() {
			s.fsWatcher.Remove(path)
			return nil
		}
		if ok {
			curWorkspace.(*WorkspaceInfo).RemoveFile(path)
		}
		return nil
	})

	if ok {
		s.workspaceCache.Delete(workspacePtr.Name)
	}
	return
}

func (s *Workspace) filterWorkspace(ev event.Event, re event.Result) {
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

func (s *Workspace) FilterWorkspace(filter fu.Filter) (ret []*config.Workspace, err *cd.Error) {
	ret = config.Workspaces()
	return
}

func (s *Workspace) filterWorkspaceFile(ev event.Event, re event.Result) {
	val := ev.Data()
	if val == nil {
		log.Errorf("filterWorkspaceFile failed, invalid event value")
		if re != nil {
			re.Set(nil, cd.NewError(cd.Unexpected, "illegal event value"))
		}
		return
	}
	workspaceName, workspaceOK := val.(string)
	if !workspaceOK {
		log.Errorf("filterWorkspaceFile failed, invalid event value:%v", val)
		if re != nil {
			re.Set(nil, cd.NewError(cd.Unexpected, "illegal event value"))
		}
		return
	}
	fileList, fileErr := s.FilterWorkspaceFileList(workspaceName, nil)
	if re != nil {
		re.Set(fileList, fileErr)
	}
}

func (s *Workspace) FilterWorkspaceFileList(workspaceName string, filter fu.Filter) (ret []string, err *cd.Error) {
	curWorkspace, ok := s.workspaceCache.Load(workspaceName)
	if !ok {
		err = cd.NewError(cd.NotFound, fmt.Sprintf("workspace:%s not found", workspaceName))
		return
	}

	if !fp.Exist(curWorkspace.(*WorkspaceInfo).workspacePtr.Path) {
		err = cd.NewError(cd.NotFound, fmt.Sprintf("workspace:%s, path:%s not exists", workspaceName, curWorkspace.(*WorkspaceInfo).workspacePtr.Path))
		return
	}

	ret = curWorkspace.(*WorkspaceInfo).Filter(filter)
	return
}
