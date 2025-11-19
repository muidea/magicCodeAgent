package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/muidea/magicCommon/foundation/log"
	"github.com/muidea/magicCommon/framework/application"
	"github.com/muidea/magicCommon/framework/service"

	_ "github.com/muidea/magicCodeAgent"

	"github.com/muidea/magicCodeAgent/internal/config"

	_ "github.com/muidea/magicCodeAgent/internal/initiators/routeregistry"
	_ "github.com/muidea/magicCodeAgent/internal/initiators/singlechecker"
	_ "github.com/muidea/magicCodeAgent/internal/modules/blocks/commander"
	_ "github.com/muidea/magicCodeAgent/internal/modules/blocks/initiator"
	_ "github.com/muidea/magicCodeAgent/internal/modules/blocks/module"
	_ "github.com/muidea/magicCodeAgent/internal/modules/blocks/project"
	_ "github.com/muidea/magicCodeAgent/internal/modules/blocks/templatemanager"
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

	// 这里是为了避免与magicCodeAgent的端口冲突，并没有其他的作用
	bindPort := config.BindPort()
	os.Setenv("BindPort", "2"+bindPort)

	config.SetApplicationInfo("magicCodeCli", buildTime, gitCommit, version, false)
	err := application.Startup(service.DefaultService())
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
