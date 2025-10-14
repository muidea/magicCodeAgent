package biz

import (
	cd "github.com/muidea/magicCommon/def"

	"github.com/muidea/magicCodeAgent/internal/config"
	imklcfg "github.com/muidea/magicCodeAgent/internal/modules/kernel/llm/config"
	imklc "github.com/muidea/magicCodeAgent/internal/modules/kernel/llm/pkg/common"
	imkwc "github.com/muidea/magicCodeAgent/internal/modules/kernel/workspace/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
)

var apiKeyMask = "sk-******"

func (s *Workspace) GetVersionInfo() *imkwc.VersionInfo {
	return &imkwc.VersionInfo{
		BuildTime:    config.BuildTime(),
		EndpointName: config.EndpointName(),
		GitCommit:    config.GitCommit(),
		Version:      config.Version(),
	}
}

func (s *Workspace) GetSystemConfig() (authPtr *imkwc.AuthorInfo, remoteModel *imkwc.ModelInfo, workspaceSetting *models.WorkspaceSetting) {
	authorConfig := config.AuthorConfig()
	if authorConfig != nil {
		authPtr = &imkwc.AuthorInfo{
			Author: authorConfig.Author,
			EMail:  authorConfig.EMail,
		}
	}

	models := config.Models()
	for _, val := range models {
		if val.Name == imklc.RemoteDeepseek {
			remoteModel = &imkwc.ModelInfo{
				APIKey:   "sk-***",
				BaseUrl:  val.BaseUrl,
				Model:    val.Model,
				Provider: val.Provider,
			}
		}
	}

	workspaceSetting = config.WorkspaceSetting()

	return
}

func (s *Workspace) SetSystemConfig(paramPtr *imkwc.SystemConfigParam) (err *cd.Error) {
	authorConfig := config.AuthorConfig()
	if authorConfig == nil {
		authorConfig = &models.AuthorInfo{}
	}
	if paramPtr.AuthorInfo != nil {
		if paramPtr.AuthorInfo.Author != "" {
			authorConfig.Author = paramPtr.AuthorInfo.Author
		}
		if paramPtr.AuthorInfo.EMail != "" {
			authorConfig.EMail = paramPtr.AuthorInfo.EMail
		}
	}

	modelList := config.Models()
	var curRemote *models.ModelInfo
	newList := []*models.ModelInfo{}
	for _, val := range modelList {
		if val.Name == imklc.RemoteDeepseek {
			curRemote = val
		}
	}

	if curRemote == nil {
		curRemote = imklcfg.GetDefaultRemoteModel()
	}

	if paramPtr.RemoteModel != nil {
		if paramPtr.RemoteModel.BaseUrl != "" {
			curRemote.BaseUrl = paramPtr.RemoteModel.BaseUrl
		}
		if paramPtr.RemoteModel.Model != "" {
			curRemote.Model = paramPtr.RemoteModel.Model
		}
		if paramPtr.RemoteModel.Provider != "" {
			curRemote.Provider = paramPtr.RemoteModel.Provider
		}
		if paramPtr.RemoteModel.APIKey != "" && paramPtr.RemoteModel.APIKey != apiKeyMask {
			curRemote.APIKey = paramPtr.RemoteModel.APIKey
		}
	}

	newList = append(newList, curRemote)
	config.UpdateConfig(authorConfig, newList, paramPtr.WorkspaceSetting)
	return
}
