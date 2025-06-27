package util

import (
	"sort"
	"sync"
)

// 文件列表管理
// 对文件列表的进行管理，提供有序文件访问
type FileManager struct {
	fileMutex sync.RWMutex
	files     []string
}

func (s *FileManager) Store(filePath string) {
	s.fileMutex.Lock()
	defer s.fileMutex.Unlock()

	// 使用二分查找确定插入位置
	i := sort.SearchStrings(s.files, filePath)
	if i < len(s.files) && s.files[i] == filePath {
		return // 已存在则不重复添加
	}
	s.files = append(s.files, "")
	copy(s.files[i+1:], s.files[i:])
	s.files[i] = filePath
}

func (s *FileManager) Delete(filePath string) {
	s.fileMutex.Lock()
	defer s.fileMutex.Unlock()

	i := sort.SearchStrings(s.files, filePath)
	if i < len(s.files) && s.files[i] == filePath {
		s.files = append(s.files[:i], s.files[i+1:]...)
	}
}

func (s *FileManager) Range(f func(fileName string) bool) {
	s.fileMutex.RLock()
	defer s.fileMutex.RUnlock()

	for _, file := range s.files {
		if !f(file) {
			break
		}
	}
}
