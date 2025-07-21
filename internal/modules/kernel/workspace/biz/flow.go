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

func (s *FlowManager) Reset() {
	s.localFlows = make([]*models.Flow, 0)
	s.currentFlowPtr = nil
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

const checkPrompt = `请按以下步骤依次执行，避免重复确认或漏检。仅在全部满足后执行attempt_completion结束任务。

---

## 1. 检查已完成项  
- 所有标记为完成的项是否真正完成？  
- 确认无错误、无异常、无遗漏。

## 2. 业务层级完整性检查（设计文档对齐）  
- **是否所有在设计文档中定义的业务模块都已建立功能模块设计？**  
  - 对照设计文档中的业务模块清单，逐项比对确认。  
  - 不允许遗漏任何业务模块。  
- 每个业务模块是否已拆分至少一个功能模块？  
- 所有功能模块是否已完成功能开发，且对齐原始业务意图？

## 3. 功能层级完成检查  
- 所有功能模块是否状态为“已完成”？  
- 是否实现正确、符合预期，未出现遗漏或冗余？  
- 每个功能模块是否可以回溯到具体业务需求？

## 4. 判断是否结束任务  
- 所有检查步骤均通过。  
- 所有业务模块覆盖完毕，未有缺失或延迟项。  
- 所有功能项完成、正确、可交付。

✅ 满足以上条件 → 执行attempt_completion
⛔ 否则 → 返回未完成项继续处理。
`

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
		childs := []*models.Flow{}
		for _, child := range parent.GetChildrens() {
			if !child.CanFinish() {
				childs = append(childs, child)
			}
		}
		if len(childs) > 0 {
			s.currentFlowPtr = childs[0]
		}
		//parent.Childrens = childs

		for {
			if s.currentFlowPtr == nil {
				break
			}

			if !s.currentFlowPtr.CanFinish() {
				subFlowPtr := s.NewFlow("继续处理下一级功能", s.currentFlowPtr)
				// 子层级还未分解，不能结束
				subFlowPtr.AppendMessage(&models.MessageItem{
					Role: common.User,
					Content: []string{
						"继续处理下一级功能",
					},
				})
				return
			}

			if s.currentFlowPtr.Level != "" {
				s.currentFlowPtr = s.currentFlowPtr.Parent
				continue
			}

			s.currentFlowPtr.AppendMessage(&models.MessageItem{
				Role: common.User,
				Content: []string{
					checkPrompt,
				},
			})
			break
		}
	}

	return
}

func (s *FlowManager) CompleteFlow(flowPtr *models.Flow) (ret *models.Flow) {
	s.Reset()
	return
}
