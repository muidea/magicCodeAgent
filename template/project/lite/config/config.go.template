package config

import (
	log "github.com/cihub/seelog"
	fu "github.com/muidea/magicCommon/foundation/util"
)

var configItem *CfgItem

const cfgPath = "/var/app/config/cfg.json"

func init() {
	cfg := &CfgItem{}
	err := fu.LoadConfig(cfgPath, cfg)
	if err != nil {
		log.Errorf("load config failed, err:%s", err.Error())
		return
	}

	configItem = cfg
}

type CfgItem struct {
}
