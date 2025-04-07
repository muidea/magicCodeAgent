package biz

import (
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/log"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicToolkit/internal/modules/base/biz"
	"github.com/muidea/magicToolkit/internal/modules/kernel/watcher/pkg/common"
	"github.com/muidea/magicToolkit/pkg/util"
)

type Watcher struct {
	biz.Base

	fsWatcher *fsnotify.Watcher
	eventChan chan HybridEvent
}

func New(
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine,
) *Watcher {
	watcherPtr, watcherErr := fsnotify.NewWatcher()
	if watcherErr != nil {
		log.Errorf("new watcher failed, err:%s", watcherErr.Error())
		return nil
	}

	ptr := &Watcher{
		Base:      biz.New(common.WatcherModule, eventHub, backgroundRoutine),
		fsWatcher: watcherPtr,
		eventChan: make(chan HybridEvent, 1024),
	}

	ptr.SubscribeFunc(common.AddWorkSpace, ptr.addWorkspace)
	ptr.SubscribeFunc(common.RemoveWorkSpace, ptr.removeWorkspace)
	ptr.SubscribeFunc(common.RefreshWorkSpace, ptr.refreshWorkspace)

	return ptr
}

func (s *Watcher) Run() (err *cd.Error) {
	s.Base.AsyncTask(func() {
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
					s.eventChan <- HybridEvent{
						EventType: "filesystem",
						Path:      event.Name,
						Operation: map[fsnotify.Op]string{
							fsnotify.Create: "create",
							fsnotify.Write:  "update",
							fsnotify.Remove: "delete",
						}[event.Op],
						Timestamp: time.Now(),
					}
				})
			case err := <-s.fsWatcher.Errors:
				log.Errorf("Watcher error: %v", err)
			}
		}
	})
	return
}

func (s *Watcher) addWorkspace(ev event.Event, re event.Result) {
	workspaceVal, workspaceOK := ev.Data().(string)
	if !workspaceOK {
		log.Errorf("add workspace failed, invalid workspace:%v", ev.Data())
		return
	}

	if !util.Exists(workspaceVal) {
		log.Errorf("add workspace failed, workspace not exists:%s", workspaceVal)
		if re != nil {
			re.Set(nil, cd.NewError(cd.Unexpected, "workspace not exists"))
		}
		return
	}

	filepath.Walk(workspaceVal, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && !isIgnored(path) {
			s.fsWatcher.Add(path)
		}
		return nil
	})

	if re != nil {
		re.Set(nil, nil)
	}
}

func (s *Watcher) removeWorkspace(ev event.Event, re event.Result) {
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
		if info.IsDir() {
			s.fsWatcher.Remove(path)
		}
		return nil
	})
	if re != nil {
		re.Set(nil, nil)
	}
}

func (s *Watcher) refreshWorkspace(ev event.Event, re event.Result) {
}
