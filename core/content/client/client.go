package client

import (
	"github.com/muidea/magicToolkit/core/util"
	"github.com/muidea/magicToolkit/model"
)

func init() {
	util.RegisterFunc("clientServiceFuncDeclare", ServiceFuncDeclare)
	util.RegisterFunc("clientServiceFuncDefine", ServiceFuncDefine)
}

func ServiceFuncDeclare(servicePtr *model.Module) string {
	return "hello world"
	//return fmt.Sprintf("service:%s", servicePtr.Name)
}

func ServiceFuncDefine(servicePtr *model.Module) string {
	return "hello world"
	//return fmt.Sprintf("service:%s", servicePtr.Name)
}
