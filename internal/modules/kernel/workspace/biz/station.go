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
	iipc "github.com/muidea/magicCodeAgent/internal/initiators/prompts/pkg/common"
	"github.com/muidea/magicCodeAgent/internal/modules/base/biz"
	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
	imklc "github.com/muidea/magicCodeAgent/internal/modules/kernel/llm/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/cache"
	"github.com/muidea/magicCodeAgent/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
	"github.com/muidea/magicCodeAgent/pkg/util"
	"github.com/muidea/magicCommon/foundation/log"
	fp "github.com/muidea/magicCommon/foundation/path"
	fu "github.com/muidea/magicCommon/foundation/util"
)

const (
	Idle      = "idle"
	Terminate = "terminate"
	Chatting  = "chatting"
	Finished  = "finished"
)

type Observer interface {
	ID() string
	OnMessage(msgPtr *models.MessageInfo)
}

type Station struct {
	bizBasePtr *biz.Base
	configPtr  *config.Workspace
	fileInfo   util.FileManager

	currentWorkStage string

	chatRunningStatus string

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

func (s *Station) envPrompt(ptr *models.Flow) string {
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
	if s.configPtr.GetStage() != common.DesignStage {
		structContent := s.promptsRegistry.ReadSectionPromptContent(fmt.Sprintf("project-%s-struct", s.getWorkspaceType()))
		if structContent != "" {
			stringBuilder.WriteString(structContent)
			stringBuilder.WriteString("\n\n")
		}
	}

	//## 项目状态
	stageContent := s.promptsRegistry.ReadSectionPromptContent(fmt.Sprintf("stage-%s", s.configPtr.GetStage()))
	if stageContent != "" {
		stringBuilder.WriteString(stageContent)
		stringBuilder.WriteString("\n\n")
	}

	// 前一阶段输出物清单
	previousLabel, previousPath := s.configPtr.GetPreviousStagePath()
	if previousPath != "" {
		stringBuilder.WriteString(fmt.Sprintf("### %s\n", previousLabel))
		s.fileInfo.Range(func(filePath string) bool {
			if strings.HasPrefix(filePath, previousPath) {
				relativePath := filePath[len(previousPath)+1:]
				stringBuilder.WriteString("- " + relativePath + "\n")
			}
			return true
		})
		stringBuilder.WriteString("\n\n")
	}

	stagePath := s.configPtr.GetCurrentStagePath()
	// 当前任务执行状态
	// 新增，删除，修改的文件清单
	stringBuilder.WriteString("### 当前任务执行状态\n")
	if ptr.AddFiles.Len() > 0 {
		stringBuilder.WriteString("#### 新增文件清单\n")
		ptr.AddFiles.Range(func(filePath string) bool {
			if strings.HasPrefix(filePath, stagePath) {
				relativePath := filePath[len(stagePath)+1:]
				stringBuilder.WriteString("- " + relativePath + "\n")
			}
			return true
		})
		stringBuilder.WriteString("\n")
	}
	if ptr.DelFiles.Len() > 0 {
		stringBuilder.WriteString("#### 删除文件清单\n")
		ptr.DelFiles.Range(func(filePath string) bool {
			if strings.HasPrefix(filePath, stagePath) {
				relativePath := filePath[len(stagePath)+1:]
				stringBuilder.WriteString("- " + relativePath + "\n")
			}
			return true
		})
		stringBuilder.WriteString("\n")
	}
	if ptr.ModFiles.Len() > 0 {
		stringBuilder.WriteString("#### 修改文件清单\n")
		ptr.ModFiles.Range(func(filePath string) bool {
			if strings.HasPrefix(filePath, stagePath) {
				relativePath := filePath[len(stagePath)+1:]
				stringBuilder.WriteString("- " + relativePath + "\n")
			}
			return true
		})
		stringBuilder.WriteString("\n")
	}
	stringBuilder.WriteString("\n\n")

	stringBuilder.WriteString("### 当前文件清单\n")
	s.fileInfo.Range(func(filePath string) bool {
		if strings.HasPrefix(filePath, stagePath) {
			relativePath := filePath[len(stagePath)+1:]
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

	curFlowPtr := s.localFlowManger.GetCurrentFlow()
	if curFlowPtr != nil {
		curFlowPtr.AddFile(filePath)
	}

	s.fileInfo.Store(filePath)
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

		curFlowPtr := s.localFlowManger.GetCurrentFlow()
		if curFlowPtr != nil {
			curFlowPtr.AddFile(path)
		}

		s.fileInfo.Store(path)
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

	curFlowPtr := s.localFlowManger.GetCurrentFlow()
	if curFlowPtr != nil {
		curFlowPtr.AddFile(filePath)
	}

	s.fileInfo.Delete(filePath)
	return true
}

func (s *Station) WriteFile(filePath string) bool {
	if !strings.HasPrefix(filePath, s.configPtr.Path) {
		return false
	}

	if s.checkIgnore(filePath) {
		return true
	}

	curFlowPtr := s.localFlowManger.GetCurrentFlow()
	if curFlowPtr != nil {
		curFlowPtr.ModFile(filePath)
	}
	return true
}

func (s *Station) ClearDir(dirPath string) bool {
	if !strings.HasPrefix(dirPath, s.configPtr.Path) {
		return false
	}

	if s.checkIgnore(dirPath) {
		return true
	}

	s.fileInfo.Range(func(filePath string) bool {
		if strings.HasPrefix(filePath, dirPath) {

			curFlowPtr := s.localFlowManger.GetCurrentFlow()
			if curFlowPtr != nil {
				curFlowPtr.AddFile(filePath)
			}

			s.fileInfo.Delete(filePath)
		}
		return true
	})

	return true
}

func (s *Station) Filter(filter fu.Filter) (ret []string) {
	s.fileInfo.Range(func(filePath string) bool {
		if filter != nil && filter.Filter(filePath) {
			ret = append(ret, filePath)
		} else if filter == nil {
			ret = append(ret, filePath)
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

func (s *Station) ChatMessage(chatContext context.Context, stage string, itemPtr *models.MessageItem) (err *cd.Error) {
	var curFlow *models.Flow
	flowIDVal := chatContext.Value(common.FlowID{})
	if flowIDVal != nil {
		curFlow = s.localFlowManger.GetFlow(flowIDVal.(string))
	} else {
		curFlow = s.localFlowManger.GetCurrentFlow()
	}
	if curFlow == nil {
		curFlow = s.localFlowManger.NewFlow(itemPtr, nil)
	} else {
		curFlow.AppendMessage(itemPtr)
	}

	s.currentWorkStage = stage
	s.chatRunningStatus = Chatting

	s.DisposeMessage(itemPtr)

	err = s.internalChat(chatContext, curFlow)
	return
}

func (s *Station) ChatTerminate(ctx context.Context, stage string) (err *cd.Error) {
	s.chatRunningStatus = Terminate

	return
}

func (s *Station) internalChat(chatContext context.Context, chatFlow *models.Flow) (err *cd.Error) {
	stationSystemPrompt := s.systemPrompt()
	stationEnvPrompt := s.envPrompt(chatFlow)
	contentContext := context.WithValue(chatContext, common.ContentSink{}, s)
	providerContext := context.WithValue(contentContext, imklc.DeepseekProvider{}, imklc.RemoteDeepseek)
	msgErr := s.localAssistant.Chat(providerContext, s.configPtr, stationSystemPrompt, stationEnvPrompt, chatFlow.GetFlowPath())
	if msgErr != nil {
		err = msgErr
		log.Errorf("internalChat, s.localAssistant.Chat error:%v", msgErr)
		return
	}

	return
}

func (s *Station) DisposeMessage(itemPtr *models.MessageItem) {
	messageInfoPtr := &models.MessageInfo{
		Item:       *itemPtr,
		CreateTime: time.Now().Local().Format(fu.CSTLayoutWithMillisecond),
		Workspace:  s.configPtr.Name,
	}

	msgPtr := s.localMessageCache.Add(messageInfoPtr)

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
	msgPtr := &models.MessageItem{
		Role:    common.Assistant,
		Content: []string{"通讯中断，需要重试！"},
	}

	s.DisposeMessage(msgPtr)
}

func (s *Station) OnRecv(event string, content string) {
	msgPtr := &models.MessageItem{
		Role:    common.Assistant,
		Content: []string{content},
	}

	curFlow := s.localFlowManger.GetCurrentFlow()
	if curFlow != nil {
		curFlow.AppendMessage(msgPtr)
	}

	s.DisposeMessage(msgPtr)

	// 如果处于终止状态，则不处理
	if s.chatRunningStatus == Terminate {
		return
	}

	s.bizBasePtr.AsyncTask(func() {
		blockList, blockErr := util.SpliteBlocks(content)
		if blockErr != nil {
			s.chatRunningStatus = Terminate
			log.Errorf("OnRecv, util.SpliteBlocks err:%v", blockErr)
			return
		}

		curFlow := s.localFlowManger.GetCurrentFlow()
		for _, block := range blockList {
			workParam, workResult, workErr := s.localWorker.DoWork(s.configPtr, block)
			if workErr != nil {
				s.chatRunningStatus = Terminate
				log.Errorf("OnRecv, s.localWorker.DoWork err:%v", workErr)
				return
			}

			//workResult.Result.Role = common.User
			switch workResult.Name {
			case "thinking":
				// noting todo
			case "ask_followup_question":
				curFlow = s.tryToContinue(curFlow, workParam, workResult)
			case "attempt_completion":
				curFlow = s.tryToFinish(curFlow, workParam, workResult)
			case "decompose_design_subtask":
				curFlow = s.tryToSwitch(curFlow, workParam, workResult)
			default:
				curFlow = s.continueRun(curFlow, workParam, workResult)
			}

			if curFlow != nil {
				s.DisposeMessage(&workResult.Result)
				err := s.internalChat(context.Background(), curFlow)
				if err != nil {
					s.chatRunningStatus = Terminate
					log.Errorf("internalChat error: %v", err)
					return
				}
			}
		}
	})
}

func (s *Station) tryToContinue(_ *models.Flow, _ imbtc.ToolParam, _ *imbtc.ToolResult) (ret *models.Flow) {
	// ask_followup_question
	// 当前阶段是需要等待用户的输入的，不用做任何事情
	return
}

func (s *Station) tryToFinish(curFlow *models.Flow, _ imbtc.ToolParam, _ *imbtc.ToolResult) (ret *models.Flow) {
	// attempt_completion
	// 检查所有的flow是否都执行完毕，如果执行完毕就
	ret = s.localFlowManger.FinishFlow(curFlow)
	return
}

func (s *Station) tryToSwitch(curFlow *models.Flow, param imbtc.ToolParam, result *imbtc.ToolResult) (ret *models.Flow) {
	// decompose_design_subtask
	// 到这里说明在开始任务分解，需要根据任务分解的层次切换当前的Flow
	// 同时这里进行校验，避免大模型跑偏了，
	level := param["level"].(string)
	msgItem := result.Result
	msgItem.Role = common.User
	switch common.CompareLevel(curFlow.GetLevel(), level) {
	case common.LevelGreater:
		s.localFlowManger.NewFlow(&msgItem, curFlow)
		// 子级任务
		//s.localFlowManger.NewFlow(chatMsg, curFlow)
	case common.LevelEqual:
		// 同级任务
		s.localFlowManger.NewFlow(&msgItem, curFlow.GetParent())
	case common.LevelLess:
		// 父级任务
		parent := curFlow.GetParent()
		if parent != nil {
			s.localFlowManger.NewFlow(&msgItem, parent.GetParent())
		} else {
			// 到这里属于异常了，先记录日志看看情况
			log.Warnf("curFlow.GetParent() is nil, %+v", curFlow.GetFlowPath())
			s.localFlowManger.NewFlow(&msgItem, nil)
		}
	}

	ret = s.localFlowManger.GetCurrentFlow()
	return
}

func (s *Station) continueRun(curFlow *models.Flow, _ imbtc.ToolParam, result *imbtc.ToolResult) (ret *models.Flow) {
	msgItem := result.Result
	msgItem.Role = common.User

	curFlow.AppendMessage(&msgItem)

	ret = curFlow
	return
}
