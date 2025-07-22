package biz

import (
	"bytes"
	"context"
	"fmt"
	"os"
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
	Paused    = "paused"
	Finished  = "finished"
)

const maxLoopCount = 100

type Observer interface {
	ID() string
	OnMessage(msgPtr *models.MessageInfo)
}

type Station struct {
	bizBasePtr      *biz.Base
	workspaceConfig models.Workspace
	fileInfo        util.FileManager

	chatRunningStatus string

	chatMaxLoopCount int64
	chatLoopCount    int64
	chatUsagePtr     *imklc.Usage

	stationMutex sync.RWMutex
	observerList []Observer

	localMessageCache    *cache.MessageSequenceCache
	localFlowManger      *FlowManager
	localWorker          Worker
	localAssistant       Assistant
	localPromptReader    PromptReader
	innerPromptsRegistry iipc.PromptsRegistry
}

func NewStation(
	bizBasePtr *biz.Base,
	workspacePtr *models.Workspace,
	worker Worker,
	assistant Assistant,
	reader PromptReader,
	promptsRegistry iipc.PromptsRegistry,
) *Station {
	ret := &Station{
		bizBasePtr:           bizBasePtr,
		workspaceConfig:      *workspacePtr,
		chatRunningStatus:    Idle,
		chatLoopCount:        0,
		chatUsagePtr:         &imklc.Usage{},
		localMessageCache:    cache.NewMessageSequenceCache(500),
		localFlowManger:      NewFlowManager(),
		localWorker:          worker,
		localAssistant:       assistant,
		localPromptReader:    reader,
		innerPromptsRegistry: promptsRegistry,
	}

	return ret
}

func (s *Station) systemPrompt() string {
	return s.innerPromptsRegistry.GetStagePromptsContent(s.workspaceConfig.GetStage())
}

func (s *Station) getWorkspaceProjectStruct() (string, error) {
	param := map[string]string{}
	return s.localPromptReader.ReadPrompt(&s.workspaceConfig, "project-struct", param)
}

func (s *Station) getWorkspaceStatus(ptr *models.Flow) (string, error) {
	var levelInfo strings.Builder
	switch ptr.GetLevel() {
	case models.ProjectLevel:
		levelInfo.WriteString(fmt.Sprintf("项目(`%s`)，项目名-%s", ptr.GetLevel(), ptr.GetOwner()))
	case models.ModuleLevel:
		levelInfo.WriteString(fmt.Sprintf("业务模块(`%s`)，业务模块名-%s", ptr.GetLevel(), ptr.GetOwner()))
		if parent := ptr.GetParent(); parent != nil {
			levelInfo.WriteString(fmt.Sprintf(", 所属项目-%s", parent.GetOwner()))
		}
	case models.BlockLevel:
		levelInfo.WriteString(fmt.Sprintf("功能模块(`%s`)，功能模块名-%s", ptr.GetLevel(), ptr.GetOwner()))
		if parent := ptr.GetParent(); parent != nil {
			levelInfo.WriteString(fmt.Sprintf(", 所属业务模块-%s", parent.GetOwner()))
		}
	default:
		return "", nil
	}

	param := map[string]string{
		"LevelInfo": levelInfo.String(),
	}
	return s.localPromptReader.ReadPrompt(&s.workspaceConfig, fmt.Sprintf("stage-%s-status", s.workspaceConfig.GetStage()), param)
}

func (s *Station) getWorkspaceSpecialInstruction() string {
	param := map[string]string{}
	return s.innerPromptsRegistry.ReadSectionPromptContent(fmt.Sprintf("special-%s-instruction", s.workspaceConfig.GetStage()), param)
}

