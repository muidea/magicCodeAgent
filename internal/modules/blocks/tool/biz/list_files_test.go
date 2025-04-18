package biz

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/task"
	"github.com/stretchr/testify/assert"
)

func TestListFiles(t *testing.T) {
	// 创建一个Tool实例用于测试
	eventHub := event.NewHub(100)
	backgroundRoutine := task.NewBackgroundRoutine(100)
	tool := New(eventHub, backgroundRoutine)

	// 保存原始的WorkspaceDir
	originalWorkspaceDir := WorkspaceDir
	// 临时修改WorkspaceDir为测试目录
	tempDir := t.TempDir()
	WorkspaceDir = tempDir
	// 确保测试结束后恢复原始值
	defer func() { WorkspaceDir = originalWorkspaceDir }()

	// 创建测试目录结构
	// tempDir/
	//   |- file1.txt
	//   |- file2.txt
	//   |- subdir/
	//        |- file3.txt
	//        |- subsubdir/
	//             |- file4.txt

	// 创建顶层文件
	file1Path := filepath.Join(tempDir, "file1.txt")
	file2Path := filepath.Join(tempDir, "file2.txt")
	err := os.WriteFile(file1Path, []byte("content1"), 0644)
	assert.NoError(t, err, "创建file1.txt应该成功")
	err = os.WriteFile(file2Path, []byte("content2"), 0644)
	assert.NoError(t, err, "创建file2.txt应该成功")

	// 创建子目录
	subdirPath := filepath.Join(tempDir, "subdir")
	err = os.Mkdir(subdirPath, 0755)
	assert.NoError(t, err, "创建子目录应该成功")

	// 创建子目录中的文件
	file3Path := filepath.Join(subdirPath, "file3.txt")
	err = os.WriteFile(file3Path, []byte("content3"), 0644)
	assert.NoError(t, err, "创建file3.txt应该成功")

	// 创建子子目录
	subsubdirPath := filepath.Join(subdirPath, "subsubdir")
	err = os.Mkdir(subsubdirPath, 0755)
	assert.NoError(t, err, "创建子子目录应该成功")

	// 创建子子目录中的文件
	file4Path := filepath.Join(subsubdirPath, "file4.txt")
	err = os.WriteFile(file4Path, []byte("content4"), 0644)
	assert.NoError(t, err, "创建file4.txt应该成功")

	// 测试用例1：非递归列出顶层目录文件
	results1, err := tool.ListFiles(".", false)
	assert.NoError(t, err, "非递归列出文件应该成功")
	assert.Len(t, results1, 3, "应该列出4个项目（2个文件和2个目录）")
	assert.Contains(t, results1, file1Path, "结果应该包含file1.txt")
	assert.Contains(t, results1, file2Path, "结果应该包含file2.txt")

	// 测试用例2：递归列出所有文件
	results2, err := tool.ListFiles(".", true)
	assert.NoError(t, err, "递归列出文件应该成功")
	assert.Len(t, results2, 4, "应该列出4个文件")
	assert.Contains(t, results2, file1Path, "结果应该包含file1.txt")
	assert.Contains(t, results2, file2Path, "结果应该包含file2.txt")
	assert.Contains(t, results2, file3Path, "结果应该包含file3.txt")
	assert.Contains(t, results2, file4Path, "结果应该包含file4.txt")

	// 测试用例3：列出子目录文件
	results3, err := tool.ListFiles("subdir", false)
	assert.NoError(t, err, "列出子目录文件应该成功")
	assert.Len(t, results3, 2, "应该列出2个项目（1个文件和1个目录）")
	assert.Contains(t, results3, file3Path, "结果应该包含file3.txt")

	// 测试用例4：递归列出子目录文件
	results4, err := tool.ListFiles("subdir", true)
	assert.NoError(t, err, "递归列出子目录文件应该成功")
	assert.Len(t, results4, 2, "应该列出2个文件")
	assert.Contains(t, results4, file3Path, "结果应该包含file3.txt")
	assert.Contains(t, results4, file4Path, "结果应该包含file4.txt")

	// 测试用例5：列出不存在的目录
	_, err = tool.ListFiles("nonexistent", false)
	assert.Error(t, err, "列出不存在的目录应该返回错误")

	// 测试用例6：列出文件而不是目录
	_, err = tool.ListFiles("file1.txt", false)
	assert.Error(t, err, "列出文件而不是目录应该返回错误")
}
