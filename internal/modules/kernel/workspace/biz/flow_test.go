package biz

import (
	"testing"
)

func TestCleanPath(t *testing.T) {
	fileList := []string{
		"demoWarehouse-项目功能.md",
		"./demoWarehouse-项目功能.md",
		"demoWarehouse/Shelf-业务模块.md",
		"demoWarehouse/Warehouse/WarehouseUpdate-功能模块.md",
		"./demoWarehouse/Warehouse/WarehouseUpdate-功能模块.md",
		"./demoWarehouse/Shelf/ShelfCreate-功能模块.md",
		"demoWarehouse/Shelf/ShelfCreate-功能模块.md",
		"demoWarehouse/Shelf/ShelfCreate.md", // 无 "-" 的文件名
		"./demoWarehouse/Shelf/ShelfCreate.md",
	}

	expected := []string{
		"demoWarehouse",
		"demoWarehouse",
		"demoWarehouse/Shelf",
		"demoWarehouse/Warehouse/WarehouseUpdate",
		"demoWarehouse/Warehouse/WarehouseUpdate",
		"demoWarehouse/Shelf/ShelfCreate",
		"demoWarehouse/Shelf/ShelfCreate",
		"demoWarehouse/Shelf/ShelfCreate",
		"demoWarehouse/Shelf/ShelfCreate",
	}

	expectedParent := []string{
		"",
		"",
		"demoWarehouse",
		"demoWarehouse/Warehouse",
		"demoWarehouse/Warehouse",
		"demoWarehouse/Shelf",
		"demoWarehouse/Shelf",
		"demoWarehouse/Shelf",
		"demoWarehouse/Shelf",
	}
	for i, f := range fileList {
		got := cleanPath(f)
		if got != expected[i] {
			t.Errorf("cleanPath(%q) = %q, want %q", f, got, expected[i])
		}
	}

	for i, f := range expected {
		got := getParentPath(f)
		if got != expectedParent[i] {
			t.Errorf("getParentPath(%q) = %q, idex:%d, want %q", f, got, i, expected[i])
		}
	}
}
