package main

import (
	"flag"
	"fmt"

	"github.com/muidea/magicCommon/foundation/log"
	"github.com/muidea/magicCommon/framework/application"
	"github.com/muidea/magicCommon/framework/service"

	_ "github.com/muidea/magicCodeAgent"

	"github.com/muidea/magicCodeAgent/internal/config"

	_ "github.com/muidea/magicCodeAgent/internal/initiators/prompts"
	_ "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools"
)

var (
	buildTime string
	gitCommit string
	version   string
)

func main() {
	flag.Parse()

	fmt.Printf("magicCodeTool %s\n", version)
	fmt.Printf("buildTime:%s, gitCommit:%s\n", buildTime, gitCommit)
	config.SetApplicationInfo("magicCodeTool", buildTime, gitCommit, version)

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
