package common

import (
	"fmt"
	"regexp"
	"strings"
)

type Chunk struct {
	ID       int64
	Text     string
	FilePath string
	Metadata map[string]string
	Vector   []float32
}

// -------------------- Markdown 分片（分层 + overlap） --------------------
func SplitMarkdown(filePath, content string, chunkSize, overlap int) []*Chunk {
	lines := strings.Split(content, "\n")

	var chunks []*Chunk
	var currentSection string
	var buffer []string
	sectionIndex := 0
	chunkIndex := 0

	headingRegex := regexp.MustCompile(`^#{1,6}\s+`)

	flushBuffer := func() {
		if len(buffer) == 0 {
			return
		}
		text := strings.TrimSpace(strings.Join(buffer, "\n"))
		parts := splitIntoChunks(text, chunkSize, overlap)
		for _, p := range parts {
			chunks = append(chunks, &Chunk{
				Text:     p,
				FilePath: filePath,
				Metadata: map[string]string{
					"source":   filePath,
					"section":  currentSection,
					"secIndex": fmt.Sprintf("%d", sectionIndex),
					"chunkId":  fmt.Sprintf("%d", chunkIndex),
				},
			})
			chunkIndex++
		}
		buffer = nil
	}

	for _, line := range lines {
		if headingRegex.MatchString(line) {
			flushBuffer()
			currentSection = strings.TrimSpace(line)
			sectionIndex++
		} else {
			buffer = append(buffer, line)
		}
	}
	flushBuffer()
	return chunks
}

func splitIntoChunks(text string, size, overlap int) []string {
	if size <= 0 {
		return []string{text}
	}
	runes := []rune(text)
	var result []string
	step := size - overlap
	if step <= 0 {
		step = size
	}
	for start := 0; start < len(runes); start += step {
		end := start + size
		if end > len(runes) {
			end = len(runes)
		}
		part := strings.TrimSpace(string(runes[start:end]))
		if part != "" {
			result = append(result, part)
		}
		if end == len(runes) {
			break
		}
	}
	return result
}
