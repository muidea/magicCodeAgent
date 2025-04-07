package biz

import (
	"io/fs"
	"path"
	"testing"

	"github.com/muidea/magicCommon/foundation/log"

	_ "github.com/muidea/magicToolkit"
	"github.com/muidea/magicToolkit/internal/config"
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
