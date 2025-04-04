package util

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	cd "github.com/muidea/magicCommon/def"
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

func GetProjectTemplatePath(typeVal string) (ret string, err error) {
	toolKitPath := os.Getenv("MAGICTOOLKIT_PATH")
	if toolKitPath == "" {
		err = fmt.Errorf("no found toolkitPath")
		return
	}

	templatePath := path.Join(toolKitPath, "template", "project", typeVal)
	existsFlag := Exists(templatePath)
	if !existsFlag {
		err = fmt.Errorf("toolkit project template does not exist")
		return
	}

	ret = templatePath
	return
}

func GetModuleTemplatePath() (ret string, err error) {
	toolKitPath := os.Getenv("MAGICTOOLKIT_PATH")
	if toolKitPath == "" {
		err = fmt.Errorf("no found toolkitPath")
		return
	}

	templatePath := path.Join(toolKitPath, "template", "module")
	existsFlag := Exists(templatePath)
	if !existsFlag {
		err = fmt.Errorf("toolkit module template does not exist")
		return
	}

	ret = templatePath
	return
}

func GetTemplateContent(templatePath, filePath, fileName string) (ret string, err error) {
	var fullFilePath string
	if filePath == "" {
		fullFilePath = fmt.Sprintf("%s/%s", templatePath, fileName)
	} else {
		fullFilePath = fmt.Sprintf("%s/%s/%s", templatePath, filePath, fileName)
	}

	existsFlag := Exists(fullFilePath)
	if !existsFlag {
		err = fmt.Errorf("template file not exist, file:%s", fullFilePath)
		return
	}

	contentVal, contentErr := os.ReadFile(fullFilePath)
	if contentErr != nil {
		err = contentErr
		return
	}

	ret = string(contentVal)
	return
}

// VerifyModuleWorkPath 校验工作路径是否合法
// 1.工作路径的父级路径必须是以下几种：xxx/xxx/internal/initators,xxx/xxx/internal/modules,xxx/xxx/internal/modules/blocks,xxx/xxx/xxx/internal/modules/kernel
// 2. 其他路径非法
// 3. 工作路径不能是空
// 4. 工作路径不能是相对路径
// 5. 非法路径举例： /home/xxx/internal/initators/xxx/
func VerifyModuleWorkPath(pathVal string) *cd.Error {
	if pathVal == "" {
		return cd.NewError(cd.IllegalParam, "work path cannot be empty")
	}

	// 检查相对路径
	if !filepath.IsAbs(pathVal) {
		return cd.NewError(cd.IllegalParam, "work path must be absolute path")
	}

	validParentPaths := []string{
		filepath.Join("internal", "initators"),
		filepath.Join("internal", "modules"),
		filepath.Join("internal", "modules", "blocks"),
		filepath.Join("internal", "modules", "kernel"),
	}

	// 标准化路径
	normalizedPath := filepath.Clean(pathVal)
	for _, validPath := range validParentPaths {
		// 检查路径是否以validPath结尾
		if strings.HasSuffix(normalizedPath, validPath) {
			return nil
		}
	}

	return cd.NewError(cd.IllegalParam,
		fmt.Sprintf("invalid work path %s, parent path must be one of: %v", normalizedPath, validParentPaths))
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
