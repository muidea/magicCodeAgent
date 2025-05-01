package main

import (
	"flag"
	"fmt"

	"github.com/muidea/magicCommon/foundation/log"
	"github.com/muidea/magicCommon/framework/application"
	"github.com/muidea/magicCommon/framework/service"

	_ "github.com/muidea/magicCodeAgent"

	"github.com/muidea/magicCodeAgent/internal/config"

	_ "github.com/muidea/magicCodeAgent/internal/initators/pprof"
	_ "github.com/muidea/magicCodeAgent/internal/initators/prompts"
	_ "github.com/muidea/magicCodeAgent/internal/initators/routeregistry"
	_ "github.com/muidea/magicCodeAgent/internal/modules/blocks/initator"
	_ "github.com/muidea/magicCodeAgent/internal/modules/blocks/module"
	_ "github.com/muidea/magicCodeAgent/internal/modules/blocks/project"
	_ "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools"
	_ "github.com/muidea/magicCodeAgent/internal/modules/blocks/tracer"
	_ "github.com/muidea/magicCodeAgent/internal/modules/kernel/llm"
	_ "github.com/muidea/magicCodeAgent/internal/modules/kernel/mcp"
	_ "github.com/muidea/magicCodeAgent/internal/modules/kernel/workspace"
)

var (
	buildTime string
	gitCommit string
)

var configFile = ""

func main() {
	flag.StringVar(&configFile, "config", configFile, "config file path")
	flag.Parse()

	fmt.Printf("magicCodeAgent V1.0\n")
	fmt.Printf("buildTime:%s, gitCommit:%s\n", buildTime, gitCommit)
	config.SetEndpoint("magicCodeAgent")

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
