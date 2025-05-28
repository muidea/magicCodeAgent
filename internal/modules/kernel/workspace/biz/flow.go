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
		Level:     common.Project,
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
		ret = s.currentFlowPtr
	}()

	// 1. 检查当前flowPtr是否可以完成
	hasChild := len(flowPtr.Childrens) > 0
	if flowPtr.GetLevel() != common.Point && !hasChild {
		subFlowPtr := s.NewFlow("继续分解下一级功能", flowPtr)
		// 子层级还未分解，不能结束
		subFlowPtr.AppendMessage(&models.MessageItem{
			Role: common.User,
			Content: []string{
				"继续分解下一级功能",
			},
		})

		return
	}

	flowPtr.Status = models.FLOW_STATUS_FINISH

	// 2. 将当前flowPtr的其中一个未完成的兄弟节点置为currentFlowPtr
	var curPtr *models.Flow
	if parent := flowPtr.GetParent(); parent != nil {
		children := parent.GetChildrens()
		for _, child := range children {
			if !child.IsFinish() {
				curPtr = child
				break
			}
		}
		// 如果所有的兄弟阶段都完成了，则该Parent也置为完成
		if curPtr == nil {
			parent.Status = models.FLOW_STATUS_FINISH
		} else {
			s.currentFlowPtr = curPtr
			return
		}

		parent = parent.GetParent()
		for {
			if parent == nil {
				// 表示所有的Flow都处理完毕
				s.currentFlowPtr = nil
				return
			}

			if parent.IsFinish() {
				parent = parent.GetParent()
				continue
			}

			children := parent.GetChildrens()
			if len(children) == 0 {
				// 如果当前Flow没有完成，并且无子Flow,把则将当前Flow设置为当前Flow
				s.currentFlowPtr = parent
				return
			}

			for _, child := range children {
				if child.IsFinish() {
					continue
				}
				parent = child
				break
			}
		}
	}
	return
}
