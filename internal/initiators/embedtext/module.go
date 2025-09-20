package embedtext

import (
	"math"
	"net/http"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/log"
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

	url := config.GetNlpEmbedAPIUrl()
	_, err = fn.HTTPPost(http.DefaultClient, url, req, &rsp)
	if err != nil {
		log.Errorf("http post url:%s, chunkText:%s, error: %s", url, chunkText, err.Error())
		return
	}

	ret = s.normalize(rsp.Vector)
	return
}

func (s *Embedtext) EmbedTextBulk(chunkText []string) (ret [][]float32, err error) {
	req := &common.EmbedBulkRequest{Texts: chunkText}
	rsp := &common.EmbedBulkResponse{}

	url := config.GetNlpEmbedBatchAPIUrl()
	_, err = fn.HTTPPost(http.DefaultClient, url, req, &rsp)
	if err != nil {
		log.Errorf("http post url:%s, chunkText:%s, error: %s", url, chunkText, err.Error())
		return
	}

	for _, v := range rsp.Vectors {
		ret = append(ret, s.normalize(v))
	}
	return
}

func (s *Embedtext) ExtractKeywords(chunkText string, topK int, businessContext string, simThreshold float32) (ret []string, err error) {
	req := &common.ExtractKeywordsRequest{
		Text:            chunkText,
		TopK:            topK,
		BusinessContext: businessContext,
		Threshold:       simThreshold,
	}
	rsp := &common.ExtractKeywordsResponse{}

	url := config.GetNlpKeywordsAPIUrl()
	_, err = fn.HTTPPost(http.DefaultClient, url, req, &rsp)
	if err != nil {
		log.Errorf("http post url:%s, chunkText:%s, error: %s", url, chunkText, err.Error())
		return
	}
	ret = rsp.Keywords
	return
}

func (s *Embedtext) normalize(vec []float32) []float32 {
	norm := float32(0.0)
	for _, v := range vec {
		norm += v * v
	}
	if norm == 0 {
		return vec
	}
	inv := 1.0 / float32(math.Sqrt(float64(norm)))
	out := make([]float32, len(vec))
	for i, v := range vec {
		out[i] = v * inv
	}
	return out
}
