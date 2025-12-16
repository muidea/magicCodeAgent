package biz

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/muidea/magicCodeAgent/pkg/models"
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

	// 创建 modules/kernel 子目录
	kernelPath := filepath.Join(tmpDir, "internal", "modules", "kernel")
	err = os.MkdirAll(kernelPath, 0755)
	assert.NoError(t, err)

	// 创建 modules/blocks 子目录
	blocksPath := filepath.Join(tmpDir, "internal", "modules", "blocks")
	err = os.MkdirAll(blocksPath, 0755)
	assert.NoError(t, err)

	// 非 Web 模板
	nonWebTemplate := &models.Template{Category: "service"}
	// Web 模板
	webTemplate := &models.Template{Category: "web"}

	// 测试用例1：Web 模板，返回 "src/pages"
	t.Run("web template returns src/pages", func(t *testing.T) {
		s := &Module{}
		ret, cdErr := s.verifyWorkingPath(kernelPath, webTemplate)
		assert.Nil(t, cdErr)
		assert.Equal(t, "src/pages", ret)
	})

	// 测试用例2：非 Web 模板，路径以 /modules/kernel 结尾，找到 go.mod
	t.Run("non-web kernel path with go.mod", func(t *testing.T) {
		s := &Module{}
		ret, cdErr := s.verifyWorkingPath(kernelPath, nonWebTemplate)
		assert.Nil(t, cdErr)
		// 期望的相对路径是 "internal/modules/kernel"
		expected := filepath.Join("internal", "modules", "kernel")
		assert.Equal(t, expected, ret)
	})

	// 测试用例3：非 Web 模板，路径以 /modules/blocks 结尾，找到 go.mod
	t.Run("non-web blocks path with go.mod", func(t *testing.T) {
		s := &Module{}
		ret, cdErr := s.verifyWorkingPath(blocksPath, nonWebTemplate)
		assert.Nil(t, cdErr)
		expected := filepath.Join("internal", "modules", "blocks")
		assert.Equal(t, expected, ret)
	})

	// 测试用例4：非 Web 模板，路径不以允许的后缀结尾
	t.Run("non-web path with invalid suffix", func(t *testing.T) {
		s := &Module{}
		invalidPath := filepath.Join(tmpDir, "internal", "something")
		_, cdErr := s.verifyWorkingPath(invalidPath, nonWebTemplate)
		assert.NotNil(t, cdErr)
		if cdErr != nil {
			assert.Equal(t, int(cd.IllegalParam), int(cdErr.Code))
			assert.Contains(t, cdErr.Message, "must be in /modules/kernel or /modules/blocks")
		}
	})

	// 测试用例5：非 Web 模板，路径以 /modules/kernel 结尾，但找不到 go.mod
	t.Run("non-web kernel path without go.mod", func(t *testing.T) {
		// 创建一个没有 go.mod 的单独临时目录
		noGoModDir := t.TempDir()
		noGoModKernel := filepath.Join(noGoModDir, "internal", "modules", "kernel")
		err := os.MkdirAll(noGoModKernel, 0755)
		assert.NoError(t, err)

		s := &Module{}
		_, cdErr := s.verifyWorkingPath(noGoModKernel, nonWebTemplate)
		assert.NotNil(t, cdErr)
		if cdErr != nil {
			assert.Equal(t, int(cd.IllegalParam), int(cdErr.Code))
			assert.Contains(t, cdErr.Message, "must be a go module")
		}
	})

	// 测试用例6：空路径
	t.Run("empty path", func(t *testing.T) {
		s := &Module{}
		_, cdErr := s.verifyWorkingPath("", nonWebTemplate)
		assert.NotNil(t, cdErr)
		if cdErr != nil {
			assert.Equal(t, int(cd.IllegalParam), int(cdErr.Code))
			assert.Contains(t, cdErr.Message, "working path is empty")
		}
	})

	// 测试用例7：相对路径
	t.Run("relative path", func(t *testing.T) {
		s := &Module{}
		_, cdErr := s.verifyWorkingPath("internal/modules/kernel", nonWebTemplate)
		assert.NotNil(t, cdErr)
		if cdErr != nil {
			assert.Equal(t, int(cd.IllegalParam), int(cdErr.Code))
			assert.Contains(t, cdErr.Message, "must be absolute path")
		}
	})

	// 测试用例8：go.mod 在祖父目录中
	t.Run("go.mod in grandparent directory", func(t *testing.T) {
		// 在 tmpDir 下创建更深层的目录
		deepPath := filepath.Join(tmpDir, "a", "b", "c", "modules", "kernel")
		err := os.MkdirAll(deepPath, 0755)
		assert.NoError(t, err)

		s := &Module{}
		ret, cdErr := s.verifyWorkingPath(deepPath, nonWebTemplate)
		assert.Nil(t, cdErr)
		expected := filepath.Join("a", "b", "c", "modules", "kernel")
		assert.Equal(t, expected, ret)
	})

	// 测试用例9：路径包含多余的斜杠
	t.Run("path with extra slashes", func(t *testing.T) {
		s := &Module{}
		extraSlashPath := kernelPath + "//"
		ret, cdErr := s.verifyWorkingPath(extraSlashPath, nonWebTemplate)
		assert.Nil(t, cdErr)
		expected := filepath.Join("internal", "modules", "kernel")
		assert.Equal(t, expected, ret)
	})
}
