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
	"github.com/muidea/magicCommon/foundation/log"
	fp "github.com/muidea/magicCommon/foundation/path"
	fu "github.com/muidea/magicCommon/foundation/util"

	"github.com/muidea/magicCodeAgent/internal/config"
	iipc "github.com/muidea/magicCodeAgent/internal/initiators/prompts/pkg/common"
	"github.com/muidea/magicCodeAgent/internal/modules/base/biz"
	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
	imklc "github.com/muidea/magicCodeAgent/internal/modules/kernel/llm/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/cache"
	"github.com/muidea/magicCodeAgent/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
	"github.com/muidea/magicCodeAgent/pkg/util"
)

const (
	Idle      = "idle"
	Terminate = "terminate"
	Chatting  = "chatting"
	Finished  = "finished"
)

const maxLoopCount = 100

type Observer interface {
	ID() string
	OnMessage(msgPtr *models.MessageInfo)
}

type Station struct {
	bizBasePtr *biz.Base
	configPtr  *config.Workspace
	fileInfo   util.FileManager

	chatRunningStatus string

	chatMaxLoopCount int64
	chatLoopCount    int64
	chatUsagePtr     *imklc.Usage

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
		chatLoopCount:     0,
		chatUsagePtr:      &imklc.Usage{},
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
	filterTags := s.configPtr.TechStack
	filterTags = filterTags.Add(s.configPtr.GetWorkspaceType())

	return s.promptsRegistry.GetPromptsContent(s.configPtr.GetStage(), filterTags)
}

