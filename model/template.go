package model

const (
	// DirItem 目录项
	DirItem = iota
	// FileItem 文件项
	FileItem
)

// Item item
type Item struct {
	Name     string  `json:"name"`
	Type     int     `json:"type"`
	Value    string  `json:"value"`
	SubItems []*Item `json:"subItems"`
}

// IsDir is directory
func (s *Item) IsDir() bool {
	return s.Type == DirItem
}

// IsFile is file
func (s *Item) IsFile() bool {
	return s.Type == FileItem
}

// Template template
type Template struct {
	Name  string  `json:"name"`
	Items []*Item `json:"items"`
}
