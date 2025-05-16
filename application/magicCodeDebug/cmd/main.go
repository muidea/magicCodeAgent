package main

import (
	"flag"
	"fmt"

	"github.com/muidea/magicCommon/foundation/log"
	"github.com/muidea/magicCommon/framework/application"
	"github.com/muidea/magicCommon/framework/service"

	_ "github.com/muidea/magicCodeAgent"

	"github.com/muidea/magicCodeAgent/internal/config"

	_ "github.com/muidea/magicCodeAgent/internal/modules/blocks/recorder"
)

var (
	buildTime string
	gitCommit string
)

var configFile = "/etc/magicCodeAgent/config.json"

func main() {
	flag.StringVar(&configFile, "config", configFile, "config file path")
	flag.Parse()

	fmt.Printf("magicCodeDebug V1.0\n")
	fmt.Printf("buildTime:%s, gitCommit:%s\n", buildTime, gitCommit)
	config.SetEndpoint("magicCodeDebug")
	config.EnableDebug()

	if configFile != "" {
		configErr := config.LoadConfig(configFile)
		if configErr != nil {
			log.Errorf("load config file failed, error:%s", configErr.Error())
			return
		}
	}
	err := application.Startup(service.DefaultService(config.EndpointName()))
	if err != nil {
		log.Errorf("application.Startup err:%s", err.Error())
		return
	}

	err = application.Run()
	if err != nil {
		log.Errorf("application.Run err:%s", err.Error())
		return
	}

	application.Shutdown()
}
