package biz

import (
	"path/filepath"
	"strings"
	"testing"

	fu "github.com/muidea/magicCommon/foundation/util"
	"github.com/stretchr/testify/assert"
)

func TestSiblingPath(t *testing.T) {
	fileList := []string{
		"/home/rangh/code/demo/a1.md",
		"/home/rangh/code/demo/a1/a11.md",
		"/home/rangh/code/demo/a1/a12.md",
		"/home/rangh/code/demo/a1/a11/a111.md",
		"/home/rangh/code/demo/a1/a11/a112.md",
		"/home/rangh/code/demo/a1/a12/a121.md",
		"/home/rangh/code/demo/b2.md",
		"/home/rangh/code/demo/b2/b22.md",
		"/home/rangh/code/demo/b2/b22/b222.md",
		"/home/rangh/code/demo/c3.md",
		"/home/rangh/code/demo/c3/c33.md",
		"/home/rangh/code/demo/c3/c33/c333.md",
	}

	rootPath := "/home/rangh/code/demo"

	siblingPath := "."
	childPath := "."
	siblingList := fu.StringSet{}
	childList := fu.StringSet{}

	level := "project"
	name := "a1"

	funcVal := func(name, level string) {
		offset := -1
		switch level {
		case "project":
			offset = 0
		case "module":
			offset = 1
		case "block":
			offset = 2
		}

		for _, val := range fileList {
			filePath := val[len(rootPath)+1:]
			items := strings.Split(filePath, "/")
			if len(items) != offset+1 {
				continue
			}
			if offset == 0 {
				childList = childList.Add(filePath)
				continue
			}
			if items[offset-1] != name {
				continue
			}
			childPath = strings.Join(items[:offset], "/")
			childList = childList.Add(filePath)
		}
		if childPath != "." {
			siblingPath = filepath.Dir(childPath)
		}
		for _, val := range fileList {
			filePath := val[len(rootPath)+1:]
			fileDir := filepath.Dir(filePath)
			if siblingPath == fileDir && childPath != "." {
				siblingList = siblingList.Add(filePath)
			}
		}
	}

	funcVal(name, level)

	assert.Equal(t, ".", siblingPath, "siblingPath error")
	assert.Equal(t, 0, len(siblingList), "siblingList error")
	assert.Equal(t, ".", childPath, "childPath error")
	assert.Equal(t, 3, len(childList), "childList error")

	siblingPath = "."
	childPath = "."
	siblingList = fu.StringSet{}
	childList = fu.StringSet{}

	name = "a1"
	level = "module"
	funcVal(name, level)
	assert.Equal(t, ".", siblingPath, "siblingPath error")
	assert.Equal(t, 3, len(siblingList), "siblingList error")
	assert.Equal(t, "a1", childPath, "childPath error")
	assert.Equal(t, 2, len(childList), "childList error")

	siblingPath = "."
	childPath = "."
	siblingList = fu.StringSet{}
	childList = fu.StringSet{}

	name = "a11"
	level = "block"
	funcVal(name, level)
	assert.Equal(t, "a1", siblingPath, "siblingPath error")
	assert.Equal(t, 2, len(siblingList), "siblingList error")
	assert.Equal(t, "a1/a11", childPath, "childPath error")
	assert.Equal(t, 2, len(childList), "childList error")

}
