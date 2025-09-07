package common

const EmbedtextInitiator = "/internal/initiators/embedtext"

type EmbedTextHelper interface {
	GetEmbedText() (ret EmbedTextWorker)
}

type EmbedTextWorker interface {
	EmbedText(chunkText string) (ret []float64, err error)
}

// embed service structures
type EmbedRequest struct {
	Text string `json:"text"`
}
type EmbedResponse struct {
	Vector []float64 `json:"vector"`
}
