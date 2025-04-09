package magictoolkit

import (
	"embed"

	"github.com/muidea/magicToolkit/internal/config"
)

//go:embed assets/**
var templateFS embed.FS

func init() {
	config.SetEmbedTemplateFS(templateFS)
}
