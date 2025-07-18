package config

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/muidea/magicCommon/foundation/log"
	"github.com/stretchr/testify/assert"
)

func TestModuleName(t *testing.T) {
	moduleName := "abc123"

	assert.True(t, namePattern.MatchString(moduleName))

	moduleName = "abc123"
	assert.True(t, namePattern.MatchString(moduleName))

	moduleName = "abc123"
	assert.True(t, namePattern.MatchString(moduleName))

	moduleName = "service/web/abc123"
	assert.False(t, namePattern.MatchString(moduleName))

	moduleName = "afdf/abc123"
	assert.False(t, namePattern.MatchString(moduleName))

	moduleName = "123abc/abc123"
	assert.False(t, namePattern.MatchString(moduleName))
}

func TestCheckSuffix(t *testing.T) {
	workingPath := "/home/rangh/codespace/magicCodeAgent/internal/initiators"
	codeRootPath := workingPath
	workingPath = filepath.ToSlash(workingPath)
	log.Infof("working path:%s", workingPath)
	suffixes := []string{
		"/internal/initiators",
		"/internal/modules/kernel",
		"/internal/modules/blocks",
	}
	for _, suffix := range suffixes {
		if strings.HasSuffix(workingPath, suffix) {
			codeRootPath = strings.TrimSuffix(workingPath, suffix)
			codeRootPath = strings.TrimRight(codeRootPath, "/")
			break
		}
	}

	log.Infof("codeRootPath:%s", codeRootPath)
}
