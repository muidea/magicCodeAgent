package biz

import (
	"io/fs"
	"path"
	"strings"
	"testing"

	"github.com/muidea/magicCommon/foundation/log"
	"github.com/stretchr/testify/assert"

	_ "github.com/muidea/magicCodeAgent"
	"github.com/muidea/magicCodeAgent/internal/config"
)

func TestTempleFS(t *testing.T) {
	const templateBasePath = "template/module"
	embeddedFS := config.GetTemplateFS()
	fs.WalkDir(embeddedFS, templateBasePath, func(filePath string, fileEntry fs.DirEntry, err error) error {
		if templateBasePath == filePath || fileEntry.IsDir() || err != nil {
			return nil
		}

		relativePath, fileName := path.Split(filePath[len(templateBasePath)+1:])
		//relativePath := path[len(basePath)+1:]
		//fileName := d.Name()
		log.Infof("path:%s, relativePath:%v, fileName:%v", filePath, relativePath, fileName)
		return nil
	})
}

func TestFileName(t *testing.T) {
	templateFileName := ".umirc.ts.template"
	items := strings.Split(templateFileName, ".")
	fileName := strings.Join(items[:len(items)-1], ".")
	fileExt := items[len(items)-1]
	assert.Equal(t, "template", fileExt)
	assert.Equal(t, ".umirc.ts", fileName)
}
