package common

import (
	"github.com/muidea/magicCodeAgent/pkg/models"
	cd "github.com/muidea/magicCommon/def"
)

const TemplatemanagerModule = "/internal/modules/blocks/templatemanager"

const (
	ListTemplate  = "/template/list"
	QueryTemplate = "/template/query"
)

type TemplateListResult struct {
	cd.Result
	Values []string `json:"values"`
}

type TemplateResult struct {
	cd.Result
	Value *models.Template `json:"value"`
}