func (s *Station) envPrompt(ptr *models.Flow) string {
	stringBuilder := bytes.NewBufferString("")
	stringBuilder.WriteString("<environment_details>\n")
	stringBuilder.WriteString("# 项目名称\n\n" + s.configPtr.Name + "\n\n")

	stringBuilder.WriteString("## 项目描述\n\n")
	stringBuilder.WriteString(s.configPtr.Description + "\n\n")

	stringBuilder.WriteString("## 技术栈\n\n")
	for _, val := range s.configPtr.TechStack {
		stringBuilder.WriteString("- " + val + "\n")
	}
	stringBuilder.WriteString("\n\n")

	param := map[string]string{}
	//## 项目规划目录
	if s.configPtr.GetStage() != common.DesignStage {
		structContent := s.promptsRegistry.ReadSectionPromptContent(fmt.Sprintf("project-struct-%s", s.getWorkspaceType()), param)
		if structContent != "" {
			stringBuilder.WriteString(structContent)
			stringBuilder.WriteString("\n\n")
		}
	}

	//## 项目状态
	param = map[string]string{
		"Level": ptr.GetLevel(),
	}
	statusContent := s.promptsRegistry.ReadSectionPromptContent(fmt.Sprintf("stage-%s-status-%s", s.configPtr.GetStage(), s.getWorkspaceType()), param)
	if statusContent != "" {
		stringBuilder.WriteString(statusContent)
		stringBuilder.WriteString("\n\n")
	}

	// 参考内容清单
	referenceStageVal := ""
	switch s.configPtr.GetStage() {
	case common.CodeStage, common.TestStage:
		referenceStageVal = common.DesignStage
	}
	referenceStage, referenceLabel, referencePath := s.configPtr.GetReferenceStagePath(referenceStageVal)
	if referencePath != "" {
		stringBuilder.WriteString(fmt.Sprintf("### %s\n\n", referenceLabel))
		s.fileInfo.Range(func(filePath string) bool {
			if strings.HasPrefix(filePath, referencePath) {
				relativePath := filePath[len(referencePath)+1:]
				stringBuilder.WriteString("- " + relativePath + "\n")
			}
			return true
		})
		stringBuilder.WriteString(fmt.Sprintf("\n* 可以使用`read_stage_file`工具, 指定stage为: %s 读取对应的内容\n", referenceStage))
		stringBuilder.WriteString("\n\n")
	}

	stagePath := s.configPtr.GetStageWorkPath()
	if ptr != nil {
		// 当前任务执行状态
		// 新增，删除，修改的文件清单
		stringBuilder.WriteString("### 当前任务执行状态\n\n")
		if ptr.AddFiles.Len() > 0 {
			stringBuilder.WriteString("#### 新增文件清单\n\n")
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
			stringBuilder.WriteString("#### 删除文件清单\n\n")
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
			stringBuilder.WriteString("#### 修改文件清单\n\n")
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
	}

	stageStuffix := ""
	if s.configPtr.GetStage() == common.DesignStage {
		stageStuffix = common.DesignStage
	} else {
		stageStuffix = common.CodeStage
	}
	stringBuilder.WriteString("### 当前文件清单\n\n")
	s.fileInfo.Range(func(filePath string) bool {
		if strings.HasPrefix(filePath, stagePath) {
			relativePath := filePath[len(stagePath)+1:]
			stringBuilder.WriteString("- " + relativePath + "\n")
		}
		return true
	})
	stringBuilder.WriteString("\n\n")

	param = map[string]string{}
	instructContent := s.promptsRegistry.ReadSectionPromptContent(fmt.Sprintf("special-%s-instruction", stageStuffix), param)
	if instructContent != "" {
		stringBuilder.WriteString(instructContent)
		stringBuilder.WriteString("\n\n")
	}

	stringBuilder.WriteString("</environment_details>")
	return stringBuilder.String()
}

func (s *Station) Available() bool {
	return s.configPtr.Available()
}

func (s *Station) checkIgnore(fileFullPath string) bool {
	if isIgnored(fileFullPath) {
		return true
	}

	stagePath := s.configPtr.GetStageWorkPath()
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

func (s *Station) BroadcastReady() {
	if s.localMessageCache.Size() > 0 {
		return
	}

	notifyMsg := &models.MessageItem{
		Role: common.Notifier,
		Content: []string{
			"会话准备完成！",
		},
	}

	s.DisposeMessage(notifyMsg)
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
		curFlow = s.localFlowManger.NewFlow(itemPtr.String(), nil)
	}
	curFlow.AppendMessage(itemPtr)

	s.chatMaxLoopCount = maxLoopCount
	s.chatRunningStatus = Chatting
	s.chatLoopCount = 0

	s.DisposeMessage(itemPtr)

	err = s.internalChat(chatContext, curFlow)
	return
}

func (s *Station) ChatTerminate(ctx context.Context, stage string) (err *cd.Error) {
	s.chatRunningStatus = Terminate

	s.localFlowManger.Reset()
	return
}

func (s *Station) internalChat(chatContext context.Context, chatFlow *models.Flow) (err *cd.Error) {
	s.chatLoopCount++

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

func (s *Station) OnRecv(event, content string, usagePtr *imklc.Usage) {
	msgPtr := &models.MessageItem{
		Role:    common.Assistant,
		Content: []string{content},
	}

	s.chatUsagePtr.Update(usagePtr)

	// 如果处于终止状态,或者已经超过最大循环次数，则不处理
	if s.chatRunningStatus != Chatting {
		return
	}

	curFlow := s.localFlowManger.GetCurrentFlow()
	if curFlow != nil {
		curFlow.AppendMessage(msgPtr)
	}

	s.DisposeMessage(msgPtr)

	s.bizBasePtr.AsyncTask(func() {
		s.DisposeContent(curFlow, content)
	})
}

func (s *Station) DisposeContent(curFlow *models.Flow, content string) {
	curTime := time.Now()
	var finalStatus = s.chatRunningStatus
	var chatLoopCount = s.chatLoopCount
	var notifyMsg = "会话被中断！"
	log.Infof("DisposeContent, workspace:%s, chatLoopCount:%d", s.configPtr.Name, chatLoopCount)
	defer func() {
		log.Infof("DisposeContent done, workspace:%s, chatLoopCount:%d, finalStatus:%s, elapseTime:%v", s.configPtr.Name, chatLoopCount, finalStatus, time.Since(curTime))
		if finalStatus == Terminate || finalStatus == Finished {
			// 发送本次消耗的Token数据
			notifyMsg := &models.MessageItem{
				Role: common.Notifier,
				Content: []string{
					formatUsage(s.chatUsagePtr, chatLoopCount),
					notifyMsg,
				},
			}
			s.DisposeMessage(notifyMsg)

			s.chatRunningStatus = finalStatus
		}
	}()

	blockList, blockErr := util.SpliteBlocks(content)
	if blockErr != nil {
		finalStatus = Terminate
		notifyMsg = fmt.Sprintf("执行出错, 分解数据出错:%v", blockErr.Error())
		log.Errorf("DisposeContent, util.SpliteBlocks err:%v", blockErr)
		return
	}

	for _, block := range blockList {
		workParam, workResult, workErr := s.localWorker.DoWork(s.configPtr, block)
		if workErr != nil {
			finalStatus = Terminate
			notifyMsg = fmt.Sprintf("执行出错, Tool执行出错:%v", workErr.Error())
			log.Errorf("DisposeContent, s.localWorker.DoWork err:%v", workErr)
			return
		}

		msgItem := &workResult.Result
		msgItem.Role = common.User
		//workResult.Result.Role = common.User
		switch workResult.Name {
		case "thinking":
			// noting todo
			continue
		case "ask_followup_question":
			curFlow = s.waitToContinue(curFlow, workParam, msgItem)
		case "attempt_completion":
			curFlow = s.tryToFinishFlow(curFlow, workParam, msgItem)
		case "decompose_design_plan", "decompose_code_plan":
			curFlow = s.tryToSwitchFlow(curFlow, workParam, msgItem)
		default:
			curFlow = s.continueRun(curFlow, workParam, msgItem)
		}

		s.DisposeMessage(msgItem)
	}

	if chatLoopCount > s.chatMaxLoopCount {
		finalStatus = Finished
		notifyMsg = "循环次数超过最大限制，自动结束\n 如果要继续执行，可以输入'继续执行'"
		return
	}

	if curFlow != nil {
		err := s.internalChat(context.Background(), curFlow)
		if err != nil {
			finalStatus = Terminate
			notifyMsg = fmt.Sprintf("执行出错, 调用大模型出错:%v", err.Error())
			log.Errorf("DisposeContent, internalChat error: %v", err)
			return
		}
	}

	// 如果所有的Flow都执行完毕，则将状态置为Finish
	if curFlow == nil {
		log.Infof("all flow finished")
		finalStatus = Finished
		notifyMsg = "执行完毕"
	}
}

func (s *Station) waitToContinue(curFlow *models.Flow, _ imbtc.ToolParam, _ *models.MessageItem) (ret *models.Flow) {
	// ask_followup_question
	// 当前阶段是需要等待用户的输入的，不用做任何事情
	log.Infof("waitting to continue flow %s", curFlow)
	return
}

func (s *Station) tryToFinishFlow(curFlow *models.Flow, _ imbtc.ToolParam, _ *models.MessageItem) (ret *models.Flow) {
	// attempt_completion
	// 检查所有的flow是否都执行完毕，如果执行完毕就
	log.Infof("try to finish flow %s", curFlow)
	ret = s.localFlowManger.FinishFlow(curFlow)
	return
}

func (s *Station) tryToSwitchFlow(curFlow *models.Flow, param imbtc.ToolParam, msgItem *models.MessageItem) (ret *models.Flow) {
	// decompose_design_plan
	// 到这里说明在开始任务分解，需要根据任务分解的层次切换当前的Flow
	// 同时这里进行校验，避免大模型跑偏了，
	log.Infof("try to switch flow %s", curFlow)
	level := param["level"].(string)
	compareVal := models.LevelGreater
	if curFlow != nil {
		compareVal = models.CompareLevel(curFlow.GetLevel(), level)
	}
	switch compareVal {
	case models.LevelGreater:
		ret = s.localFlowManger.NewFlow(msgItem.String(), curFlow)
		// 子级任务
		//s.localFlowManger.NewFlow(chatMsg, curFlow)
	case models.LevelEqual:
		// 同级任务
		if curFlow.CanFinish() {
			curFlow.Status = models.FLOW_STATUS_FINISH
		}
		ret = s.localFlowManger.NewFlow(msgItem.String(), curFlow.GetParent())
	case models.LevelLess:
		// 父级任务
		parent := curFlow.GetParent()
		if parent != nil {
			ret = s.localFlowManger.NewFlow(msgItem.String(), parent.GetParent())
		} else {
			// 到这里属于异常了，先记录日志看看情况
			log.Warnf("curFlow.GetParent() is nil, %+v", curFlow.GetFlowPath())
			ret = s.localFlowManger.NewFlow(msgItem.String(), nil)
		}
	}

	ret.AppendMessage(msgItem)
	return
}

func (s *Station) continueRun(curFlow *models.Flow, _ imbtc.ToolParam, msgItem *models.MessageItem) (ret *models.Flow) {
	log.Infof("continue run flow %s", curFlow)
	curFlow.AppendMessage(msgItem)

	ret = curFlow
	return
}

func (s *Station) GetFlows() []*models.Flow {
	return s.localFlowManger.GetLocalFlows()
}

func (s *Station) DebugContent(ctx context.Context, stage string) (err *cd.Error) {
	debugPath := filepath.Join(config.CurrentWorkPath(), "debug")
	fileList, fileErr := fp.ListPath(debugPath, "*.md", false)
	if fileErr != nil {
		log.Errorf("debug file error:%v", fileErr)
		return
	}
	if len(fileList) < 3 {
		log.Infof("not enough debug file")
		return
	}

	var debugFlow *models.Flow
	flowIDVal := ctx.Value(common.FlowID{})
	if flowIDVal != nil {
		debugFlow = s.localFlowManger.GetFlow(flowIDVal.(string))
	} else {
		debugFlow = s.localFlowManger.GetCurrentFlow()
	}
	if debugFlow == nil {
		debugFlow = s.localFlowManger.NewFlow("debug", nil)
	}

	systemPromptFile := fileList[0]
	systemByte, systemErr := os.ReadFile(systemPromptFile)
	if systemErr != nil {
		log.Errorf("read system prompt file error:%v", systemErr)
		return
	}
	userPromptFiles := fileList[1 : len(fileList)-1]
	envPromptFile := fileList[len(fileList)-1]
	envByte, envErr := os.ReadFile(envPromptFile)
	if envErr != nil {
		log.Errorf("read env prompt file error:%v", envErr)
		return
	}
	log.Infof("debug system prompt: %v", systemPromptFile)
	for _, val := range userPromptFiles {
		userByte, userErr := os.ReadFile(val)
		if userErr != nil {
			log.Errorf("read user prompt file error:%v", userErr)
			return
		}
		debugFlow.AppendMessage(&models.MessageItem{
			Role: common.User,
			Content: []string{
				string(userByte),
			},
		})
		log.Infof("debug user file: %s", val)
	}
	log.Infof("debug env prompt: %v", envPromptFile)

	s.chatMaxLoopCount = 2
	s.chatRunningStatus = Chatting

	s.chatLoopCount++
	systemPrompt := string(systemByte)
	envPrompt := string(envByte)
	contentContext := context.WithValue(ctx, common.ContentSink{}, s)
	providerContext := context.WithValue(contentContext, imklc.DeepseekProvider{}, imklc.RemoteDeepseek)
	msgErr := s.localAssistant.Chat(providerContext, s.configPtr, systemPrompt, envPrompt, debugFlow.GetFlowPath())
	if msgErr != nil {
		log.Errorf("chat error:%v", msgErr)
		return
	}

	return
}
