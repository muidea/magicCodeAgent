package embedtext

import (
	"net/http"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	fn "github.com/muidea/magicCommon/foundation/net"
	"github.com/muidea/magicCommon/framework/plugin/initiator"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicCodeAgent/internal/initiators/embedtext/config"
	"github.com/muidea/magicCodeAgent/internal/initiators/embedtext/pkg/common"
)

func init() {
	initiator.Register(New())
}

type Embedtext struct{}

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

func (s *Embedtext) GetEmbedText() common.EmbedTextWorker {
	return s
}

func (s *Embedtext) EmbedText(chunkText string) (ret []float32, err error) {
	req := &common.EmbedRequest{Text: chunkText}
	rsp := &common.EmbedResponse{}

	_, err = fn.HTTPPost(http.DefaultClient, config.GetBgeAPIUrl(), req, &rsp)
	if err != nil {
		return
	}

	ret = rsp.Vector
	return
}
