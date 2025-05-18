package biz

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	cd "github.com/muidea/magicCommon/def"

	"github.com/muidea/magicCodeAgent/internal/config"
	iipc "github.com/muidea/magicCodeAgent/internal/initators/prompts/pkg/common"
	"github.com/muidea/magicCodeAgent/internal/modules/base/biz"
	imklc "github.com/muidea/magicCodeAgent/internal/modules/kernel/llm/pkg/common"
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
	bizBasePtr *biz.Base
	configPtr  *config.Workspace
	fileInfo   sync.Map

	currentWorkStage string

	stationMutex sync.RWMutex
	observerList []Observer

	localMessageCache *cache.MessageSequenceCache
	localFlowManger   *FlowManager
	localWorker       Worker
	localAssistant    Assistant
	promptsRegistry   iipc.PromptsRegistry
}

func NewStation(
	bizBasePtr *biz.Base,
	workspacePtr *config.Workspace,
	worker Worker,
	assistant Assistant,
	promptsRegistry iipc.PromptsRegistry,
) *Station {
	ret := &Station{
		bizBasePtr:        bizBasePtr,
		configPtr:         workspacePtr,
		localMessageCache: cache.NewMessageSequenceCache(500),
		localFlowManger:   NewFlowManager(),
		localWorker:       worker,
		localAssistant:    assistant,
		promptsRegistry:   promptsRegistry,
	}

	return ret
}

func (s *Station) getWorkspaceType() string {
	return s.configPtr.GetWorkspaceType()

}

func (s *Station) systemPrompt() string {
	return s.promptsRegistry.GetPromptsContent(s.configPtr.GetStage(), s.configPtr.TechStack)
}

func (s *Station) envPrompt() string {
	stringBuilder := bytes.NewBufferString("")
	stringBuilder.WriteString("<environment_details>\n")
	stringBuilder.WriteString("# 项目名称：" + s.configPtr.Name + "\n\n\n")

	stringBuilder.WriteString("## 项目描述\n")
	stringBuilder.WriteString(s.configPtr.Description + "\n\n\n")

	stringBuilder.WriteString("## 技术栈\n")
	for _, val := range s.configPtr.TechStack {
		stringBuilder.WriteString("- " + val + "\n")
	}
	stringBuilder.WriteString("\n\n")

	//## 项目规划目录
	structContent := s.promptsRegistry.ReadSectionPromptContent(fmt.Sprintf("project-%s-struct", s.getWorkspaceType()))
	if structContent != "" {
		stringBuilder.WriteString(structContent)
		stringBuilder.WriteString("\n\n")
	}

	//## 项目状态
	stageContent := s.promptsRegistry.ReadSectionPromptContent(fmt.Sprintf("stage-%s", s.configPtr.GetStage()))
	if stageContent != "" {
		stringBuilder.WriteString(stageContent)
		stringBuilder.WriteString("\n\n")
	}

	previousLabel, previousPath := s.configPtr.GetPreviousStagePath()
	if previousPath != "" {
		stringBuilder.WriteString(fmt.Sprintf("### %s\n", previousLabel))
		s.fileInfo.Range(func(key, value interface{}) bool {
			if strings.HasPrefix(key.(string), previousPath) {
				relativePath := key.(string)[len(previousPath)+1:]
				stringBuilder.WriteString("- " + relativePath + "\n")
			}
			return true
		})
		stringBuilder.WriteString("\n\n")
	}

	stagePath := s.configPtr.GetCurrentStagePath()
	stringBuilder.WriteString("### 当前文件清单\n")
	s.fileInfo.Range(func(key, value interface{}) bool {
		if strings.HasPrefix(key.(string), stagePath) {
			relativePath := key.(string)[len(stagePath)+1:]
			stringBuilder.WriteString("- " + relativePath + "\n")
		}
		return true
	})
	stringBuilder.WriteString("\n\n")

	stringBuilder.WriteString("</environment_name>")
	return stringBuilder.String()
}

func (s *Station) Available() bool {
	return s.configPtr.Available()
}

