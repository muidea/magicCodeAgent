package biz

import (
	"fmt"

	"github.com/muidea/magicCommon/foundation/log"

	"github.com/muidea/magicCodeAgent/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"

	iipc "github.com/muidea/magicCodeAgent/internal/initiators/prompts/pkg/common"
)

type FlowManager struct {
	localFlows     []*models.Flow
	currentFlowPtr *models.Flow
}

func NewFlowManager() *FlowManager {
	return &FlowManager{
		localFlows:     make([]*models.Flow, 0),
		currentFlowPtr: nil,
	}
}

func (s *FlowManager) Reset() {
	s.localFlows = make([]*models.Flow, 0)
	s.currentFlowPtr = nil
}

func (s *FlowManager) newRootFlow(name string) *models.Flow {
	flowPtr := &models.Flow{
		Id:        fmt.Sprintf("flow_%d", len(s.localFlows)),
		Name:      name,
		Message:   []*models.MessageItem{},
		Status:    models.FLOW_STATUS_WAIT,
		Parent:    nil,
		Childrens: make([]*models.Flow, 0),
		Level:     models.RootLevel,
	}

	return s.appendFlow(flowPtr, nil)
}

func (s *FlowManager) NewFlow(name string, parentPtr *models.Flow) *models.Flow {
	if parentPtr == nil {
		parentPtr = s.currentFlowPtr
	}

	if parentPtr == nil {
		return s.newRootFlow(name)
	}

	flowPtr := &models.Flow{
		Id:      fmt.Sprintf("%s_%d", parentPtr.Id, len(parentPtr.Childrens)),
		Name:    name,
		Message: []*models.MessageItem{},
		Status:  models.FLOW_STATUS_WAIT,
		Parent:  parentPtr,
		Level:   models.GetChildLevel(parentPtr.Level),
	}

	return s.appendFlow(flowPtr, parentPtr)
}

func (s *FlowManager) GetCurrentFlow() *models.Flow {
	return s.currentFlowPtr
}

func (s *FlowManager) GetLocalFlows() []*models.Flow {
	return s.localFlows
}

func (s *FlowManager) GetFlow(flowID string) *models.Flow {
	for _, flow := range s.localFlows {
		if flow.Id == flowID {
			return flow
		}
	}

	return nil
}

func (s *FlowManager) appendFlow(flowPtr *models.Flow, parentPtr *models.Flow) (ret *models.Flow) {
	defer func() {
		s.currentFlowPtr = flowPtr
		ret = s.currentFlowPtr
		log.Infof("appendFlow: %s", flowPtr)
	}()

	if parentPtr != nil {
		flowPtr.Parent = parentPtr
		parentPtr.Childrens = append(parentPtr.Childrens, flowPtr)
		return
	}

	if s.currentFlowPtr != nil {
		flowPtr.Parent = s.currentFlowPtr
		s.currentFlowPtr.Childrens = append(s.currentFlowPtr.Childrens, flowPtr)
		return
	}

	s.localFlows = append(s.localFlows, flowPtr)
	return
}

func (s *FlowManager) FinishFlow(stage string, flowPtr *models.Flow) (ret *models.Flow) {
	defer func() {
		if s.currentFlowPtr != nil {
			log.Infof("continue run flow %s", s.currentFlowPtr)
		} else {
			log.Infof("finish all flow")
		}

		ret = s.currentFlowPtr
	}()

	// 1. 检查当前flowPtr是否分解完成，如果没有分解完成，则要继续分解
	if !flowPtr.CanFinish() {
		param := map[string]string{
			"Name": flowPtr.GetName(),
		}
		subLevelCheckPromptVal, subLevelCheckPromptErr := iipc.ReadInnerSystemPromptsWithParam(fmt.Sprintf("sublevel-%s-check", stage), param)
		if subLevelCheckPromptErr != nil {
			log.Errorf("read sublevel-check prompt failed, err:%s", subLevelCheckPromptErr)
			return
		}
		flowPtr.AppendMessage(&models.MessageItem{
			Role:    common.User,
			Content: []string{subLevelCheckPromptVal},
		})
		return
	}

	s.currentFlowPtr = flowPtr
	for {
		s.currentFlowPtr = s.currentFlowPtr.GetParent()
		if s.currentFlowPtr == nil {
			break
		}

		if s.currentFlowPtr.Level == models.RootLevel {
			s.currentFlowPtr = nil
			return
		}

		// 2. 判断Parent是否可完成，如果可以完成，则继续取父级
		if s.currentFlowPtr.CanFinish() {
			continue
		}

		// 3. 深度优先，获取最小层级未完成的子流程
		s.currentFlowPtr = s.currentFlowPtr.GetLatestUnFinishFlow()
		if s.currentFlowPtr != nil {
			param := map[string]string{
				"Name": s.currentFlowPtr.GetName(),
			}

			subLevelCheckPromptVal, subLevelCheckPromptErr := iipc.ReadInnerSystemPromptsWithParam(fmt.Sprintf("sublevel-%s-check", stage), param)
			if subLevelCheckPromptErr != nil {
				log.Errorf("read sublevel-check prompt failed, err:%s", subLevelCheckPromptErr)
				return
			}
			s.currentFlowPtr.AppendMessage(&models.MessageItem{
				Role:    common.User,
				Content: []string{subLevelCheckPromptVal},
			})
		}
		return
	}
	return
}

func (s *FlowManager) CompleteFlow(flowPtr *models.Flow) (ret *models.Flow) {
	s.Reset()
	return
}
