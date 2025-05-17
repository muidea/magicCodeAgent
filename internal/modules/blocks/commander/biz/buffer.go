package biz

import (
	"bytes"
	"strings"
	"sync"
)

const defaultMaxBufferLinesSize = 20

// CommandBuffer 封装调试缓冲区功能
type CommandBuffer struct {
	bufferMutex        sync.Mutex
	bufferLines        []string
	maxBufferLinesSize int
	partialLine        string // 用于存储不完整的行
}

// NewCommandBuffer 创建一个新的调试缓冲区
func NewCommandBuffer(maxLineSize int) *CommandBuffer {
	if maxLineSize <= 0 {
		maxLineSize = defaultMaxBufferLinesSize // 默认缓冲20行记录
	}
	return &CommandBuffer{
		maxBufferLinesSize: maxLineSize,
	}
}

func (s *CommandBuffer) Reset() {
	s.bufferMutex.Lock()
	defer s.bufferMutex.Unlock()

	s.bufferLines = []string{}
	s.partialLine = ""
}

// Write 实现io.Writer接口，将数据按行写入缓冲区
func (b *CommandBuffer) Write(p []byte) (n int, err error) {
	b.bufferMutex.Lock()
	defer b.bufferMutex.Unlock()

	// 将新数据与之前的不完整行合并
	combined := b.partialLine + string(p)

	// 按换行符分割
	lines := bytes.Split([]byte(combined), []byte("\n"))

	// 最后一行可能是不完整的，保存起来
	b.partialLine = string(lines[len(lines)-1])

	// 处理完整的行（除了最后一行）
	for i := 0; i < len(lines)-1; i++ {
		line := string(lines[i])
		b.bufferLines = append(b.bufferLines, line)
	}

	// 如果缓冲区超过最大大小，移除最旧的记录
	if len(b.bufferLines) > b.maxBufferLinesSize {
		b.bufferLines = b.bufferLines[len(b.bufferLines)-b.maxBufferLinesSize:]
	}

	return len(p), nil
}

// SetBufferSize 设置缓冲区大小
func (b *CommandBuffer) SetBufferSize(size int) {
	if size <= 0 {
		return
	}

	b.bufferMutex.Lock()
	defer b.bufferMutex.Unlock()

	b.maxBufferLinesSize = size

	// 如果当前缓冲区超过新的大小限制，则裁剪
	if len(b.bufferLines) > b.maxBufferLinesSize {
		b.bufferLines = b.bufferLines[len(b.bufferLines)-b.maxBufferLinesSize:]
	}
}

// GetFunctionRecords 查询包含指定函数名的记录
func (b *CommandBuffer) GetFunctionRecords(functionName string) []string {
	b.bufferMutex.Lock()
	defer b.bufferMutex.Unlock()

	var results []string
	for _, line := range b.bufferLines {
		// 检查记录是否包含指定的函数名
		if functionName == "" || strings.Contains(line, functionName) {
			results = append(results, line)
		}
	}

	return results
}

// GetAllRecords 获取所有缓冲的记录
func (b *CommandBuffer) GetAllRecords() []string {
	b.bufferMutex.Lock()
	defer b.bufferMutex.Unlock()

	// 返回缓冲区记录的副本
	result := make([]string, len(b.bufferLines))
	copy(result, b.bufferLines)

	return result
}

// GetLatestRecords 获取最新的N行记录
func (b *CommandBuffer) GetLatestRecords(count int) []string {
	b.bufferMutex.Lock()
	defer b.bufferMutex.Unlock()

	// 如果请求的行数小于等于0，返回空切片
	if count <= 0 {
		return []string{}
	}

	// 如果缓冲区为空且没有不完整行，返回空切片
	if len(b.bufferLines) == 0 && b.partialLine == "" {
		return []string{}
	}

	// 计算起始索引
	startIndex := len(b.bufferLines) - count
	if startIndex < 0 {
		startIndex = 0 // 如果请求的行数超过缓冲区大小，则从头开始
	}

	// 获取最新的N行记录
	result := make([]string, len(b.bufferLines)-startIndex)
	copy(result, b.bufferLines[startIndex:])

	// 如果存在不完整行，将其添加到结果中
	if b.partialLine != "" {
		result = append(result, b.partialLine)
	}

	return result
}