func (s *Station) envPrompt(ptr *models.Flow) (string, error) {
	stringBuilder := bytes.NewBufferString("")
	stringBuilder.WriteString("<environment_details>\n")
	stringBuilder.WriteString("# 项目名称\n\n" + s.workspaceConfig.Name + "\n\n")

	stringBuilder.WriteString("## 项目描述\n\n")
	stringBuilder.WriteString(s.workspaceConfig.Description + "\n\n")

	stringBuilder.WriteString("## 技术栈\n\n")
	for _, val := range s.workspaceConfig.TechStack {
		stringBuilder.WriteString("- " + val + "\n")
	}
	stringBuilder.WriteString("\n\n")

	//## 项目规划目录
	if s.workspaceConfig.GetStage() != models.DesignStage {
		structContent, structErr := s.getWorkspaceProjectStruct()
		if structErr != nil {
			log.Errorf("get workspace project struct failed, %s", structErr)
			return "", structErr
		}

		if structContent != "" {
			stringBuilder.WriteString(structContent)
			stringBuilder.WriteString("\n\n")
		}
	}

	//## 项目状态
	statusContent, statusErr := s.getWorkspaceStatus(ptr)
	if statusErr != nil {
		log.Errorf("get workspace status failed, flow:%v error:%s", ptr, statusErr)
		return "", statusErr
	}
	if statusContent != "" {
		stringBuilder.WriteString(statusContent)
		stringBuilder.WriteString("\n\n")
	}

	// 参考内容清单
	referenceStageVal := models.GetStageReferenceStage(s.workspaceConfig.GetStage())
	referenceLabel := models.GetStageReferenceLabel(referenceStageVal)
	referencePath := s.workspaceConfig.GetStageWorkingPath(referenceStageVal)
	if referencePath != "" {
		hasReference := false
		stringBuilder.WriteString(fmt.Sprintf("### %s\n\n", referenceLabel))
		s.fileInfo.Range(func(filePath string) bool {
			if strings.HasPrefix(filePath, referencePath) {
				relativePath := filePath[len(referencePath)+1:]
				stringBuilder.WriteString("- " + relativePath + "\n")
				hasReference = true
			}
			return true
		})
		if hasReference {
			stringBuilder.WriteString(fmt.Sprintf("\n* 可以使用`read_stage_file`工具, 指定stage为: %s 读取对应的内容\n", referenceStageVal))
			stringBuilder.WriteString("\n\n")
		}
	}

	workingPath := s.workspaceConfig.GetWorkingPath()
	/*
		if ptr != nil {
			// 当前进度
			// 新增，删除，修改的文件清单
			stringBuilder.WriteString("### 当前进度\n\n")
			if ptr.AddFiles.Len() > 0 {
				stringBuilder.WriteString("#### 新增清单\n\n")
				ptr.AddFiles.Range(func(filePath string) bool {
					if strings.HasPrefix(filePath, workingPath) {
						relativePath := filePath[len(workingPath)+1:]
						stringBuilder.WriteString("- " + relativePath + "\n")
					}
					return true
				})
				stringBuilder.WriteString("\n")
			}
			if ptr.DelFiles.Len() > 0 {
				stringBuilder.WriteString("#### 删除清单\n\n")
				ptr.DelFiles.Range(func(filePath string) bool {
					if strings.HasPrefix(filePath, workingPath) {
						relativePath := filePath[len(workingPath)+1:]
						stringBuilder.WriteString("- " + relativePath + "\n")
					}
					return true
				})
				stringBuilder.WriteString("\n")
			}
			if ptr.ModFiles.Len() > 0 {
				stringBuilder.WriteString("#### 修改清单\n\n")
				ptr.ModFiles.Range(func(filePath string) bool {
					if strings.HasPrefix(filePath, workingPath) {
						relativePath := filePath[len(workingPath)+1:]
						stringBuilder.WriteString("- " + relativePath + "\n")
					}
					return true
				})
				stringBuilder.WriteString("\n")
			}
			stringBuilder.WriteString("\n\n")
		}
	*/

	stringBuilder.WriteString("### 当前输出清单\n\n")
	s.fileInfo.Range(func(filePath string) bool {
		if strings.HasPrefix(filePath, workingPath) {
			relativePath := filePath[len(workingPath)+1:]
			stringBuilder.WriteString("- " + relativePath + "\n")
		}
		return true
	})
	stringBuilder.WriteString("\n\n")

	instructContent := s.getWorkspaceSpecialInstruction()
	if instructContent != "" {
		stringBuilder.WriteString(instructContent)
		stringBuilder.WriteString("\n\n")
	}

	stringBuilder.WriteString("</environment_details>")
	return stringBuilder.String(), nil
}

func (s *Station) Available() bool {
	return s.workspaceConfig.Available()
}

func (s *Station) checkIgnore(fileFullPath string) bool {
	for _, ignore := range s.workspaceConfig.Ignores {
		if strings.Contains(fileFullPath, ignore) {
			return true
		}

		continue
	}
	return false
}

func (s *Station) AddFile(filePath string) bool {
	if !fp.IsSubPath(s.workspaceConfig.GetRootPath(), filePath) {
		return false
	}

	if s.checkIgnore(filePath) {
		return true
	}

	//log.Infof("AddFile filePath:%s", filePath)
	curFlowPtr := s.localFlowManger.GetCurrentFlow()
	if curFlowPtr != nil {
		curFlowPtr.AddFile(filePath)
	}

	s.fileInfo.Store(filePath)
	return true
}

