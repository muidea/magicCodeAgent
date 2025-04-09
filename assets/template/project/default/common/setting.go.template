package common

import (
	cc "github.com/muidea/magicCas/common"
	cd "github.com/muidea/magicCommon/def"
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
	cd.Result
	Profile    *cc.AccountView `json:"profile"`
	Total      int64           `json:"total"`
	OperateLog []*LogView      `json:"operateLog"`
}

type SettingResult struct {
	cd.Result
	Item []*Content `json:"item"`
}
