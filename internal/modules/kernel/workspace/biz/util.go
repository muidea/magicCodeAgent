package biz

import (
	"path/filepath"
	"strings"
	"time"
)

type HybridEvent struct {
	EventType string // "filesystem" 或 "git"
	Path      string
	Operation string // create/update/delete 或 commit
	Timestamp time.Time
}

func (s HybridEvent) String() string {
	return s.EventType + " " + s.Operation + " " + s.Path + " " + s.Timestamp.String()
}

type EventHandler func(HybridEvent)

// 判断是否忽略
// .git .tmp 隐藏文件都忽略，不用处理
func isIgnored(path string) bool {
	base := filepath.Base(path)
	return strings.Contains(path, ".git") ||
		strings.HasSuffix(path, ".tmp") ||
		strings.HasPrefix(base, ".")
}
