package biz

import "github.com/muidea/magicCodeAgent/pkg/models"

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

func (s *FlowManager) AppendFlow(flowPtr *models.Flow, parentPtr *models.Flow) {
	defer func() {
		s.currentFlowPtr = flowPtr
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
}

func (s *FlowManager) FinishFlow(flowPtr *models.Flow) {
	// 1. 将当前flowPtr的前一个兄弟节点置为currentFlowPtr
	if parent := flowPtr.GetParent(); parent != nil {
		children := parent.GetChildrens()
		for i, child := range children {
			if child == flowPtr && i > 0 {
				s.currentFlowPtr = children[i-1]
				break
			}
		}
	}

	// 2. 将当前flowPtr的status置为FLOW_STATUS_FINISH
	flowPtr.Status = models.FLOW_STATUS_FINISH
}
