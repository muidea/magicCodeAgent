package main

import (
	"flag"
	"fmt"

	"github.com/muidea/magicCommon/foundation/log"
	"github.com/muidea/magicCommon/framework/application"
	"github.com/muidea/magicCommon/framework/service"

	_ "github.com/muidea/magicCodeAgent"

	"github.com/muidea/magicCodeAgent/internal/config"

	_ "github.com/muidea/magicCodeAgent/internal/initiators/embedtext"
	_ "github.com/muidea/magicCodeAgent/internal/initiators/pprof"
	_ "github.com/muidea/magicCodeAgent/internal/initiators/prompts"
	_ "github.com/muidea/magicCodeAgent/internal/initiators/routeregistry"
	_ "github.com/muidea/magicCodeAgent/internal/initiators/singlechecker"
	_ "github.com/muidea/magicCodeAgent/internal/modules/blocks/commander"
	_ "github.com/muidea/magicCodeAgent/internal/modules/blocks/initiator"
	_ "github.com/muidea/magicCodeAgent/internal/modules/blocks/module"
	_ "github.com/muidea/magicCodeAgent/internal/modules/blocks/project"
	_ "github.com/muidea/magicCodeAgent/internal/modules/blocks/recorder"
	_ "github.com/muidea/magicCodeAgent/internal/modules/blocks/templatemanager"
	_ "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools"
	_ "github.com/muidea/magicCodeAgent/internal/modules/blocks/tracer"
	_ "github.com/muidea/magicCodeAgent/internal/modules/kernel/knowledgebase"
	_ "github.com/muidea/magicCodeAgent/internal/modules/kernel/llm"
	_ "github.com/muidea/magicCodeAgent/internal/modules/kernel/mcp"
	_ "github.com/muidea/magicCodeAgent/internal/modules/kernel/workspace"
)

var (
	buildTime   string
	gitCommit   string
	version     string
	enableTrace string
)

var configFile = ""

func main() {
	flag.StringVar(&configFile, "config", configFile, "config file path")
	flag.Parse()

	fmt.Printf("magicCodeAgent %s\n", version)
	fmt.Printf("buildTime:%s, gitCommit:%s\n", buildTime, gitCommit)
	config.SetApplicationInfo("magicCodeAgent", buildTime, gitCommit, version, enableTrace == "true")

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
