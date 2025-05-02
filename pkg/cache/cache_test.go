package cache

import (
	"testing"
	"time"

	"github.com/muidea/magicCodeAgent/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestMessageSequenceCache(t *testing.T) {
	cache := NewMessageSequenceCache(100)

	// 测试Add和Size
	for i := 0; i < 50; i++ {
		cache.Add(&models.MessageInfo{
			Item:       &models.MessageItem{Role: "user", Content: "test"},
			CreateTime: time.Now().Format(time.RFC3339),
		})
	}
	assert.Equal(t, 50, cache.Size())

	// 测试GetAll
	all := cache.GetAll()
	assert.Equal(t, 50, len(all))
	for i, msg := range all {
		assert.Equal(t, int64(i+1), msg.SN)
	}

	// 测试GetPage
	t.Run("TestGetPage", func(t *testing.T) {
		// 正常分页
		page1, total := cache.GetPage(1, 10)
		assert.Equal(t, 50, total)
		assert.Equal(t, 10, len(page1))
		assert.Equal(t, int64(1), page1[0].SN)

		// 最后一页不足
		page5, _ := cache.GetPage(5, 10)
		assert.Equal(t, 10, len(page5))

		// 超出页码
		page6, _ := cache.GetPage(6, 10)
		assert.Equal(t, 0, len(page6))

		// 无效页大小
		pageInvalid, _ := cache.GetPage(1, 0)
		assert.Equal(t, 0, len(pageInvalid))
	})

	// 测试Clear
	cache.Clear()
	assert.Equal(t, 0, cache.Size())
	assert.Equal(t, 0, len(cache.GetAll()))
}

func TestMessageSequenceCacheConcurrent(t *testing.T) {
	cache := NewMessageSequenceCache(1000)

	// 并发写入测试
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				cache.Add(&models.MessageInfo{
					Item:       &models.MessageItem{Role: "user", Content: "test"},
					CreateTime: time.Now().Format(time.RFC3339),
				})
			}
		}()
	}

	// 等待并发写入完成
	time.Sleep(1 * time.Second)
	assert.Equal(t, 1000, cache.Size())

	// 并发读取测试
	for i := 0; i < 10; i++ {
		go func() {
			page, _ := cache.GetPage(1, 100)
			assert.Equal(t, 100, len(page))
		}()
	}
}
