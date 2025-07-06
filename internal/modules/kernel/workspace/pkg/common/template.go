package common

import (
	cd "github.com/muidea/magicCommon/def"
	fu "github.com/muidea/magicCommon/foundation/util"
)

const WorkspaceTemplateBasePath = "assets/template/workspace"

type Template struct {
	Name        string       `json:"name"`        // 工作区名称
	Description string       `json:"description"` // 工作区描述
	TechStack   fu.StringSet `json:"techStack"`   // 技术栈列表
	Ignores     fu.StringSet `json:"ignores"`     // 忽略文件/目录列表
}

type WorkspaceTemplateResult struct {
	cd.Result
	Value *Template `json:"value"`
}
