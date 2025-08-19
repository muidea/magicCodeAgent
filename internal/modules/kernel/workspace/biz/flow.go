package biz

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/muidea/magicCommon/foundation/log"

	"github.com/muidea/magicCodeAgent/pkg/models"
)

// 实现将类似
// ./demoWarehouse-项目功能.md
// ./demoWarehouse/Warehouse-功能模块.md
// 这样的路径清理成 demoWarehouse, demoWarehouse/Warehouse
// 需要考虑./是否存在的情况
func cleanPath(path string) (ret string) {
	// 去除前缀 "./"（如果存在）
	path = strings.TrimPrefix(path, "./")

	// 获取不带扩展名的文件名部分（例如 "Warehouse-功能模块"）
	base := filepath.Base(path)
	base = strings.TrimSuffix(base, filepath.Ext(base))

	// 判断是否是顶层路径（无斜杠），如 demoWarehouse-项目功能.md
	// 若是顶层，返回文件名主干中 "-" 前的部分
	if !strings.Contains(path, "/") {
		if idx := strings.Index(base, "-"); idx != -1 {
			return base[:idx]
		}
		return base
	}

	// 否则返回目录 + 文件主干中 "-" 前的部分
	dir := filepath.Dir(path)
	if idx := strings.Index(base, "-"); idx != -1 {
		base = base[:idx]
	}
	return filepath.Join(dir, base)
}

func getParentPath(path string) (ret string) {
	idx := strings.LastIndex(path, "/")
	if idx == -1 {
		return ""
	}

	return path[:idx]
}

func getID(parentPtr *models.Flow) string {
	return fmt.Sprintf("%s_%d", parentPtr.Id, len(parentPtr.Childrens))
}

type FlowManager struct {
	localFlows     []*models.Flow
	currentFlowPtr *models.Flow
}

func NewFlowManager(rootName string) *FlowManager {
	ptr := &FlowManager{
		localFlows:     make([]*models.Flow, 0),
		currentFlowPtr: nil,
	}

	ptr.localFlows = append(ptr.localFlows, &models.Flow{
		Id:      "0",
		Name:    rootName,
		Level:   models.RootLevel,
		Stage:   "",
		Status:  models.StatusDoing,
		Path:    "*",
		Message: make([]*models.MessageItem, 0),
	})
	ptr.currentFlowPtr = ptr.localFlows[0]
	return ptr
}

func (s *FlowManager) insertFlow(current *models.Flow, stage string, parts []string, depth int) {
	if depth >= len(parts) {
		return
	}

	name := parts[depth]
	level := ""
	switch depth {
	case 0:
		level = "project"
	case 1:
		level = "module"
	default:
		level = "block"
	}

	// 查找是否已存在该子节点
	var child *models.Flow
	for _, c := range current.Childrens {
		if c.Name == name && c.Stage == stage && c.Level == level {
			child = c
			break
		}
	}

	// 不存在就创建新节点
	if child == nil {
		child = &models.Flow{
			Id:      getID(current),
			Name:    name,
			Level:   level,
			Stage:   stage,
			Status:  models.StatusInit,
			Path:    strings.Join(parts, "/"),
			Message: []*models.MessageItem{},
			Parent:  current,
		}

		s.localFlows = append(s.localFlows, child)
		current.Childrens = append(current.Childrens, child)
	}

	// 递归处理下一层
	s.insertFlow(child, stage, parts, depth+1)
}

func (s *FlowManager) InitFlow(flowPath []string, stage string) (err error) {
	if stage == models.DesignStage {
		return
	}

	for _, path := range flowPath {
		trimmed := strings.TrimSuffix(path, "-项目功能.md")
		trimmed = strings.TrimSuffix(trimmed, "-业务模块.md")
		trimmed = strings.TrimSuffix(trimmed, "-功能模块.md")
		parts := strings.Split(trimmed, "/")
		s.insertFlow(s.currentFlowPtr, stage, parts, 0)
	}
	return
}

func (s *FlowManager) Reset() {
	s.localFlows = s.localFlows[:1]
	s.currentFlowPtr = s.localFlows[0]
	s.currentFlowPtr.Childrens = []*models.Flow{}
}

func (s *FlowManager) getID(parentPtr *models.Flow) string {
	return fmt.Sprintf("%s_%d", parentPtr.Id, len(parentPtr.Childrens))
}

func (s *FlowManager) GetStageRoot(stage string) (ret *models.Flow) {
	rootPtr := s.localFlows[0]
	for _, val := range rootPtr.Childrens {
		if val.Stage == stage {
			ret = val
			break
		}
	}
	return
}

func (s *FlowManager) GetCurrentFlow() *models.Flow {
	return s.currentFlowPtr
}

func (s *FlowManager) GetLocalFlows() []*models.Flow {
	return s.localFlows
}

