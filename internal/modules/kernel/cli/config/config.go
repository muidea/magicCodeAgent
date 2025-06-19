package config

import (
	"flag"
	"fmt"
	"regexp"

	"github.com/muidea/magicCodeAgent/internal/initiators/prompts/pkg/common"
	cd "github.com/muidea/magicCommon/def"
)

const (
	ModuleVal    = "module"
	InitiatorVal = "initiator"
	ProjectVal   = "project"
	CreateVal    = "create"
	DestroyVal   = "destroy"
	RefreshVal   = "refresh"
)

// actionType: module,initiator,project
var actionType string

// executeAction: create,destroy,refresh
var executeAction string
var projectType string
var projectName string
var projectDescription string
var moduleName string
var enableService bool
var enableDao bool

func init() {
	flag.StringVar(&actionType, "actionType", "", "action type (required, options: module/initiator/project)")
	flag.StringVar(&executeAction, "executeAction", "", "execute action (required, options: create/destroy/refresh)")
	flag.StringVar(&projectName, "projectName", "", "project name (required, must start with letter and contain only alphanumeric chars)")
	flag.StringVar(&projectType, "projectType", "", "project type (required, options: web/service)")
	flag.StringVar(&projectDescription, "projectDescription", "", "project description")
	flag.StringVar(&moduleName, "moduleName", "", "module name (required, must start with letter and contain only alphanumeric chars)")
	flag.BoolVar(&enableService, "enableService", false, "enable service (optional, default: false)")
	flag.BoolVar(&enableDao, "enableDao", false, "enable dao (optional, default: false)")
}

func GetActionType() string {
	return actionType
}

func GetExecuteAction() string {
	return executeAction
}

func GetModuleName() string {
	return moduleName
}

func GetProjectName() string {
	return projectName
}

func GetProjectType() string {
	return projectType
}

func GetProjectDescription() string {
	return projectDescription
}

func GetEnableService() bool {
	return enableService
}

func GetEnableDao() bool {
	return enableDao
}

var namePattern = regexp.MustCompile(`^(?:[a-zA-Z][a-zA-Z0-9]*)$`)

// 校验输入参数
// 1. moduleType只能是module,initiator,project其中一种，必选参数，如果参数为空则失败，返回参数错误
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
	if actionType == "" {
		err = cd.NewError(cd.IllegalParam, "actionType不能为空")
		return
	}
	if actionType != "module" && actionType != "initiator" && actionType != "project" {
		err = cd.NewError(cd.IllegalParam, fmt.Sprintf("actionType:%s, 必须是module/initiator/project之一", actionType))
		return
	}

	// 验证executeAction
	if executeAction == "" {
		err = cd.NewError(cd.IllegalParam, "executeAction不能为空")
		return
	}
	if executeAction != "create" && executeAction != "destroy" && executeAction != "refresh" {
		err = cd.NewError(cd.IllegalParam, fmt.Sprintf("executeAction:%s, 必须是create/destroy/refresh之一", executeAction))
		return
	}

	if projectType != common.Service && projectType != common.Web {
		err = cd.NewError(cd.IllegalParam, fmt.Sprintf("projectType:%s, 必须是service/web之一", projectType))
		return
	}

	switch actionType {
	case "module", "initiator":
		if !namePattern.MatchString(moduleName) {
			err = cd.NewError(cd.IllegalParam, fmt.Sprintf("moduleName:%s, 只能是字母开头的普通单词", moduleName))
			return
		}
	case "project":
		if !namePattern.MatchString(projectName) {
			err = cd.NewError(cd.IllegalParam, fmt.Sprintf("projectName:%s, 只能是字母开头的普通单词", projectName))
			return
		}
	}

	return
}
