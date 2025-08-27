package embedtext

import (
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/framework/plugin/initiator"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicCodeAgent/internal/initiators/embedtext/pkg/common"
)

func init() {
	initiator.Register(New())
}

type Embedtext struct {}

func New() *Embedtext {
	return &Embedtext{}
}

func (s *Embedtext) ID() string {
	return common.EmbedtextInitiator
}

func (s *Embedtext) Setup(eventHub event.Hub, backgroundRoutine task.BackgroundRoutine) (err *cd.Error) {
	return nil
}

func (s *Embedtext) Run() (err *cd.Error) {
	return
}

func (s *Embedtext) Teardown() {
}
