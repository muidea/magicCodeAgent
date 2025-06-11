package main

import (
	"flag"
	"fmt"

	"github.com/muidea/magicCommon/foundation/log"
	"github.com/muidea/magicCommon/framework/application"
	"github.com/muidea/magicCommon/framework/service"

	_ "github.com/muidea/magicCodeAgent"

	"github.com/muidea/magicCodeAgent/internal/config"

	_ "github.com/muidea/magicCodeAgent/internal/modules/blocks/initiator"
	_ "github.com/muidea/magicCodeAgent/internal/modules/blocks/module"
	_ "github.com/muidea/magicCodeAgent/internal/modules/blocks/project"
	_ "github.com/muidea/magicCodeAgent/internal/modules/kernel/cli"
)

var (
	buildTime string
	gitCommit string
	version   string
)

func main() {
	flag.Parse()

	fmt.Printf("magicCodeCli %s\n", version)
	fmt.Printf("buildTime:%s, gitCommit:%s\n", buildTime, gitCommit)
	config.SetApplicationInfo("magicCodeCli", buildTime, gitCommit, version)
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
