package common

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChunk_String(t *testing.T) {
	tests := []struct {
		name     string
		chunk    *Chunk
		expected string
	}{
		{
			name: "empty vector",
			chunk: &Chunk{
				FilePath: "test.md",
				Vector:   []float32{},
			},
			expected: "filePath:test.md, vector len:0",
		},
		{
			name: "with vector",
			chunk: &Chunk{
				FilePath: "test.md",
				Vector:   []float32{0.1, 0.2, 0.3},
			},
			expected: "filePath:test.md, vector len:3",
		},
		{
			name: "nil vector",
			chunk: &Chunk{
				FilePath: "test.md",
				Vector:   nil,
			},
			expected: "filePath:test.md, vector len:0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.chunk.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSplitMarkdown_EmptyContent(t *testing.T) {
	chunks := SplitMarkdown("test.md", "", 100, 20)
	assert.Empty(t, chunks)
}

func TestSplitMarkdown_SingleSection(t *testing.T) {
	content := `# Title
This is some content.
More content here.`

	chunks := SplitMarkdown("test.md", content, 100, 20)
	assert.Len(t, chunks, 1)
	assert.Equal(t, "This is some content.\nMore content here.", chunks[0].Text)
	assert.Equal(t, "test.md", chunks[0].FilePath)
	assert.Equal(t, "# Title", chunks[0].Metadata["section"])
	assert.Equal(t, "1", chunks[0].Metadata["secIndex"])
	assert.Equal(t, "0", chunks[0].Metadata["chunkId"])
}

func TestSplitMarkdown_MultipleSections(t *testing.T) {
	content := `# Section 1
Content for section 1.

## Subsection 1.1
Content for subsection.

# Section 2
Content for section 2.`

	chunks := SplitMarkdown("test.md", content, 100, 20)
	assert.Len(t, chunks, 3)

	assert.Equal(t, "Content for section 1.", chunks[0].Text)
	assert.Equal(t, "# Section 1", chunks[0].Metadata["section"])
	assert.Equal(t, "1", chunks[0].Metadata["secIndex"])

	assert.Equal(t, "Content for subsection.", chunks[1].Text)
	assert.Equal(t, "## Subsection 1.1", chunks[1].Metadata["section"])
	assert.Equal(t, "2", chunks[1].Metadata["secIndex"])

	assert.Equal(t, "Content for section 2.", chunks[2].Text)
	assert.Equal(t, "# Section 2", chunks[2].Metadata["section"])
	assert.Equal(t, "3", chunks[2].Metadata["secIndex"])
}

func TestSplitMarkdown_LargeContent(t *testing.T) {
	// Create content that exceeds chunk size
	content := `# Large Section
This is a very long line that should be split into multiple chunks because it exceeds the specified chunk size limit.
The splitIntoChunksByLine function should handle this properly by breaking it into smaller pieces while maintaining line integrity.`

	chunks := SplitMarkdown("test.md", content, 50, 10)

	// The exact number of chunks depends on the splitting algorithm
	assert.True(t, len(chunks) > 1, "Should split into multiple chunks")

	// Verify that all chunks contain meaningful content
	for i, chunk := range chunks {
		assert.NotEmpty(t, strings.TrimSpace(chunk.Text), "Chunk %d should not be empty", i)
	}

	// Verify that the original content is preserved when concatenated
	allText := ""
	for _, chunk := range chunks {
		allText += chunk.Text + "\n"
	}

	// Check for key phrases instead of exact matches due to splitting
	assert.Contains(t, allText, "very long line")
	assert.Contains(t, allText, "chunk size limit")
	assert.Contains(t, allText, "split into")      // This phrase is split across chunks
	assert.Contains(t, allText, "multiple chunks") // This phrase is split across chunks
}

func TestSplitMarkdown_ZeroChunkSize(t *testing.T) {
	content := `# Test Section
Some content here.`

	chunks := SplitMarkdown("test.md", content, 0, 0)
	assert.Len(t, chunks, 1)
	assert.Equal(t, "Some content here.", chunks[0].Text)
}

func TestSplitMarkdown_OnlyHeadings(t *testing.T) {
	content := `# Heading 1
## Heading 2
### Heading 3`

	chunks := SplitMarkdown("test.md", content, 100, 20)
	assert.Empty(t, chunks) // No content between headings, so no chunks
}

func TestSplitIntoChunksByLine_EmptyText(t *testing.T) {
	result := splitIntoChunksByLine("", 100, 20)
	assert.Empty(t, result)
}

func TestSplitIntoChunksByLine_SingleLine(t *testing.T) {
	text := "This is a single line of text."
	result := splitIntoChunksByLine(text, 100, 20)
	assert.Len(t, result, 1)
	assert.Equal(t, text, result[0])
}

func TestSplitIntoChunksByLine_MultipleLines(t *testing.T) {
	text := "Line 1\nLine 2\nLine 3"
	result := splitIntoChunksByLine(text, 100, 5) // Use larger chunk size to fit all lines
	assert.Len(t, result, 1)                      // All lines should fit within chunk size
	assert.Equal(t, text, result[0])
}

func TestSplitIntoChunksByLine_MultipleLinesWithSmallChunk(t *testing.T) {
	text := "Line 1\nLine 2\nLine 3"
	result := splitIntoChunksByLine(text, 10, 2) // Small chunk size that forces splitting
	// Should split into multiple chunks due to small size
	assert.True(t, len(result) > 1, "Should split into multiple chunks with small chunk size")

	// Verify all content is preserved
	combined := strings.Join(result, "\n")
	assert.Contains(t, combined, "Line 1")
	assert.Contains(t, combined, "Line 2")
	assert.Contains(t, combined, "Line 3")
}

func TestSplitIntoChunksByLine_LongLineForcedSplit(t *testing.T) {
	text := "This is a very long line that needs to be split because it exceeds the chunk size limit."
	result := splitIntoChunksByLine(text, 30, 10)
	assert.Len(t, result, 3)

	// Check that the split preserves content
	combined := result[0] + result[1] + result[2]
	assert.Contains(t, combined, "This is a very long line")
	assert.Contains(t, combined, "exceeds the chunk size limit")
}

func TestSplitIntoChunksByLine_ZeroSize(t *testing.T) {
	text := "Some text content"
	result := splitIntoChunksByLine(text, 0, 0)
	assert.Len(t, result, 1)
	assert.Equal(t, text, result[0])
}

func TestSplitIntoChunksByLine_WithOverlap(t *testing.T) {
	text := "First line\nSecond line which is longer\nThird line"
	result := splitIntoChunksByLine(text, 25, 10)

	// Should split into chunks with overlap
	assert.True(t, len(result) >= 2)

	// Check that chunks don't exceed size limit
	for _, chunk := range result {
		assert.LessOrEqual(t, len([]rune(chunk)), 25)
	}
}

func TestSplitIntoChunksByLine_WhitespaceOnly(t *testing.T) {
	text := "   \n\t\n   "
	result := splitIntoChunksByLine(text, 100, 20)
	assert.Empty(t, result)
}

func TestSplitMarkdown(t *testing.T) {
	md := `
# 第一章
这是第一章的内容，不会超过 chunkSize。

# 第二章
这一章的内容比较多。
每一行都不算长。
但是整体加起来会超过 chunkSize，所以应该被切分成多个块。
再来几行文字，增加长度。
最后一行，收尾。

# 第三章
这一行非常非常非常非常非常非常非常非常非常非常非常非常非常非常非常非常长，
长到必须被强制切分，否则会超出 chunkSize 限制。
`

	chunks := SplitMarkdown("test.md", md, 40, 5) // 设置小一点方便测试
	for i, c := range chunks {
		t.Logf("Chunk %d:\n%s\n----", i, c.Text)
	}

	// 基本断言
	if len(chunks) == 0 {
		t.Errorf("期望至少切分出一些 chunk，结果为 0")
	}
}

func TestSplitIntoChunksByLine(t *testing.T) {
	text := strings.Repeat("a", 30) + "\n" + strings.Repeat("b", 30) + "\nshort line"
	parts := splitIntoChunksByLine(text, 20, 5)

	for i, p := range parts {
		t.Logf("Part %d: %q", i, p)
	}

	if len(parts) < 3 {
		t.Errorf("期望至少切分成 3 部分，结果=%d", len(parts))
	}
}
