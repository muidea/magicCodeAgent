package config

import (
	"flag"
	"regexp"

	cd "github.com/muidea/magicCommon/def"
)

const (
	ModuleVal   = "module"
	InitatorVal = "initator"
	ProjectVal  = "project"
	CreateVal   = "create"
	DestroyVal  = "destroy"
	RefreshVal  = "refresh"
)

// moduleType: module,initator,project
var moduleType string

// executeAction: create,destroy,refresh
var executeAction string
var moduleName string
var enableService bool
var enableDao bool

func init() {
	flag.StringVar(&moduleType, "moduleType", "", "module type (required, options: module/initator/project)")
	flag.StringVar(&executeAction, "executeAction", "", "execute action (required, options: create/destroy/refresh)")
	flag.StringVar(&moduleName, "moduleName", "", "module name (required, must start with letter and contain only alphanumeric chars)")
	flag.BoolVar(&enableService, "enableService", false, "enable service (optional, default: false)")
	flag.BoolVar(&enableDao, "enableDao", false, "enable dao (optional, default: false)")
}

func GetModuleType() string {
	return moduleType
}

func GetExecuteAction() string {
	return executeAction
}

func GetModuleName() string {
	return moduleName
}

func GetEnableService() bool {
	return enableService
}

func GetEnableDao() bool {
	return enableDao
}

// 校验输入参数
// 1. moduleType只能是module,initator,project其中一种，必选参数，如果参数为空则失败，返回参数错误
// 2. executeAction只能是create,destroy,refresh其中一种，必选参数，如果参数为空则失败，返回参数错误
// 3. moduleName为必选参数，只能是字母开头的普通单词，不能包含其他字符，不能是数字或其他字符，如果参数为空则失败，返回参数错误
// 4. enableService,enableDao为可选参数，如果参数为空则默认为false，如果参数不为空则判断是否为true或false，如果参数不为true或false则失败，返回参数错误
func VerifyVar() (err *cd.Error) {
	defer func() {
		if err != nil {
			flag.PrintDefaults()
		}
	}()

	// 验证moduleType
	if moduleType == "" {
		err = cd.NewError(cd.IllegalParam, "moduleType不能为空")
		return
	}
	if moduleType != "module" && moduleType != "initator" && moduleType != "project" {
		err = cd.NewError(cd.IllegalParam, "moduleType必须是module/initator/project之一")
		return
	}

	// 验证executeAction
	if executeAction == "" {
		err = cd.NewError(cd.IllegalParam, "executeAction不能为空")
		return
	}
	if executeAction != "create" && executeAction != "destroy" && executeAction != "refresh" {
		err = cd.NewError(cd.IllegalParam, "executeAction必须是create/destroy/refresh之一")
		return
	}

	// 验证moduleName
	if moduleName == "" {
		err = cd.NewError(cd.IllegalParam, "moduleName不能为空")
		return
	}
	if !regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]*$`).MatchString(moduleName) {
		err = cd.NewError(cd.IllegalParam, "moduleName必须以字母开头且只能包含字母数字")
		return
	}

	return
}
