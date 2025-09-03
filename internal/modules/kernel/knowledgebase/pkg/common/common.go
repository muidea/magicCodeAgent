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

type ListWarehouseResult struct {
	cd.Result
	Total  int64    `json:"total"`
	Values []string `json:"values"`
}

type CreateWarehouseParam struct {
	Name string `json:"name"`
}

type CreateWarehouseResult struct {
	cd.Result
}

type DestroyWarehouseParam struct {
	Name string `json:"name"`
}
type DestroyWarehouseResult struct {
	cd.Result
}

type KnowledgeItem struct {
	ID        uint32 `json:"id"`
	Name      string `json:"name"`
	FilePath  string `json:"filePath"`
	Warehouse string `json:"warehouse"`
	Status    int    `json:"status"`
	CreatedAt int64  `json:"createdAt"`
}

type ListKnowledgeResult struct {
	cd.Result
	Values []KnowledgeItem `json:"values"`
}
