package main

import (
	"fmt"
	fu "github.com/muidea/magicCommon/foundation/util"
	"github.com/muidea/magicToolkit/model"
	"path"
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

	defineFile := path.Join(rootPath, ".project", "define.json")
	err = fu.SaveConfig(defineFile, cfg)
	return
}
