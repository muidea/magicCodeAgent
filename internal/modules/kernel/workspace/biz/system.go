package biz

import (
	cd "github.com/muidea/magicCommon/def"

	"github.com/muidea/magicCodeAgent/internal/config"
	imklcfg "github.com/muidea/magicCodeAgent/internal/modules/kernel/llm/config"
	imklc "github.com/muidea/magicCodeAgent/internal/modules/kernel/llm/pkg/common"
	imkwc "github.com/muidea/magicCodeAgent/internal/modules/kernel/workspace/pkg/common"
)

var apiKeyMask = "sk-***"

func (s *Workspace) GetSystemConfig() (authPtr *imkwc.AuthorInfo, localModel *imkwc.ModelInfo, remoteModel *imkwc.ModelInfo) {
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
		if val.Name == imklc.LocalDeepseek {
			localModel = &imkwc.ModelInfo{
				APIKey:   "sk-***",
				BaseUrl:  val.BaseUrl,
				Model:    val.Model,
				Provider: val.Provider,
			}
		}
	}

	return
}

func (s *Workspace) SetSystemConfig(paramPtr *imkwc.SystemConfigParam) (err *cd.Error) {
	authorConfig := config.AuthorConfig()
	if authorConfig == nil {
		authorConfig = &config.AuthorConfigInfo{}
	}
	if paramPtr.AuthorInfo != nil {
		if paramPtr.AuthorInfo.Author != "" {
			authorConfig.Author = paramPtr.AuthorInfo.Author
		}
		if paramPtr.AuthorInfo.EMail != "" {
			authorConfig.EMail = paramPtr.AuthorInfo.EMail
		}
	}

	models := config.Models()
	var curLocal *config.ModelInfo
	var curRemote *config.ModelInfo
	newList := []*config.ModelInfo{}
	for _, val := range models {
		if val.Name == imklc.LocalDeepseek {
			curLocal = val
		}
		if val.Name == imklc.RemoteDeepseek {
			curRemote = val
		}
	}

	if curLocal == nil {
		curLocal = imklcfg.GetDefaultLocalModel()
	}
	if curRemote == nil {
		curRemote = imklcfg.GetDefaultRemoteModel()
	}

	if paramPtr.LocalModel != nil {
		if paramPtr.LocalModel.BaseUrl != "" {
			curLocal.BaseUrl = paramPtr.LocalModel.BaseUrl
		}
		if paramPtr.LocalModel.Model != "" {
			curLocal.Model = paramPtr.LocalModel.Model
		}
		if paramPtr.LocalModel.Provider != "" {
			curLocal.Provider = paramPtr.LocalModel.Provider
		}
		if paramPtr.LocalModel.APIKey != "" && paramPtr.LocalModel.APIKey != apiKeyMask {
			curLocal.APIKey = paramPtr.LocalModel.APIKey
		}
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

	newList = append(newList, curLocal)
	newList = append(newList, curRemote)
	config.UpdateConfig(authorConfig, newList)
	return
}
