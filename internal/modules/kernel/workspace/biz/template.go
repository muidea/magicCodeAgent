package biz

import (
	"encoding/json"
	"io/fs"
	"path"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/foundation/log"

	"github.com/muidea/magicCodeAgent/internal/config"
	imkwc "github.com/muidea/magicCodeAgent/internal/modules/kernel/workspace/pkg/common"
)

func (s *Workspace) QueryWorkspaceTemplate(workspaceType string) (ret *imkwc.Template, err *cd.Error) {
	ret = &imkwc.Template{
		TechStack: []string{},
		Ignores:   []string{},
	}
	fullTemplatePath := path.Join(imkwc.WorkspaceTemplateBasePath, workspaceType, "template.json")
	_, fileErr := fs.Stat(config.GetTemplateFS(), fullTemplatePath)
	if fileErr != nil {
		log.Errorf("stat template.json failed, err:%v", fileErr)
		return
	}

	contentVal, contentErr := fs.ReadFile(config.GetTemplateFS(), fullTemplatePath)
	if contentErr != nil {
		log.Errorf("open template.json failed, err:%v", contentErr)
		return
	}

	if byteErr := json.Unmarshal(contentVal, ret); byteErr != nil {
		log.Errorf("unmarshal template.json failed, err:%v", byteErr)
		return
	}
	return
}
