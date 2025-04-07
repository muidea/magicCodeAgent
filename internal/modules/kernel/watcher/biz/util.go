package biz

import (
	"strings"
	"time"
)

type HybridEvent struct {
	EventType string // "filesystem" 或 "git"
	Path      string
	Operation string // create/update/delete 或 commit
	Timestamp time.Time
}

type EventHandler func(HybridEvent)

func isIgnored(path string) bool {
	return strings.Contains(path, ".git") ||
		strings.HasSuffix(path, ".tmp")
}