func (s *Station) RemoveFile(filePath string) bool {
	if !fp.IsSubPath(s.workspaceConfig.WorkspaceDir, filePath) {
		return false
	}

	if s.checkIgnore(filePath) {
		return true
	}

	//log.Infof("RemoveFile filePath:%s", filePath)
	curFlowPtr := s.localFlowManger.GetCurrentFlow()
	if curFlowPtr != nil {
		curFlowPtr.AddFile(filePath)
	}

	s.fileInfo.Delete(filePath)
	return true
}

func (s *Station) WriteFile(filePath string) bool {
	if !fp.IsSubPath(s.workspaceConfig.WorkspaceDir, filePath) {
		return false
	}

	if s.checkIgnore(filePath) {
		return true
	}

	//log.Infof("WriteFile filePath:%s", filePath)
	curFlowPtr := s.localFlowManger.GetCurrentFlow()
	if curFlowPtr != nil {
		curFlowPtr.ModFile(filePath)
	}
	s.fileInfo.Store(filePath)
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

func (s *Station) ChatMessage(chatContext context.Context, itemPtr *models.MessageItem) (err *cd.Error) {
	if s.chatRunningStatus == Terminate || s.chatRunningStatus == Finished {
		// 表示从结束状态恢复过来，目前是重新开始，这里要主动的重置
		s.localFlowManger.Reset()
	}

	var curFlow *models.Flow
	flowIDVal := chatContext.Value(common.FlowID{})
	if flowIDVal != nil {
		curFlow = s.localFlowManger.GetFlow(flowIDVal.(string))
	} else {
		curFlow = s.localFlowManger.GetCurrentFlow()
	}
	if curFlow == nil {
		curFlow = s.localFlowManger.NewFlow(itemPtr.String(), "workspace", nil)
	}
	curFlow.AppendMessage(itemPtr)

	s.chatMaxLoopCount = maxLoopCount
	s.chatRunningStatus = Chatting
	s.chatLoopCount = 0

	s.DisposeMessage(itemPtr)

	err = s.internalChat(chatContext, curFlow)
	return
}

func (s *Station) ChatTerminate(ctx context.Context) (err *cd.Error) {
	s.chatRunningStatus = Terminate

	//s.localFlowManger.Reset()
	return
}

func (s *Station) internalChat(chatContext context.Context, chatFlow *models.Flow) (err *cd.Error) {
	// 如果处于终止状态,或者已经超过最大循环次数，则不处理
	if s.chatRunningStatus != Chatting {
		return
	}

	s.chatLoopCount++

	sysPrompt := s.systemPrompt()
	envPrompt, envErr := s.envPrompt(chatFlow)
	if envErr != nil {
		err = cd.NewError(cd.Unexpected, "get env prompt failed")
		return
	}

	contentContext := context.WithValue(chatContext, common.ContentSink{}, s)
	providerContext := context.WithValue(contentContext, imklc.DeepseekProvider{}, imklc.RemoteDeepseek)
	msgErr := s.localAssistant.Chat(providerContext, &s.workspaceConfig, sysPrompt, envPrompt, chatFlow.GetFlowPath())
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
		Workspace:  s.workspaceConfig.Name,
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
	s.chatUsagePtr.Update(usagePtr)

	s.bizBasePtr.AsyncTask(func() {
		s.DisposeContent(content)
	})
}

func (s *Station) DisposeContent(content string) {
	curTime := time.Now()
	var finalStatus = s.chatRunningStatus
	var chatLoopCount = s.chatLoopCount
	var notifyMsg = "会话被中断！"
	log.Infof("DisposeContent, workspace:%s, chatLoopCount:%d", s.workspaceConfig.Name, chatLoopCount)
	defer func() {
		log.Infof("DisposeContent done, workspace:%s, chatLoopCount:%d, finalStatus:%s, elapseTime:%v", s.workspaceConfig.Name, chatLoopCount, finalStatus, time.Since(curTime))
		if finalStatus != Chatting {
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

	msgPtr := &models.MessageItem{
		Role:    common.Assistant,
		Content: []string{content},
	}

	curFlow := s.localFlowManger.GetCurrentFlow()
	s.DisposeMessage(msgPtr)

	blockList, blockErr := util.SpliteBlocks(content)
	if blockErr != nil {
		finalStatus = Terminate
		notifyMsg = fmt.Sprintf("执行出错, 分解数据出错:%v", blockErr.Error())
		log.Errorf("DisposeContent, util.SpliteBlocks err:%v", blockErr)
		return
	}

	for _, block := range blockList {
		workParam, workResult, workErr := s.localWorker.DoWork(&s.workspaceConfig, block)
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
		case "attempt_finish":
			curFlow.AppendMessage(msgPtr)
			curFlow = s.tryToFinishFlow(curFlow, workParam, msgItem)
		case "decompose_design_plan", "decompose_code_plan":
			curFlow = s.tryToSwitchFlow(curFlow, workParam, msgItem)
			curFlow.AppendMessage(msgPtr)
			curFlow.AppendMessage(msgItem)
		case "attempt_completion":
			curFlow.AppendMessage(msgPtr)
			curFlow = s.tryToCompleteFlow(curFlow, workParam, msgItem)
		default:
			curFlow.AppendMessage(msgPtr)
			curFlow = s.continueRun(curFlow, workParam, msgItem)
		}

		s.DisposeMessage(msgItem)
	}

	if chatLoopCount > s.chatMaxLoopCount {
		finalStatus = Paused
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

func (s *Station) tryToFinishFlow(curFlow *models.Flow, _ imbtc.ToolParam, _ *models.MessageItem) (ret *models.Flow) {
	// attempt_finish
	// 检查所有的flow是否都执行完毕，如果执行完毕就
	log.Infof("try to finish flow %v", curFlow)
	ret = s.localFlowManger.FinishFlow(curFlow)
	return
}

func (s *Station) tryToSwitchFlow(curFlow *models.Flow, param imbtc.ToolParam, msgItem *models.MessageItem) (ret *models.Flow) {
	// decompose_design_plan
	// 到这里说明在开始任务分解，需要根据任务分解的层次切换当前的Flow
	// 同时这里进行校验，避免大模型跑偏了，
	log.Infof("try to switch flow %v", curFlow)
	level := param["level"].(string)
	name := param["name"].(string)
	compareVal := models.LevelGreater
	if curFlow != nil {
		compareVal = models.CompareLevel(curFlow.GetLevel(), level)
	}
	switch compareVal {
	case models.LevelGreater:
		ret = s.localFlowManger.NewFlow(msgItem.String(), name, curFlow)
		// 子级任务
		//s.localFlowManger.NewFlow(chatMsg, curFlow)
	case models.LevelEqual:
		// 同级任务
		if curFlow.CanFinish() {
			curFlow.Status = models.FLOW_STATUS_FINISH
		}
		ret = s.localFlowManger.NewFlow(msgItem.String(), name, curFlow.GetParent())
	case models.LevelLess:
		// 父级任务
		parent := curFlow.GetParent()
		if parent != nil {
			ret = s.localFlowManger.NewFlow(msgItem.String(), name, parent.GetParent())
		} else {
			// 到这里属于异常了，先记录日志看看情况
			log.Warnf("curFlow.GetParent() is nil, %+v", curFlow.GetFlowPath())
			ret = s.localFlowManger.NewFlow(msgItem.String(), name, nil)
		}
	}

	//ret.AppendMessage(msgItem)
	return
}

func (s *Station) tryToCompleteFlow(curFlow *models.Flow, _ imbtc.ToolParam, _ *models.MessageItem) (ret *models.Flow) {
	// attempt_completion
	// 到这里说明当前的Flow已经完成了，直接结束当前的Flow
	log.Infof("try to complete flow %v", curFlow)
	ret = s.localFlowManger.CompleteFlow(curFlow)
	return
}

func (s *Station) continueRun(curFlow *models.Flow, _ imbtc.ToolParam, msgItem *models.MessageItem) (ret *models.Flow) {
	log.Infof("continue run flow %v", curFlow)
	if curFlow != nil {
		curFlow.AppendMessage(msgItem)
	}

	ret = curFlow
	return
}

func (s *Station) GetFlows() []*models.Flow {
	return s.localFlowManger.GetLocalFlows()
}

func (s *Station) DebugContent(ctx context.Context) (err *cd.Error) {
	debugPath := filepath.Join(config.CurrentWorkingPath(), "debug")
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
		debugFlow = s.localFlowManger.NewFlow("debug", "debug", nil)
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
	msgErr := s.localAssistant.Chat(providerContext, &s.workspaceConfig, systemPrompt, envPrompt, debugFlow.GetFlowPath())
	if msgErr != nil {
		log.Errorf("chat error:%v", msgErr)
		return
	}

	return
}
