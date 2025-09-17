package biz

import "github.com/muidea/magicCodeAgent/pkg/models"

type KnowledgeList []string

type LocalKnowledge interface {
	ReCallKnowledge(workspacePtr *models.Workspace, flowPtr *models.Flow, thinkingContent string) (KnowledgeList, error)
}

func (s *Workspace) ReCallKnowledge(workspacePtr *models.Workspace, flowPtr *models.Flow, thinkingContent string) (KnowledgeList, error) {
	return nil, nil
}
