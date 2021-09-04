package common

import (
	casCommon "github.com/muidea/magicCas/common"
	commonDef "github.com/muidea/magicCommon/def"
)

const (
	ViewSettingProfile = "/setting/profile/query/"

	ViewSetting = "/setting/query/"
)

const SettingModule = "/kernel/setting"

type Content struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type ProfileResult struct {
	commonDef.Result
	Profile    *casCommon.AccountView `json:"profile"`
	Total      int64                  `json:"total"`
	OperateLog []*LogView             `json:"operateLog"`
}

type SettingResult struct {
	commonDef.Result
	Item []*Content `json:"item"`
}
