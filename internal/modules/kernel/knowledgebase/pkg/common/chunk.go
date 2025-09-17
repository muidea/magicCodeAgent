package common

import (
	"fmt"
	"regexp"
	"strings"
)

type Chunk struct {
	ID         int64             `json:"id"`
	Text       string            `json:"text"`
	FilePath   string            `json:"filePath"`
	Metadata   map[string]string `json:"metadata"`
	Vector     []float32         `json:"-"`
	Similarity float64           `json:"similarity"`
}

func (s *Chunk) String() string {
	return fmt.Sprintf("filePath:%s, vector len:%d", s.FilePath, len(s.Vector))
}

// -------------------- Markdown 分片（章节优先，章节超长才切分，保证行完整） --------------------
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
		if text == "" {
			buffer = nil
			return
		}

		if len([]rune(text)) <= chunkSize || chunkSize <= 0 {
			// 不超过 chunkSize，直接作为一个chunk
			chunks = append(chunks, &Chunk{
				Text:     text,
				FilePath: filePath,
				Metadata: map[string]string{
					"source":   filePath,
					"section":  currentSection,
					"secIndex": fmt.Sprintf("%d", sectionIndex),
					"chunkId":  fmt.Sprintf("%d", chunkIndex),
				},
			})
			chunkIndex++
		} else {
			// 超过 chunkSize，再进行切分（保证行完整）
			parts := splitIntoChunksByLine(text, chunkSize, overlap)
			for _, p := range parts {
				if strings.TrimSpace(p) == "" {
					continue
				}
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

// -------------------- 保证行完整的切分 --------------------
func splitIntoChunksByLine(text string, size, overlap int) []string {
	if size <= 0 {
		if strings.TrimSpace(text) != "" {
			return []string{text}
		}
		return nil
	}

	lines := strings.Split(text, "\n")
	var result []string
	var buffer []string
	currentLen := 0

	for _, line := range lines {
		lineRunes := []rune(line)
		lineLen := len(lineRunes)

		// 如果单行比 chunkSize 还长，强制切这行
		if lineLen > size {
			if len(buffer) > 0 {
				result = append(result, strings.TrimSpace(strings.Join(buffer, "\n")))
				buffer = nil
				currentLen = 0
			}
			for start := 0; start < lineLen; start += size {
				end := start + size
				if end > lineLen {
					end = lineLen
				}
				part := strings.TrimSpace(string(lineRunes[start:end]))
				if part != "" {
					result = append(result, part)
				}
			}
			continue
		}

		// 如果加上这一行超过 chunkSize，先输出当前buffer
		if currentLen+lineLen+1 > size { // +1 代表换行
			if len(buffer) > 0 {
				result = append(result, strings.TrimSpace(strings.Join(buffer, "\n")))
			}
			// overlap 简单实现：保留最后 overlap 个字符到下一个chunk
			if overlap > 0 && len(result) > 0 {
				last := result[len(result)-1]
				runes := []rune(last)
				if overlap < len(runes) {
					buffer = []string{string(runes[len(runes)-overlap:])}
					currentLen = len(buffer[0])
				} else {
					buffer = nil
					currentLen = 0
				}
			} else {
				buffer = nil
				currentLen = 0
			}
		}

		buffer = append(buffer, line)
		currentLen += lineLen + 1
	}

	if len(buffer) > 0 {
		chunk := strings.TrimSpace(strings.Join(buffer, "\n"))
		if chunk != "" {
			result = append(result, chunk)
		}
	}

	return result
}