func (s *FlowManager) SwitchFlow(name, level, path, stage string) (ret *models.Flow, err error) {
	defer func() {
		// 这里对ret进行校验，避免出现越级的情况，以防止模型跑偏
		if ret == nil {
			return
		}
		parentPtr := ret.GetParent()
		if parentPtr == nil {
			return
		}

		if parentPtr.Status == models.StatusInit {
			switch stage {
			case models.DesignStage:
				err = errors.New("出现越级设计，需要重新调整")
			case models.CodeStage:
				err = errors.New("出现越级开发，需要重新调整")
			}
		}
	}()

	finalPath := cleanPath(path)
	// 只有DesignStage才需要添加名称，这是因为设计阶段path是输出对象不是参照对象
	if stage == models.DesignStage {
		finalPath = filepath.Join(finalPath, name)
	}

	parentPath := getParentPath(finalPath)
	if level == models.ProjectLevel {
		parentPath = "*"
	}
	log.Infof("name:%s, level:%s, path:%s, parentPath:%s, stage:%s, flowSize:%d", name, level, finalPath, parentPath, stage, len(s.localFlows))
	var currentPtr *models.Flow
	var parentPtr *models.Flow
	idx := 0
	for {
		flow := s.localFlows[idx]
		log.Infof("idx:%d, %v", idx, flow)
		if flow.GetName() == name && flow.GetPath() == finalPath && flow.GetLevel() == level && flow.GetStage() == stage {
			currentPtr = flow
			break
		}
		if flow.GetPath() == parentPath && (flow.GetStage() == stage || flow.GetStage() == "") {
			parentPtr = flow
		}
		idx++
		if idx >= len(s.localFlows) {
			break
		}
	}
	if currentPtr != nil {
		if currentPtr.Status > models.StatusDoing {
			log.Errorf("任务已经执行，需要核对输入参数name:%s, level:%s, path:%s是否正确", name, level, finalPath)
			err = fmt.Errorf("任务%s已经执行完成，继续处理其他任务", currentPtr.GetName())
			return
		}

		currentPtr.Status = models.StatusDoing
		s.currentFlowPtr = currentPtr
		ret = currentPtr
		return
	}

	if parentPtr == nil {
		log.Errorf("父级任务不存在，需要核对输入参数name:%s, level:%s, path:%s是否正确", name, level, finalPath)
		err = fmt.Errorf("父级任务不存在，需要核对输入参数path:%s是否正确", path)
		return
	}

	if models.GetChildLevel(parentPtr.GetLevel()) != level {
		log.Errorf("越级任务，需要核对输入参数name:%s, level:%s, path:%s是否正确", name, level, finalPath)
		err = fmt.Errorf("越级任务，需要核对输入参数level:%s是否正确", level)
		return
	}

	var flowPtr *models.Flow
	for _, val := range parentPtr.Childrens {
		if val.Name == name {
			flowPtr = val
		}
	}

	if flowPtr == nil {
		flowPtr = &models.Flow{
			Id:      s.getID(parentPtr),
			Name:    name,
			Level:   level,
			Stage:   stage,
			Path:    finalPath,
			Message: []*models.MessageItem{},
			Parent:  parentPtr,
		}
		parentPtr.Childrens = append(parentPtr.Childrens, flowPtr)
		s.localFlows = append(s.localFlows, flowPtr)
	}
	flowPtr.Status = models.StatusDoing

	s.currentFlowPtr = flowPtr
	ret = flowPtr
	return
}

func (s *FlowManager) FinishFlow(name, level, path string, needContinue bool, stage string) (ret *models.Flow, err error) {
	defer func() {
		if s.currentFlowPtr != nil {
			log.Infof("continue run flow %s", s.currentFlowPtr)
		} else {
			log.Infof("finish all flow")
		}

		ret = s.currentFlowPtr
	}()

	finalPath := cleanPath(path)
	log.Infof("name:%s, level:%s, path:%s, finalPath:%s, stage:%s, flowSize:%d", name, level, path, finalPath, stage, len(s.localFlows))

	var finishFlow *models.Flow
	idx := 0
	for {
		flow := s.localFlows[idx]
		log.Infof("idx:%d, %v", idx, flow)
		if flow.GetName() == name && flow.GetLevel() == level && flow.GetPath() == finalPath && flow.GetStage() == stage {
			finishFlow = flow
			break
		}

		idx++
		if idx >= len(s.localFlows) {
			break
		}
	}

	if finishFlow == nil {
		err = fmt.Errorf("找不到对应的任务, 需要核对输入参数name:%s,level:%s, path:%s，提供正确的参数", name, level, path)
		return
	}

	if needContinue || !finishFlow.CanFinish() {
		finishFlow.Status = models.StatusFinish
	} else {
		finishFlow.Status = models.StatusDone
	}

	// 1. 检查当前flowPtr是否分解完成，如果没有分解完成，则要继续分解
	if !finishFlow.CanFinish() {
		s.currentFlowPtr = finishFlow
		return
	}

	s.currentFlowPtr = finishFlow.GetParent()
	if s.currentFlowPtr == nil || s.currentFlowPtr.Level == models.RootLevel {
		s.currentFlowPtr = nil
		return
	}

	subFlowPtr := s.currentFlowPtr.GetLatestUnFinishFlow()
	if subFlowPtr != nil {
		s.currentFlowPtr = subFlowPtr
		return
	}

	return
}

func (s *FlowManager) CompleteFlow(flowPtr *models.Flow) (ret *models.Flow) {
	s.Reset()
	return
}
