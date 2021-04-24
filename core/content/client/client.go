package client

import (
	"github.com/muidea/magicToolkit/core/util"
	"github.com/muidea/magicToolkit/model"
)

func init() {
	util.RegisterFunc("serviceDeclare", ServiceDeclare)
}

func ServiceDeclare(servicePtr *model.Module) string {
	return "hello world"
	//return fmt.Sprintf("service:%s", servicePtr.Name)
}
