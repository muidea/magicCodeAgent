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

// ContextType 表示当前查询的上下文类型
type ContextType int

const (
	ContextFAQ          ContextType = iota // 精确问答
	ContextOpenDialogue                    // 开放式对话
	ContextMixed                           // 混合场景
)

const (
	ListKnowledgebaseWarehouse     = "/knowledgebase/warehouse/list"
	CreateKnowledgebaseWarehouse   = "/knowledgebase/warehouse/create"
	DestroyKnowledgebaseWarehouse  = "/knowledgebase/warehouse/destroy"
	UploadWarehouseKnowledge       = "/knowledgebase/knowledge/upload"
	DeleteWarehouseKnowledge       = "/knowledgebase/knowledge/delete"
	ListWarehouseKnowledge         = "/knowledgebase/knowledge/list"
	GetWarehouseKnowledgeContent   = "/knowledgebase/knowledge/content"
	SearchWarehouseKnowledge       = "/knowledgebase/knowledge/search"
	HybridSearchWarehouseKnowledge = "/knowledgebase/knowledge/hybrid_search"
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

type GetKnowledgeContentResult struct {
	cd.Result
	Value string `json:"value"`
}

type SearchKnowledgeParam struct {
	Warehouse string   `json:"warehouse"`
	Text      []string `json:"text"`
	TopK      int      `json:"topK"`
	Threshold float32  `json:"threshold"`
}

type SearchKnowledgeResult struct {
	cd.Result
	Values []*Chunk `json:"values"`
}

type HybridSearchKnowledgeParam struct {
	Warehouse   string      `json:"warehouse"`
	Text        string      `json:"text"`
	TopK        int         `json:"topK"`
	Threshold   float32     `json:"threshold"`
	ContextType ContextType `json:"contextType"`
}

type HybridSearchKnowledgeResult struct {
	cd.Result
	Values []*Chunk `json:"values"`
}
