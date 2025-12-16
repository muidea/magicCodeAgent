package biz

import (
	"os"
	"path/filepath"
	"testing"

	cd "github.com/muidea/magicCommon/def"
	"github.com/stretchr/testify/assert"
)

func TestVerifyWorkingPath(t *testing.T) {
	// 创建临时目录作为测试根目录
	tmpDir := t.TempDir()

	// 在临时目录中创建 go.mod
	goModPath := filepath.Join(tmpDir, "go.mod")
	err := os.WriteFile(goModPath, []byte("module test\n"), 0644)
	assert.NoError(t, err)

	// 创建 initiators 子目录
	initiatorsPath := filepath.Join(tmpDir, "internal", "initiators")
	err = os.MkdirAll(initiatorsPath, 0755)
	assert.NoError(t, err)

	// 测试用例1：有效路径，go.mod 在父目录
	t.Run("valid path with go.mod in parent", func(t *testing.T) {
		s := &Initiator{}
		ret, err := s.verifyWorkingPath(initiatorsPath)
		assert.Nil(t, err)
		// 期望的相对路径是 "internal/initiators"
		expected := filepath.Join("internal", "initiators")
		assert.Equal(t, expected, ret)
	})

	// 测试用例2：路径不以 /initiators 结尾
	t.Run("path not ending with initiators", func(t *testing.T) {
		s := &Initiator{}
		invalidPath := filepath.Join(tmpDir, "internal", "something")
		_, cdErr := s.verifyWorkingPath(invalidPath)
		assert.NotNil(t, cdErr)
		if cdErr != nil {
			assert.Equal(t, int(cd.IllegalParam), int(cdErr.Code))
		}
	})

	// 测试用例3：非绝对路径
	t.Run("relative path", func(t *testing.T) {
		s := &Initiator{}
		_, cdErr := s.verifyWorkingPath("internal/initiators")
		assert.NotNil(t, cdErr)
		if cdErr != nil {
			assert.Equal(t, int(cd.IllegalParam), int(cdErr.Code))
		}
	})

	// 测试用例4：空路径
	t.Run("empty path", func(t *testing.T) {
		s := &Initiator{}
		_, cdErr := s.verifyWorkingPath("")
		assert.NotNil(t, cdErr)
		if cdErr != nil {
			assert.Equal(t, int(cd.IllegalParam), int(cdErr.Code))
		}
	})

	// 测试用例5：go.mod 不在任何父目录中
	t.Run("no go.mod in parent directories", func(t *testing.T) {
		// 创建一个没有 go.mod 的单独临时目录
		noGoModDir := t.TempDir()
		noGoModInitiators := filepath.Join(noGoModDir, "sub", "initiators")
		err := os.MkdirAll(noGoModInitiators, 0755)
		assert.NoError(t, err)

		s := &Initiator{}
		_, cdErr := s.verifyWorkingPath(noGoModInitiators)
		assert.NotNil(t, cdErr)
		if cdErr != nil {
			assert.Equal(t, int(cd.IllegalParam), int(cdErr.Code))
			assert.Contains(t, cdErr.Message, "no go.mod found")
		}
	})

	// 测试用例6：go.mod 在祖父目录中
	t.Run("go.mod in grandparent directory", func(t *testing.T) {
		// 在 tmpDir 下创建更深层的目录
		deepPath := filepath.Join(tmpDir, "a", "b", "c", "initiators")
		err := os.MkdirAll(deepPath, 0755)
		assert.NoError(t, err)

		s := &Initiator{}
		ret, err := s.verifyWorkingPath(deepPath)
		assert.Nil(t, err)
		expected := filepath.Join("a", "b", "c", "initiators")
		assert.Equal(t, expected, ret)
	})

	// 测试用例7：路径包含多余的斜杠和清理
	t.Run("path with extra slashes", func(t *testing.T) {
		s := &Initiator{}
		extraSlashPath := initiatorsPath + "//"
		ret, err := s.verifyWorkingPath(extraSlashPath)
		assert.Nil(t, err)
		expected := filepath.Join("internal", "initiators")
		assert.Equal(t, expected, ret)
	})
}
