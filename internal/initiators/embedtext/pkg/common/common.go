package common

const EmbedtextInitiator = "/internal/initiators/embedtext"

type EmbedTextHelper interface {
	GetEmbedText() (ret EmbedTextWorker)
}

type EmbedTextWorker interface {
	EmbedText(chunkText string) (ret []float32, err error)
}

// embed service structures
type EmbedRequest struct {
	Text string `json:"text"`
}
type EmbedResponse struct {
	Vector []float32 `json:"vector"`
}
