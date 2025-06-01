package biz

import (
	"fmt"

	"github.com/muidea/magicCommon/foundation/log"

	"github.com/muidea/magicCodeAgent/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
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

func (s *FlowManager) newRootFlow(title string) *models.Flow {
	flowPtr := &models.Flow{
		Id:        fmt.Sprintf("flow_%d", len(s.localFlows)),
		Title:     title,
		Message:   []*models.MessageItem{},
		Status:    models.FLOW_STATUS_WAIT,
		Parent:    nil,
		Childrens: make([]*models.Flow, 0),
		Level:     models.RootLevel,
	}

	return s.appendFlow(flowPtr, nil)
}

func (s *FlowManager) NewFlow(title string, parentPtr *models.Flow) *models.Flow {
	if parentPtr == nil {
		parentPtr = s.currentFlowPtr
	}

	if parentPtr == nil {
		return s.newRootFlow(title)
	}

	flowPtr := &models.Flow{
		Id:      fmt.Sprintf("%s_%d", parentPtr.Id, len(parentPtr.Childrens)),
		Title:   title,
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

func (s *FlowManager) FinishFlow(flowPtr *models.Flow) (ret *models.Flow) {
	defer func() {
		if s.currentFlowPtr != nil {
			log.Infof("continue run flow %s", s.currentFlowPtr)
		} else {
			log.Infof("finish all flow")
		}

		ret = s.currentFlowPtr
	}()

	// 1. 检查当前flowPtr是否分解完成，如果没有分解完成，则要继续分解
	hasChild := len(flowPtr.Childrens) > 0
	if flowPtr.GetLevel() != models.BlockLevel && !hasChild {
		subFlowPtr := s.NewFlow("继续处理下一级功能", flowPtr)
		// 子层级还未分解，不能结束
		subFlowPtr.AppendMessage(&models.MessageItem{
			Role: common.User,
			Content: []string{
				"继续处理下一级功能",
			},
		})
		return
	}

	// 这里只直接返回上一层，而不是继续根据Flow状态进行判断的原因是由于大模型并不能按照预期一次性完成所有的分解
	// 直接返回上一次可以让他根据当前的情况，重新验算一遍，以便有遗漏，但是带来的坏处是多执行了一些步骤，会额外占用Token开销
	if flowPtr.CanFinish() {
		flowPtr.Status = models.FLOW_STATUS_FINISH
	}
	parent := flowPtr.GetParent()
	// 这里优先处理尚未完成的子Flow
	s.currentFlowPtr = parent
	if parent != nil {
		for _, child := range parent.GetChildrens() {
			if !child.CanFinish() {
				s.currentFlowPtr = child
				break
			}
		}
		s.currentFlowPtr.AppendMessage(&models.MessageItem{
			Role: common.User,
			Content: []string{
				"继续检查并处理剩余功能",
			},
		})
	}

	return
}
