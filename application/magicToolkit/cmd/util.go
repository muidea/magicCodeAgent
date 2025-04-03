package main

import (
	"fmt"
	"os"
	"path"

	fp "github.com/muidea/magicCommon/foundation/path"
	fu "github.com/muidea/magicCommon/foundation/util"
	"github.com/muidea/magicToolkit/model"
)

func loadProjectDefine(rootPath string) (ret *model.Define, err error) {
	defineFile := path.Join(rootPath, ".project", "define.json")
	ptr := &model.Define{}
	err = fu.LoadConfig(defineFile, ptr)
	if err != nil {
		return
	}

	ret = ptr
	return
}

func saveProjectDefine(rootPath string, cfg *model.Define) (err error) {
	if cfg == nil {
		err = fmt.Errorf("illegal project define")
		return
	}

	filePath := path.Join(rootPath, ".project")
	if !fp.Exist(filePath) {
		err = os.MkdirAll(filePath, os.ModePerm)
	}
	if err != nil {
		return
	}

	defineFile := path.Join(filePath, "define.json")
	err = fu.SaveConfig(defineFile, cfg)
	return
}
