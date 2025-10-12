package biz

import (
	"strings"

	"github.com/muidea/magicCommon/event"

	"github.com/muidea/magicCodeAgent/internal/config"
	imkkc "github.com/muidea/magicCodeAgent/internal/modules/kernel/knowledgebase/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
)

type KnowledgeList []string

type LocalKnowledge interface {
	ReCallKnowledge(workspacePtr *models.Workspace, flowPtr *models.Flow, thinkingContent string) (KnowledgeList, error)
}

func (s *Workspace) ReCallKnowledge(workspacePtr *models.Workspace, flowPtr *models.Flow, thinkingContent string) (KnowledgeList, error) {
	workspaceSetting := config.WorkspaceSetting()
	paramPtr := &imkkc.HybridSearchKnowledgeParam{
		Warehouse:       workspacePtr.KnowledgeWarehouse,
		Text:            thinkingContent,
		TopK:            workspaceSetting.TopK,
		BusinessContext: strings.Join(workspacePtr.TechStack, ","),
		Similarity:      workspaceSetting.VectSimilarity,
		ContextType:     imkkc.ContextFAQ,
	}

	rcEvent := event.NewEvent(imkkc.HybridSearchWarehouseKnowledge, s.ID(), imkkc.KnowledgeBaseModule, nil, paramPtr)
	rcResult := s.SendEvent(rcEvent)
	chunkList, chunkErr := event.GetAs[[]*imkkc.Chunk](rcResult)
	if chunkErr != nil {
		return nil, chunkErr
	}

	knowledgeList := KnowledgeList{}
	for _, chunk := range chunkList {
		if chunk.Text == "" {
			continue
		}
		knowledgeList = append(knowledgeList, chunk.Text)
	}

	return knowledgeList, nil
}
