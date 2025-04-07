package main

import (
	"flag"
	"fmt"

	"github.com/muidea/magicCommon/foundation/log"
	"github.com/muidea/magicCommon/framework/application"
	"github.com/muidea/magicCommon/framework/service"

	_ "github.com/muidea/magicToolkit"

	"github.com/muidea/magicToolkit/internal/config"

	_ "github.com/muidea/magicToolkit/internal/modules/executor"
	_ "github.com/muidea/magicToolkit/internal/modules/initator"
	_ "github.com/muidea/magicToolkit/internal/modules/module"
	_ "github.com/muidea/magicToolkit/internal/modules/project"
)

func main() {
	flag.Parse()

	fmt.Printf("magicToolkit V1.0\n")
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
