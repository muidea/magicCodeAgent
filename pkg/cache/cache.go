package cache

import (
	"sync"
	"sync/atomic"

	"github.com/muidea/magicCodeAgent/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
)

// MessageSequenceCache 是一个MessageInfo内存顺序缓存
// 实现将MessageInfo顺序缓存到内存中，用于快速获取消息列表
// 实际最大缓存记录条数为1w条(暂定，可以通过配置项来调整)
type MessageSequenceCache struct {
	snVal    int64
	mu       sync.RWMutex
	messages []*models.MessageInfo
	maxSize  int
}

// NewMessageSequenceCache 创建新的缓存实例
func NewMessageSequenceCache(maxSize int) *MessageSequenceCache {
	if maxSize <= 0 {
		maxSize = 10000 // 默认1w条
	}
	return &MessageSequenceCache{
		messages: make([]*models.MessageInfo, 0, maxSize),
		maxSize:  maxSize,
	}
}

// Add 添加消息到缓存
func (c *MessageSequenceCache) Add(msg *models.MessageInfo) {
	msg.SN = atomic.AddInt64(&c.snVal, 1)
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.messages) >= c.maxSize {
		// 达到最大容量时，移除最早的消息
		c.messages = c.messages[1:]
	}
	c.messages = append(c.messages, msg)
}

// GetAll 获取所有缓存消息（按添加顺序）
func (c *MessageSequenceCache) GetAll() []*models.MessageInfo {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// 返回副本避免外部修改
	result := make([]*models.MessageInfo, len(c.messages))
	copy(result, c.messages)
	return result
}

// GetPage 分页获取缓存消息
func (c *MessageSequenceCache) GetPage(page, pageSize int) ([]*models.MessageInfo, int) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	total := len(c.messages)
	if total == 0 || pageSize <= 0 {
		return []*models.MessageInfo{}, total
	}

	// 计算总页数
	totalPages := (total + pageSize - 1) / pageSize
	if page <= 0 {
		page = 1
	} else if page > totalPages {
		return []*models.MessageInfo{}, total
	}

	// 计算起始和结束索引
	start := (page - 1) * pageSize
	end := start + pageSize
	if end > total {
		end = total
	}

	// 返回分页数据和总数
	result := make([]*models.MessageInfo, end-start)
	copy(result, c.messages[start:end])
	return result, total
}

func (s *MessageSequenceCache) FilterMessage(filter *common.MessageFilter) (ret []*models.MessageInfo, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	pageSize := filter.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	pageIndex := filter.PageIndex
	if pageIndex <= 0 {
		pageIndex = 1
	}

	itemCount := 0
	offsetVal := (pageIndex - 1) * pageSize
	for _, msg := range s.messages {
		if filter.Workspace != "" {
			if msg.Workspace != filter.Workspace {
				continue
			}
		}
		itemCount++
		if itemCount < offsetVal {
			continue
		}

		ret = append(ret, msg)
		if len(ret) >= pageSize {
			break
		}
	}
	return
}

// Clear 清空缓存
func (c *MessageSequenceCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.messages = make([]*models.MessageInfo, 0, c.maxSize)
}

// Size 获取当前缓存大小
func (c *MessageSequenceCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.messages)
}
