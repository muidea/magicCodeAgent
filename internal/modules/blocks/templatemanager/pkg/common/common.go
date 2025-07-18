package common

import (
	"github.com/muidea/magicCodeAgent/pkg/models"
	cd "github.com/muidea/magicCommon/def"
)

const TemplatemanagerModule = "/internal/modules/blocks/templatemanager"

const (
	ListTemplate = "/template/list/"
	LoadTemplate = "/template/load/"
)

type ListTemplateResult struct {
	cd.Result
	Values []string `json:"values"`
}

type LoadTemplateResult struct {
	cd.Result
	Value *models.Template `json:"value"`
}
