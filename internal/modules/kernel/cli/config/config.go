package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/foundation/log"

	"github.com/muidea/magicCodeAgent/internal/config"
	"github.com/muidea/magicCodeAgent/pkg/models"
)

const (
	ModuleKey    = "module"
	InitiatorKey = "initiator"
	ProjectKey   = "project"
	CommanderKey = "commander"
	CreateKey    = "create"
	DestroyKey   = "destroy"
	RefreshKey   = "refresh"
	ExecuteKey   = "execute"
	WorkspaceKey = "workspace"
)

// actionType: module,initiator,project
var actionType string

// executeAction: create,destroy,refresh
var executeAction string

var executeCommand string

var currentWorkspace *models.Workspace

func init() {
	flag.StringVar(&actionType, "actionType", "", "action type (required, options: module/initiator/project/commander)")
	flag.StringVar(&executeAction, "executeAction", "", "execute action (required, options: create/destroy/refresh/execute)")
	flag.StringVar(&executeCommand, "executeCommand", "", "execute command(optional, only for executeAction=execute)")
}

func GetActionType() string {
	return actionType
}

func GetExecuteAction() string {
	return executeAction
}

func GetExecuteCommand() string {
	return executeCommand
}

func GetCurrentWorkspace() *models.Workspace {
	return currentWorkspace
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
	if actionType != ModuleKey && actionType != InitiatorKey && actionType != ProjectKey && actionType != CommanderKey {
		err = cd.NewError(cd.IllegalParam, fmt.Sprintf("actionType:%s, 必须是module/initiator/project/command之一", actionType))
		return
	}

	// 验证executeAction
	if executeAction == "" {
		err = cd.NewError(cd.IllegalParam, "executeAction不能为空")
		return
	}
	if executeAction != CreateKey && executeAction != DestroyKey && executeAction != RefreshKey && executeAction != ExecuteKey {
		err = cd.NewError(cd.IllegalParam, fmt.Sprintf("executeAction:%s, 必须是create/destroy/refresh/execute之一", executeAction))
		return
	}

	workspaceName, workspaceOK := os.LookupEnv(WorkspaceKey)
	if workspaceOK {
		currentWorkspace = config.GetWorkspace(workspaceName)
	} else {
		// 读取inner service创建Workspace
		// 所以这里只能默认为是基于magic框架的golang项目
		currentWorkspace = getDefaultWorkspace()
	}

	if currentWorkspace == nil {
		err = cd.NewError(cd.IllegalParam, "当前工作目录不是有效的项目目录")
		return
	}

	return
}

func getDefaultWorkspace() *models.Workspace {
	workingPath := config.CurrentWorkingPath()

	codeRootPath := workingPath
	workingPath = filepath.ToSlash(workingPath)
	log.Infof("working path:%s", workingPath)
	suffixes := []string{
		"/internal/initiators",
		"/internal/modules/kernel",
		"/internal/modules/blocks",
	}
	for _, suffix := range suffixes {
		if strings.HasSuffix(workingPath, suffix) {
			codeRootPath = strings.TrimSuffix(workingPath, suffix)
			codeRootPath = strings.TrimRight(codeRootPath, "/")
			break
		}
	}

	log.Infof("codeRootPath:%s", codeRootPath)
	goModPath := filepath.Join(codeRootPath, "go.mod")
	if _, err := os.Stat(goModPath); err == nil {
		content, err := os.ReadFile(goModPath)
		if err == nil {
			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "module ") {
					modulePath := strings.TrimSpace(strings.TrimPrefix(line, "module"))
					parts := strings.Split(modulePath, "/")
					name := parts[len(parts)-1]

					return &models.Workspace{
						Name:         name,
						Template:     "inner/service/service_v0.0.1",
						Status:       "enable",
						Stage:        "code",
						WorkspaceDir: workingPath,
					}
				}
			}
		}
	}

	return nil
}
