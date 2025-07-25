package biz

import (
	"bytes"
	"strings"
	"sync"
)

const defaultMaxBufferLinesSize = 50

// LogBuffer 封装调试缓冲区功能
type LogBuffer struct {
	bufferMutex        sync.Mutex
	bufferLines        []string
	maxBufferLinesSize int
	partialLine        string // 用于存储不完整的行
}

// NewLogBuffer 创建一个新的调试缓冲区
func NewLogBuffer(maxLineSize int) *LogBuffer {
	if maxLineSize <= 0 {
		maxLineSize = defaultMaxBufferLinesSize // 默认缓冲20行记录
	}
	return &LogBuffer{
		maxBufferLinesSize: maxLineSize,
	}
}

func (s *LogBuffer) Reset() {
	s.bufferMutex.Lock()
	defer s.bufferMutex.Unlock()

	s.bufferLines = []string{}
	s.partialLine = ""
}

// Write 实现io.Writer接口，将数据按行写入缓冲区
func (s *LogBuffer) Write(p []byte) (n int, err error) {
	s.bufferMutex.Lock()
	defer s.bufferMutex.Unlock()

	// 将新数据与之前的不完整行合并
	combined := s.partialLine + string(p)

	// 按换行符分割
	lines := bytes.Split([]byte(combined), []byte("\n"))

	// 最后一行可能是不完整的，保存起来
	s.partialLine = string(lines[len(lines)-1])

	// 处理完整的行（除了最后一行）
	for i := 0; i < len(lines)-1; i++ {
		line := string(lines[i])
		s.bufferLines = append(s.bufferLines, line)
	}

	// 如果缓冲区超过最大大小，移除最旧的记录
	if len(s.bufferLines) > s.maxBufferLinesSize {
		s.bufferLines = s.bufferLines[len(s.bufferLines)-s.maxBufferLinesSize:]
	}

	return len(p), nil
}

// SetBufferSize 设置缓冲区大小
func (s *LogBuffer) SetBufferSize(size int) {
	if size <= 0 {
		return
	}

	s.bufferMutex.Lock()
	defer s.bufferMutex.Unlock()

	s.maxBufferLinesSize = size

	// 如果当前缓冲区超过新的大小限制，则裁剪
	if len(s.bufferLines) > s.maxBufferLinesSize {
		s.bufferLines = s.bufferLines[len(s.bufferLines)-s.maxBufferLinesSize:]
	}
}

// GetFunctionRecords 查询包含指定函数名的记录
func (s *LogBuffer) GetFunctionRecords(functionName string) []string {
	s.bufferMutex.Lock()
	defer s.bufferMutex.Unlock()

	var results []string
	for _, line := range s.bufferLines {
		// 检查记录是否包含指定的函数名
		if functionName == "" || strings.Contains(line, functionName) {
			results = append(results, line)
		}
	}

	return results
}

// GetAllRecords 获取所有缓冲的记录
func (s *LogBuffer) GetAllRecords() []string {
	s.bufferMutex.Lock()
	defer s.bufferMutex.Unlock()

	// 返回缓冲区记录的副本
	result := make([]string, len(s.bufferLines))
	copy(result, s.bufferLines)

	return result
}

// GetLatestRecords 获取最新的N行记录
func (s *LogBuffer) GetLatestRecords(count int) []string {
	s.bufferMutex.Lock()
	defer s.bufferMutex.Unlock()

	// 如果请求的行数小于等于0，返回空切片
	if count <= 0 {
		return []string{}
	}

	// 如果缓冲区为空且没有不完整行，返回空切片
	if len(s.bufferLines) == 0 && s.partialLine == "" {
		return []string{}
	}

	// 计算起始索引
	startIndex := len(s.bufferLines) - count
	if startIndex < 0 {
		startIndex = 0 // 如果请求的行数超过缓冲区大小，则从头开始
	}

	// 获取最新的N行记录
	result := make([]string, len(s.bufferLines)-startIndex)
	copy(result, s.bufferLines[startIndex:])

	// 如果存在不完整行，将其添加到结果中
	if s.partialLine != "" {
		result = append(result, s.partialLine)
	}

	return result
}

func (s *LogBuffer) IsEmpty() bool {
	s.bufferMutex.Lock()
	defer s.bufferMutex.Unlock()
	return len(s.bufferLines) == 0 && s.partialLine == ""
}
