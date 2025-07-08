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
	iipc "github.com/muidea/magicCodeAgent/internal/initiators/prompts/pkg/common"
	"github.com/muidea/magicCodeAgent/internal/modules/base/biz"
	imkwc "github.com/muidea/magicCodeAgent/internal/modules/kernel/workspace/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/common"
)

type Workspace struct {
	biz.Base

	fsWatcher             *fsnotify.Watcher
	workspaceStationCache sync.Map

	promptsRegistry iipc.PromptsRegistry
}

func New(
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine,
	promptsRegistry iipc.PromptsRegistry,
) *Workspace {
	workspacePtr, workspaceErr := fsnotify.NewWatcher()
	if workspaceErr != nil {
		log.Errorf("new workspace failed, err:%s", workspaceErr.Error())
		return nil
	}

	ptr := &Workspace{
		Base:            biz.New(imkwc.WorkspaceModule, eventHub, backgroundRoutine),
		fsWatcher:       workspacePtr,
		promptsRegistry: promptsRegistry,
	}

	ptr.SubscribeFunc(imkwc.QueryWorkspace, ptr.queryWorkspace)
	ptr.SubscribeFunc(imkwc.CreateWorkspace, ptr.createWorkspace)
	ptr.SubscribeFunc(imkwc.DestroyWorkspace, ptr.destroyWorkspace)
	ptr.SubscribeFunc(imkwc.AddWorkspace, ptr.addWorkspace)
	ptr.SubscribeFunc(imkwc.UpdateWorkspace, ptr.updateWorkspace)
	ptr.SubscribeFunc(imkwc.RemoveWorkspace, ptr.removeWorkspace)
	ptr.SubscribeFunc(imkwc.RefreshWorkspace, ptr.refreshWorkspace)
	ptr.SubscribeFunc(imkwc.FilterWorkspace, ptr.filterWorkspace)
	ptr.SubscribeFunc(imkwc.FilterWorkspaceFile, ptr.filterWorkspaceFile)

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
			//time.AfterFunc(500*time.Millisecond, func() {
			switch event.Op {
			case fsnotify.Create:
				s.disposeCreate(event)
			case fsnotify.Remove:
				s.disposeRemove(event)
			case fsnotify.Rename:
				s.disposeRename(event)
			case fsnotify.Write:
				s.disposeWrite(event)
			}
			//})
		case err := <-s.fsWatcher.Errors:
			log.Errorf("Workspace error: %v", err)
		}
	}
}

