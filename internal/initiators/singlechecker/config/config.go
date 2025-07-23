package config

import "flag"

var forceRestart bool

func init() {
	flag.BoolVar(&forceRestart, "forceRestart", false, "force restart (optional, default: false)")
}

func IsForceRestart() bool {
	return forceRestart
}
