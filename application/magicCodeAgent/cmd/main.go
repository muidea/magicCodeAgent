package main

import (
	"flag"
	"fmt"

	"github.com/muidea/magicCommon/foundation/log"
	"github.com/muidea/magicCommon/framework/application"
	"github.com/muidea/magicCommon/framework/service"

	_ "github.com/muidea/magicToolkit"

	"github.com/muidea/magicToolkit/internal/config"

	_ "github.com/muidea/magicToolkit/internal/modules/blocks/initator"
	_ "github.com/muidea/magicToolkit/internal/modules/blocks/module"
	_ "github.com/muidea/magicToolkit/internal/modules/blocks/project"
	_ "github.com/muidea/magicToolkit/internal/modules/kernel/mcp"
	_ "github.com/muidea/magicToolkit/internal/modules/kernel/watcher"
)

var (
	buildTime string
	gitCommit string
)

var configFile = ""

func main() {
	flag.StringVar(&configFile, "Config", configFile, "config file path")
	flag.Parse()

	fmt.Printf("magicCodeAgent V1.0\n")
	fmt.Printf("buildTime:%s, gitCommit:%s\n", buildTime, gitCommit)
	if configFile != "" {
		configErr := config.LoadConfig(configFile)
		if configErr != nil {
			log.Errorf("load config file failed, error:%s", configErr.Error())
			return
		}
	}
	err := application.Startup(service.HoldService(config.EndpointName()))
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
