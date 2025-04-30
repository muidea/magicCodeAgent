package common

import (
	cd "github.com/muidea/magicCommon/def"

	"github.com/muidea/magicCodeAgent/pkg/models"
)

const TracerModule = "/internal/modules/blocks/tracer"

const (
	FilterMessage = "/message/filter"
	AppendMessage = "/message/append"
)

type MessageInfoListResult struct {
	cd.Result
	Values []*models.MessageInfo `json:"values"`
}

type Tracer interface {
	AppendMessage(msg *models.MessageInfo)
}
