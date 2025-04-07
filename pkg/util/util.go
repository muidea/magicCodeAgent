package util

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicToolkit/internal/config"
)

const (
	ProjectDefaultTemplatePath = "default"
)

func Exists(itemPath string) bool {
	_, err := os.Stat(itemPath)
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}

	return true
}

func ToLower(val string) string {
	return strings.ToLower(val)
}

func ToTitle(val string) string {
	caser := cases.Title(language.English)
	return caser.String(val)
}

func ToCamelCase(val string) string {
	if len(val) == 1 {
		return strings.ToLower(val)
	}

	return strings.ToLower(val[0:1]) + val[1:]
}

func GetTemplateContent(templatePath, filePath, fileName string) (ret string, err error) {
	fullFilePath := path.Join(templatePath, filePath, fileName)
	if _, fileErr := fs.Stat(config.GetTemplateFS(), fullFilePath); fileErr != nil {
		err = fileErr
		return
	}

	contentVal, contentErr := fs.ReadFile(config.GetTemplateFS(), fullFilePath)
	if contentErr != nil {
		err = contentErr
		return
	}

	ret = string(contentVal)
	return
}

var sep = string(filepath.Separator)
var validModuleParentPaths = []string{
	filepath.Join("internal", "modules"),
	filepath.Join("internal", "modules", "blocks"),
	filepath.Join("internal", "modules", "kernel"),
}

// VerifyModuleWorkPath 校验工作路径是否合法，并返回路径中ProjectName
// 1.工作路径的父级路径必须是以下几种：xxx/[ProjectName]/internal/modules,xxx/[ProjectName]/internal/modules/blocks,xxx/xxx/[ProjectName]/internal/modules/kernel
// 2. 其他路径非法
// 3. 工作路径不能是空
// 4. 工作路径不能是相对路径
// 5. 非法路径举例： /home/xxx/internal/modules/xxx/
func VerifyModuleWorkPath(pathVal string) (projectName, prefixPath string, err *cd.Error) {
	projectName, prefixPath, err = verifyWorkPath(pathVal, validModuleParentPaths)
	return
}

var validInitatorParentPaths = []string{
	filepath.Join("internal", "initators"),
}

// VerifyInitatorWorkPath 校验工作路径是否合法
// 1.工作路径的父级路径必须是以下几种：xxx/xxx/internal/initators
// 2. 其他路径非法
// 3. 工作路径不能是空
// 4. 工作路径不能是相对路径
// 5. 非法路径举例： /home/xxx/internal/modules/xxx/
func VerifyInitatorWorkPath(pathVal string) (projectName, prefixPath string, err *cd.Error) {
	projectName, prefixPath, err = verifyWorkPath(pathVal, validInitatorParentPaths)
	return
}

func verifyWorkPath(pathVal string, validParentPath []string) (projectName, prefixPath string, err *cd.Error) {
	if pathVal == "" {
		return "", "", cd.NewError(cd.IllegalParam, "work path cannot be empty")
	}

	// 检查相对路径
	if !filepath.IsAbs(pathVal) {
		return "", "", cd.NewError(cd.IllegalParam, "work path must be absolute path")
	}

	// 标准化路径
	normalizedPath := filepath.Clean(pathVal)
	for _, validPath := range validParentPath {
		// 检查路径是否以validPath结尾
		if strings.HasSuffix(normalizedPath, validPath) {
			prefix := normalizedPath[:len(normalizedPath)-len(validPath)-1]
			projectName = prefix[strings.LastIndex(prefix, sep)+1:]
			prefixPath = validPath
			return
		}
	}

	return "", "", cd.NewError(cd.IllegalParam,
		fmt.Sprintf("invalid work path %s, parent path must be one of: %v", normalizedPath, validParentPath))
}

// VerifyProjectWorkPath 校验项目工作目录
// 1. 当前目录下无子目录或子文件
// 2. 当天目录下同时存在如下文件go.mod,go.sum,Makefile, 及application,internal,pkg,vendor目录
func VerifyProjectWorkPath(pathVal string) *cd.Error {
	// 空路径视为合法
	if pathVal == "" {
		return cd.NewError(cd.IllegalParam, "project work path cannot be empty")
	}

	// 检查目录是否为空
	entries, readErr := os.ReadDir(pathVal)
	if readErr != nil {
		return cd.NewError(cd.IllegalParam, fmt.Sprintf("failed to read directory: %v", readErr))
	}
	if len(entries) == 0 {
		return nil
	}

	// 检查相对路径
	if !filepath.IsAbs(pathVal) {
		return cd.NewError(cd.IllegalParam, "work path must be absolute path")
	}

	// 检查必需文件
	requiredFiles := []string{"go.mod", "go.sum", "Makefile"}
	for _, file := range requiredFiles {
		filePath := filepath.Join(pathVal, file)
		if !Exists(filePath) {
			return cd.NewError(cd.IllegalParam, fmt.Sprintf("required file %s not found in project", file))
		}
	}

	// 检查必需目录
	requiredDirs := []string{"application", "internal", "pkg", "vendor"}
	for _, dir := range requiredDirs {
		dirPath := filepath.Join(pathVal, dir)
		if !Exists(dirPath) {
			return cd.NewError(cd.IllegalParam, fmt.Sprintf("required directory %s not found in project", dir))
		}
	}

	return nil
}
