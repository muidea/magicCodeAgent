package biz

import (
	"context"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	cd "github.com/muidea/magicCommon/def"

	"github.com/muidea/magicCodeAgent/internal/config"
	imklc "github.com/muidea/magicCodeAgent/internal/modules/kernel/llm/pkg/common"
	imkwc "github.com/muidea/magicCodeAgent/internal/modules/kernel/workspace/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/cache"
	"github.com/muidea/magicCodeAgent/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
	"github.com/muidea/magicCodeAgent/pkg/util"
	"github.com/muidea/magicCommon/foundation/log"
	fp "github.com/muidea/magicCommon/foundation/path"
	fu "github.com/muidea/magicCommon/foundation/util"
)

type Observer interface {
	ID() string
	OnMessage(msgPtr *models.MessageInfo)
}

type Station struct {
	configPtr *config.Workspace
	fileInfo  sync.Map

	stationMutex sync.RWMutex
	observerList []Observer

	localMessageCache *cache.MessageSequenceCache
	localFlowManger   *FlowManager
	localWorker       Worker
	localAssistant    Assistant
}

func NewStation(workspacePtr *config.Workspace, worker Worker, assistant Assistant) *Station {
	ret := &Station{
		configPtr:         workspacePtr,
		localMessageCache: cache.NewMessageSequenceCache(500),
		localFlowManger:   NewFlowManager(),
		localWorker:       worker,
		localAssistant:    assistant,
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

func (s *Station) GetCurrentFlow() *models.Flow {
	return s.localFlowManger.GetCurrentFlow()
}

func (s *Station) GetFlow(id string) *models.Flow {
	return s.localFlowManger.GetFlow(id)
}

func (s *Station) ChatMessage(ctx context.Context, paramPtr *imkwc.ChatMessage) (err *cd.Error) {
	messageItemPtr := paramPtr.ToMessageItem()
	messageInfoPtr := &models.MessageInfo{
		Item:       messageItemPtr,
		CreateTime: time.Now().Local().Format(fu.CSTLayoutWithMillisecond),
		Workspace:  s.configPtr.Name,
	}
	var curFlow *models.Flow
	flowIDVal := ctx.Value(common.FlowID{})
	if flowIDVal != nil {
		curFlow = s.GetFlow(flowIDVal.(string))
	} else {
		curFlow = s.GetCurrentFlow()
	}
	if curFlow == nil {
		curFlow = s.localFlowManger.NewFlow(messageInfoPtr, nil)
	}

	s.DisposeSendMessage(messageInfoPtr, curFlow)

	contentContext := context.WithValue(ctx, common.ContentSink{}, s)
	providerContext := context.WithValue(contentContext, imklc.DeepseekProvider{}, imklc.RemoteDeepseek)

	msgErr := s.localAssistant.Chat(providerContext, s.configPtr.Name, messageItemPtr, curFlow)
	if msgErr != nil {
		err = msgErr
		log.Errorf("ChatMessage, err:%v", msgErr)
		return
	}

	return
}

func (s *Station) DisposeRecvMessage(msgPtr *models.MessageInfo) {
	msgPtr = s.localMessageCache.Add(msgPtr)

	curFlow := s.localFlowManger.GetCurrentFlow()
	if curFlow != nil {
		curFlow.AppendRecv(msgPtr)
	}

	//log.Infof("DisposeRecvMessage, :%v", msgPtr)
	s.stationMutex.RLock()
	defer s.stationMutex.RUnlock()

	for _, val := range s.observerList {
		val.OnMessage(msgPtr)
	}
}

func (s *Station) DisposeSendMessage(msgPtr *models.MessageInfo, flowPtr *models.Flow) {
	msgPtr = s.localMessageCache.Add(msgPtr)

	curFlow := s.localFlowManger.GetCurrentFlow()
	if curFlow == nil {
		s.localFlowManger.NewFlow(msgPtr, flowPtr)
	} else {
		curFlow.AppendSend(msgPtr)
	}

	//log.Infof("DisposeSendMessage, :%v", msgPtr)
	s.stationMutex.RLock()
	defer s.stationMutex.RUnlock()

	for _, val := range s.observerList {
		val.OnMessage(msgPtr)
	}
}

func (s *Station) FilterMessage(filter *common.MessageFilter) (ret []*models.MessageInfo, err *cd.Error) {
	messageList, messageErr := s.localMessageCache.FilterMessage(filter)
	if messageErr != nil {
		err = cd.NewError(cd.Unexpected, messageErr.Error())
		return
	}

	ret = messageList
	return
}

func (s *Station) OnClose() {
	msgPtr := &models.MessageInfo{
		Item: &models.MessageItem{
			Role:    common.Recorder,
			Content: "terminate",
		},
		CreateTime: time.Now().Local().Format(fu.CSTLayoutWithMillisecond),
		Workspace:  s.configPtr.Name,
	}

	s.DisposeRecvMessage(msgPtr)
}

func (s *Station) OnRecv(event string, content string) {
	msgPtr := &models.MessageInfo{
		Item: &models.MessageItem{
			Role:    common.Assistant,
			Content: content,
		},
		CreateTime: time.Now().Local().Format(fu.CSTLayoutWithMillisecond),
		Workspace:  s.configPtr.Name,
	}

	s.DisposeRecvMessage(msgPtr)

	blockList, blockErr := util.SpliteBlocks(content)
	if blockErr != nil {
		log.Errorf("OnRecv, err:%v", blockErr)
		return
	}

	for _, block := range blockList {
		resultVal, resultFinal := s.localWorker.DoWork(s.configPtr.Path, block)
		if resultFinal || resultVal == "" {
			continue
		}

		s.ChatMessage(context.Background(), &imkwc.ChatMessage{
			Workspace: s.configPtr.Name,
			Content:   resultVal,
		})
	}
}
