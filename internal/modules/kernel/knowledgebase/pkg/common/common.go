package common

import (
	cd "github.com/muidea/magicCommon/def"
)

const KnowledgeBaseModule = "/internal/modules/kernel/knowledgebase"

const (
	BasePath   = "knowledgebase"
	DataPath   = "data"
	StaticPath = "static"
)

const (
	IdleState = iota
	LoadedState
)

const (
	ListKnowledgebaseWarehouse    = "/knowledgebase/warehouse/list"
	CreateKnowledgebaseWarehouse  = "/knowledgebase/warehouse/create"
	DestroyKnowledgebaseWarehouse = "/knowledgebase/warehouse/destroy"
	UploadWorkspaceKnowledge      = "/knowledgebase/knowledge/upload"
	ImportWorkspaceKnowledge      = "/knowledgebase/knowledge/import"
	DeleteWorkspaceKnowledge      = "/knowledgebase/knowledge/delete"
	ListWorkspaceKnowledge        = "/knowledgebase/knowledge/list"
)

type WarehouseItem struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   int64  `json:"createdAt"`
}

type ListWarehouseResult struct {
	cd.Result
	Values []*WarehouseItem `json:"values"`
}

type CreateWarehouseParam struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateWarehouseResult ListWarehouseResult

type DestroyWarehouseResult ListWarehouseResult

type KnowledgeItem struct {
	ID          uint64 `json:"-"`
	Name        string `json:"name"`
	Description string `json:"description"`
	FilePath    string `json:"filePath"`
	Warehouse   string `json:"warehouse"`
	Status      int    `json:"status"`
	CreatedAt   int64  `json:"createdAt"`
}

type ImportKnowledgeParam struct {
	FilePath  string `json:"filePath"`
	Warehouse string `json:"warehouse"`
}

type ListKnowledgeResult struct {
	cd.Result
	Values []*KnowledgeItem `json:"values"`
}