func (s *Station) checkIgnore(fileFullPath string) bool {
	if isIgnored(fileFullPath) {
		return true
	}

	stagePath := s.configPtr.GetCurrentStagePath()
	filePath, fileName := path.Split(fileFullPath)
	for _, ignore := range s.configPtr.Ignores {
		// 只需要匹配code目录下的文件
		ignorePath := filepath.Join(stagePath, ignore)
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
		if info.IsDir() {
			return nil
		}

		if s.checkIgnore(path) {
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

func (s *Station) ChatMessage(ctx context.Context, stage string, messageItemPtr *models.MessageItem) (err *cd.Error) {
	messageInfoPtr := &models.MessageInfo{
		Item:       *messageItemPtr,
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

	s.currentWorkStage = stage

	curFlow = s.DisposeSendMessage(messageInfoPtr, curFlow)

	stationSystemPrompt := s.systemPrompt()
	stationEnvPrompt := s.envPrompt()
	contentContext := context.WithValue(ctx, common.ContentSink{}, s)
	providerContext := context.WithValue(contentContext, imklc.DeepseekProvider{}, imklc.RemoteDeepseek)
	msgErr := s.localAssistant.Chat(providerContext, s.configPtr, stationSystemPrompt, stationEnvPrompt, curFlow.GetFlowPath())
	if msgErr != nil {
		err = msgErr
		log.Errorf("ChatMessage, s.localAssistant.Chat error:%v", msgErr)
		return
	}

	return
}

func (s *Station) DisposeRecvMessage(msgPtr *models.MessageInfo) {
	msgPtr = s.localMessageCache.Add(msgPtr)

	curFlow := s.localFlowManger.GetCurrentFlow()
	if curFlow != nil {
		curFlow.AppendMessage(msgPtr)
	}

	//log.Infof("DisposeRecvMessage, :%v", msgPtr)
	s.stationMutex.RLock()
	defer s.stationMutex.RUnlock()

	for _, val := range s.observerList {
		val.OnMessage(msgPtr)
	}
}

func (s *Station) DisposeSendMessage(msgPtr *models.MessageInfo, flowPtr *models.Flow) (ret *models.Flow) {
	msgPtr = s.localMessageCache.Add(msgPtr)

	curFlow := s.localFlowManger.GetCurrentFlow()
	if curFlow == nil {
		curFlow = s.localFlowManger.NewFlow(msgPtr, flowPtr)
	} else {
		curFlow.AppendMessage(msgPtr)
	}

	//log.Infof("DisposeSendMessage, :%v", msgPtr)
	s.stationMutex.RLock()
	defer s.stationMutex.RUnlock()

	for _, val := range s.observerList {
		val.OnMessage(msgPtr)
	}

	ret = curFlow
	return
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
		Item: models.MessageItem{
			Role:    common.Recorder,
			Content: []string{"terminate"},
		},
		CreateTime: time.Now().Local().Format(fu.CSTLayoutWithMillisecond),
		Workspace:  s.configPtr.Name,
	}

	s.DisposeRecvMessage(msgPtr)
}

func (s *Station) OnRecv(event string, content string) {
	msgPtr := &models.MessageInfo{
		Item: models.MessageItem{
			Role:    common.Assistant,
			Content: []string{content},
		},
		CreateTime: time.Now().Local().Format(fu.CSTLayoutWithMillisecond),
		Workspace:  s.configPtr.Name,
	}

	s.DisposeRecvMessage(msgPtr)

	s.bizBasePtr.AsyncTask(func() {
		blockList, blockErr := util.SpliteBlocks(content)
		if blockErr != nil {
			log.Errorf("OnRecv, err:%v", blockErr)
			return
		}

		for _, block := range blockList {
			resultVal, resultFinal := s.localWorker.DoWork(s.configPtr, block)
			if resultFinal {
				//s.DisposeRecvMessage(&models.MessageInfo{
				//	Item:       resultVal,
				//	CreateTime: time.Now().Local().Format(fu.CSTLayoutWithMillisecond),
				//	Workspace:  s.configPtr.Name,
				//})

				continue
			}

			// 这里强制切换Role为User，否则会导致deepseek解析报错
			resultVal.Role = common.User
			s.ChatMessage(context.Background(), s.currentWorkStage, &resultVal)
		}
	})
}
