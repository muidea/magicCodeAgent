package biz

import (
	"encoding/json"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicCommon/foundation/log"

	"github.com/muidea/magicCodeAgent/internal/config"
	"github.com/muidea/magicCodeAgent/internal/modules/base/biz"
	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/templatemanager/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
)

type Templatemanager struct {
	biz.Base

	innerFS fs.FS
	localFS fs.FS
}

func New(
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine,
) *Templatemanager {
	ptr := &Templatemanager{
		Base:    biz.New(imbtc.TemplatemanagerModule, eventHub, backgroundRoutine),
		innerFS: config.GetTemplateFS(),
		localFS: os.DirFS(config.RootConfigDir()),
	}

	ptr.SubscribeFunc(imbtc.ListTemplate, ptr.listTemplate)
	ptr.SubscribeFunc(imbtc.LoadTemplate, ptr.loadTemplate)

	return ptr
}

func (s *Templatemanager) Run() (err *cd.Error) {
	return
}

func (s *Templatemanager) listTemplate(ev event.Event, re event.Result) {
	templateList := s.ListTemplate()
	if re != nil {
		re.Set(templateList, nil)
	}
}

func (s *Templatemanager) loadTemplate(ev event.Event, re event.Result) {
	templateVal := ev.Data()
	if templateVal == nil {
		if re != nil {
			re.Set(nil, cd.NewError(cd.IllegalParam, "invalid template name"))
		}
		return
	}
	templateNameVal, templateNameOK := templateVal.(string)
	if !templateNameOK {
		if re != nil {
			re.Set(nil, cd.NewError(cd.IllegalParam, "invalid template name"))
		}
		return
	}
	templatePtr, templateErr := s.LoadTemplate(templateNameVal)
	if re != nil {
		re.Set(templatePtr, templateErr)
	}
}

func (s *Templatemanager) ListTemplate() (ret []string) {
	var templateList []string
	entryList, entryErr := fs.ReadDir(s.localFS, models.TemplatePath)
	if entryErr == nil {
		for _, entry := range entryList {
			if !entry.IsDir() {
				continue
			}

			subEntryList, subEntryErr := fs.ReadDir(s.localFS, path.Join(models.TemplatePath, entry.Name()))
			if subEntryErr != nil {
				continue
			}

			for _, subEntry := range subEntryList {
				if !subEntry.IsDir() {
					continue
				}
				templateList = append(templateList, path.Join(models.LocalKey, entry.Name(), subEntry.Name()))
			}
		}
	}

	entryList, entryErr = fs.ReadDir(s.innerFS, models.TemplatePath)
	if entryErr == nil {
		for _, entry := range entryList {
			if !entry.IsDir() {
				continue
			}

			subEntryList, subEntryErr := fs.ReadDir(s.innerFS, path.Join(models.TemplatePath, entry.Name()))
			if subEntryErr != nil {
				continue
			}
			for _, subEntry := range subEntryList {
				if !subEntry.IsDir() {
					continue
				}
				templateList = append(templateList, path.Join(models.InnerKey, entry.Name(), subEntry.Name()))
			}
		}
	}

	ret = templateList
	return
}

func (s *Templatemanager) LoadTemplate(templateName string) (ret *models.Template, err *cd.Error) {
	templateFS := s.localFS
	items := strings.Split(templateName, "/")
	if len(items) == 3 {
		templateName = strings.Join(items[1:], "/")
		if items[0] == models.InnerKey {
			templateFS = s.innerFS
		} else {
			templateFS = s.localFS
		}
	}

	templateFinalName := filepath.Join(models.TemplatePath, templateName, "template.json")
	byteVal, byteErr := fs.ReadFile(templateFS, templateFinalName)
	if byteErr != nil {
		log.Errorf("load template failed, finalName:%s, error:%s", templateFinalName, byteErr.Error())
		err = cd.NewError(cd.Unexpected, "load template failed")
		return
	}

	templatePtr := &models.Template{}
	jsonErr := json.Unmarshal(byteVal, templatePtr)
	if jsonErr != nil {
		log.Errorf("load template failed, finalName:%s, error:%s", templateFinalName, jsonErr.Error())
		err = cd.NewError(cd.Unexpected, "load template failed")
		return
	}
	templatePtr.HostRootFS = templateFS
	if validErr := templatePtr.Validate(); validErr != nil {
		log.Errorf("valid template failed, finalName:%s, error:%s", templateFinalName, validErr.Error())
		err = cd.NewError(cd.Unexpected, "valid template failed")
		return
	}

	ret = templatePtr
	return
}
