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

	// 已经有分解，则要校验所属ParentFlow是否有未完成
	for {
		if flowPtr == nil {
			// 所有的Flow都处理完成，返回nil
			s.currentFlowPtr = nil
			return
		}

		// 先将当前Flow状态至为FINISH
		// 再检查其Parent是否可以完成，如果不能完成，则从中挑选要给未完成的继续
		if flowPtr.CanFinish() {
			flowPtr.Status = models.FLOW_STATUS_FINISH
			flowPtr = flowPtr.GetParent()
			log.Infof("get parent flow：%s", flowPtr)
			continue
		}
		break
	}

	// 到这里说明flowPtr不能完成，要其中找挑选一个未完成的继续
	allFinished := true
	for {
		childs := flowPtr.GetChildrens()
		for _, childPtr := range childs {
			if childPtr.CanFinish() {
				childPtr.Status = models.FLOW_STATUS_FINISH
				continue
			}

			flowPtr = childPtr
			log.Infof("get child flow：%s", flowPtr)
			allFinished = false
			break
		}

		if allFinished {
			break
		}
	}

	s.currentFlowPtr = flowPtr
	return
}