func (s *Workspace) disposeCreate(event fsnotify.Event) {
	if fp.IsDir(event.Name) {
		s.workspaceStationCache.Range(func(key, value interface{}) bool {
			workspaceVal, workspaceOK := value.(*Station)
			if !workspaceOK {
				return true
			}
			if workspaceVal.AddDir(event.Name) {
				return false
			}
			return true
		})
	} else {
		s.workspaceStationCache.Range(func(key, value interface{}) bool {
			workspaceVal, workspaceOK := value.(*Station)
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
		s.workspaceStationCache.Range(func(key, value interface{}) bool {
			workspaceVal, workspaceOK := value.(*Station)
			if !workspaceOK {
				return true
			}

			if workspaceVal.ClearDir(event.Name) {
				return false
			}
			return true
		})
	} else {
		s.workspaceStationCache.Range(func(key, value interface{}) bool {
			workspaceVal, workspaceOK := value.(*Station)
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

func (s *Workspace) disposeRename(event fsnotify.Event) {
	if fp.IsDir(event.Name) {
		s.workspaceStationCache.Range(func(key, value interface{}) bool {
			workspaceVal, workspaceOK := value.(*Station)
			if !workspaceOK {
				return true
			}

			if workspaceVal.ClearDir(event.Name) {
				return false
			}
			return true
		})
	} else {
		s.workspaceStationCache.Range(func(key, value interface{}) bool {
			workspaceVal, workspaceOK := value.(*Station)
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

func (s *Workspace) disposeWrite(event fsnotify.Event) {
	s.workspaceStationCache.Range(func(key, value interface{}) bool {
		workspaceVal, workspaceOK := value.(*Station)
		if !workspaceOK {
			return true
		}
		if workspaceVal.WriteFile(event.Name) {
			return false
		}
		return true
	})
}

func (s *Workspace) queryWorkspace(ev event.Event, re event.Result) {
	workspaceVal, workspaceOK := ev.Data().(*config.Workspace)
	if !workspaceOK {
		log.Errorf("query workspace failed, invalid workspace:%v", ev.Data())
		return
	}
	workspacePtr, workspaceErr := s.QueryWorkspace(workspaceVal.Name)
	if re != nil {
		re.Set(workspacePtr, workspaceErr)
	}
}

func (s *Workspace) QueryWorkspace(workspaceName string) (ret *config.Workspace, err *cd.Error) {
	workspacePtr := config.GetWorkspace(workspaceName)
	if workspacePtr == nil {
		log.Warnf("query workspace, workspace not exists:%s", workspaceName)
		err = cd.NewError(cd.Unexpected, fmt.Sprintf("query workspace, workspace not exists:%s", workspaceName))
		return
	}

	ret = workspacePtr
	return
}

func (s *Workspace) createWorkspace(ev event.Event, re event.Result) {
	workspaceVal, workspaceOK := ev.Data().(*config.Workspace)
	if !workspaceOK {
		log.Errorf("create workspace failed, invalid workspace:%v", ev.Data())
		return
	}
	workspacePtr, workspaceErr := s.CreateWorkspace(workspaceVal)
	if re != nil {
		re.Set(workspacePtr, workspaceErr)
	}
}

func (s *Workspace) CreateWorkspace(workspacePtr *config.Workspace) (ret *config.Workspace, err *cd.Error) {
	// 检查路径是否为空目录
	if fp.Exist(workspacePtr.Path) {
		emptyVal, emptyErr := fp.IsDirEmpty(workspacePtr.Path)
		if !emptyVal || emptyErr != nil {
			err = cd.NewError(cd.Unexpected, fmt.Sprintf("workspace path:%s is not empty", workspacePtr.Path))
			return
		}
	} else {
		// 创建目录
		if e := os.MkdirAll(workspacePtr.Path, 0755); e != nil {
			err = cd.NewError(cd.Unexpected, e.Error())
			return
		}
	}

	workspacePtr.Status = config.EnableStatus
	workspacePtr.UpdateTime = time.Now().Format(fu.CSTLayout)
	if eVal := config.AddWorkspace(workspacePtr); eVal != nil {
		log.Errorf("create workspace failed, add workspace failed:%s", workspacePtr.Path)
		err = cd.NewError(cd.Unexpected, eVal.Error())
		return
	}

	err = s.enableWorkspaceInternal(workspacePtr)
	if err != nil {
		return
	}

	ret = workspacePtr
	return
}

func (s *Workspace) destroyWorkspace(ev event.Event, re event.Result) {
	workspaceVal, workspaceOK := ev.Data().(*config.Workspace)
	if !workspaceOK {
		log.Errorf("destroy workspace failed, invalid workspace:%v", ev.Data())
		return
	}
	workspaceErr := s.DestroyWorkspace(workspaceVal.Name)
	if re != nil {
		re.Set(nil, workspaceErr)
	}
}

func (s *Workspace) DestroyWorkspace(workspaceName string) (err *cd.Error) {
	workspacePtr := config.GetWorkspace(workspaceName)
	if workspacePtr == nil {
		log.Warnf("destroy workspace, workspace not exists:%s", workspaceName)
		return
	}

	err = s.disableWorkspaceInternal(workspacePtr)
	if err != nil {
		return
	}
	os.RemoveAll(workspacePtr.Path)
	config.DelWorkspace(workspaceName)
	return
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
	workspace, ok := s.workspaceStationCache.Load(workspacePtr.Name)
	if ok {
		ret = workspace.(*Station).configPtr
		return
	}

	if fp.Exist(workspacePtr.Path) {
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
	curWorkspace, ok := s.workspaceStationCache.Load(workspaceName)
	if ok {
		workspacePtr := curWorkspace.(*Station).configPtr
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
	if wErr != nil {
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
	if !fp.Exist(workspacePtr.Path) || !workspacePtr.Available() {
		log.Warnf("enable workspace, workspace not exists:%s", workspacePtr.Path)
		return
	}

	workspace := NewStation(
		&s.Base,
		workspacePtr, s, s, s.promptsRegistry)
	s.workspaceStationCache.Store(workspacePtr.Name, workspace)

	if aErr := s.fsWatcher.Add(workspacePtr.Path); aErr != nil {
		log.Errorf("enable workspace failed, add workspace watcher:%s failed, error:%v", workspacePtr.Path, aErr)
	}

	filepath.Walk(workspacePtr.Path, func(path string, info os.FileInfo, err error) error {
		if isIgnored(path) {
			return nil
		}

		if info.IsDir() {
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
	if !fp.Exist(workspacePtr.Path) {
		log.Warnf("remove workspace, workspace not exists:%s", workspacePtr.Path)
		s.workspaceStationCache.Delete(workspacePtr.Name)
		return
	}

	if rErr := s.fsWatcher.Remove(workspacePtr.Path); rErr != nil {
		log.Errorf("disable workspace failed, remove workspace watcher:%s failed, error:%v", workspacePtr.Path, rErr)
	}
	curWorkspace, ok := s.workspaceStationCache.Load(workspacePtr.Name)
	filepath.Walk(workspacePtr.Path, func(path string, info os.FileInfo, err error) error {
		if isIgnored(path) {
			return nil
		}

		if info.IsDir() {
			return nil
		}
		if ok {
			curWorkspace.(*Station).RemoveFile(path)
		}
		return nil
	})

	if ok {
		s.workspaceStationCache.Delete(workspacePtr.Name)
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
	for _, val := range ret {
		if !common.ValidStage(val.Stage) {
			val.Stage = common.DesignStage
		}
	}
	return
}
