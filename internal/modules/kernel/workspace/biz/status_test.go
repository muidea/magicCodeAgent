package biz

import (
	"bytes"
	"testing"
)

func TestDumpFlowStatusTree(t *testing.T) {
	// 构造测试数据
	root := &FlowStatus{
		Name:   "Root",
		Level:  "High",
		Status: "Running",
		Childrens: []*FlowStatus{
			{
				Name:   "Child1",
				Level:  "Medium",
				Status: "Completed",
				Childrens: []*FlowStatus{
					{Name: "GrandChild1", Level: "Low", Status: "Pending"},
					{Name: "GrandChild2", Level: "Low", Status: "Failed"},
				},
			},
			{
				Name:   "Child2",
				Level:  "Medium",
				Status: "Running",
				Childrens: []*FlowStatus{
					{Name: "GrandChild3", Level: "Low", Status: "Success"},
				},
			},
		},
	}

	var sb bytes.Buffer
	DumpFlowStatusTree(root, &sb)
	t.Log("\n" + sb.String())
}

func TestNewFlowStatus(t *testing.T) {
	flowList := []string{
		"demoWarehouse-项目功能.md",
		"demoWarehouse/Shelf-业务模块.md",
		"demoWarehouse/Shelf/ShelfCreate-功能模块.md",
		"demoWarehouse/Shelf/ShelfDelete-功能模块.md",
		"demoWarehouse/Shelf/ShelfQuery-功能模块.md",
		"demoWarehouse/Shelf/ShelfUpdate-功能模块.md",
		"demoWarehouse/Warehouse-业务模块.md",
		"demoWarehouse/Warehouse/WarehouseCreate-功能模块.md",
		"demoWarehouse/Warehouse/WarehouseDelete-功能模块.md",
		"demoWarehouse/Warehouse/WarehouseQuery-功能模块.md",
		"demoWarehouse/Warehouse/WarehouseUpdate-功能模块.md",
	}
	root := NewFlowStatus(flowList)

	var buf bytes.Buffer
	DumpFlowStatusTree(root, &buf)
	t.Log("\n" + buf.String())
}
