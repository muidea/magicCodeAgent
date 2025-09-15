package common

const EmbedtextInitiator = "/internal/initiators/embedtext"

type EmbedTextHelper interface {
	GetEmbedText() (ret EmbedTextWorker)
}

type EmbedTextWorker interface {
	EmbedText(chunkText string) (ret []float32, err error)
	EmbedTextBulk(chunkText []string) (ret [][]float32, err error)
	ExtractKeywords(chunkText string, topK int) (ret []string, err error)
}

// embed service structures
type EmbedRequest struct {
	Text string `json:"text"`
}
type EmbedResponse struct {
	Vector []float32 `json:"vector"`
}

type EmbedBulkRequest struct {
	Texts []string `json:"texts"`
}
type EmbedBulkResponse struct {
	Vectors [][]float32 `json:"vectors"`
}

type ExtractKeywordsRequest struct {
	Text string `json:"text"`
	TopK int    `json:"topK"`
}

type ExtractKeywordsResponse struct {
	Keywords []string `json:"keywords"`
}
