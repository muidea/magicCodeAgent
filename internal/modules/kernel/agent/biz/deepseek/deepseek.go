package deepseek

import (
	"fmt"
	"net/http"
	"net/url"
	"path"

	fn "github.com/muidea/magicCommon/foundation/net"

	"github.com/muidea/magicCodeAgent/internal/modules/kernel/agent/config"
	"github.com/muidea/magicCodeAgent/internal/modules/kernel/agent/pkg/common"
)

func ChatCompletion(paramPtr *common.ChatCompletionRequest) (ret *common.ChatCompletionResponse, err error) {
	resultPtr := &common.ChatCompletionResponse{}
	clnt := http.DefaultClient
	fullUrl := path.Join(config.GetDeepSeekUrl(), common.ChatCompletion)
	ctx := url.Values{}
	ctx.Set("Content-Type", "application/json;charset=utf-8")
	ctx.Set("Accept", "application/json;charset=utf-8")
	ctx.Set("Authorization", fmt.Sprintf("Bearer %s", config.GetDeepSeekAPIKey()))
	_, err = fn.HTTPPost(clnt, fullUrl, paramPtr, resultPtr, ctx)
	if err != nil {
		return
	}

	ret = resultPtr
	return
}

func BetaCompletion(paramPtr *common.BetaCompletionRequest) (ret *common.BetaCompletionResponse, err error) {
	resultPtr := &common.BetaCompletionResponse{}
	clnt := http.DefaultClient
	fullUrl := path.Join(config.GetDeepSeekUrl(), common.BetaCompletion)
	ctx := url.Values{}
	ctx.Set("Content-Type", "application/json;charset=utf-8")
	ctx.Set("Accept", "application/json;charset=utf-8")
	ctx.Set("Authorization", fmt.Sprintf("Bearer %s", config.GetDeepSeekAPIKey()))
	_, err = fn.HTTPPost(clnt, fullUrl, paramPtr, resultPtr, ctx)
	if err != nil {
		return
	}

	ret = resultPtr
	return
}

func ModelsList() (ret *common.ModelsListResponse, err error) {
	resultPtr := &common.ModelsListResponse{}
	clnt := http.DefaultClient
	fullUrl := path.Join(config.GetDeepSeekUrl(), common.ModelsList)
	ctx := url.Values{}
	ctx.Set("Content-Type", "application/json;charset=utf-8")
	ctx.Set("Accept", "application/json;charset=utf-8")
	ctx.Set("Authorization", fmt.Sprintf("Bearer %s", config.GetDeepSeekAPIKey()))
	_, err = fn.HTTPGet(clnt, fullUrl, resultPtr, ctx)
	if err != nil {
		return
	}
	ret = resultPtr
	return
}

func UserBalance() (ret *common.UserBalanceResponse, err error) {
	resultPtr := &common.UserBalanceResponse{}
	clnt := http.DefaultClient
	fullUrl := path.Join(config.GetDeepSeekUrl(), common.UserBalance)
	ctx := url.Values{}
	ctx.Set("Content-Type", "application/json;charset=utf-8")
	ctx.Set("Accept", "application/json;charset=utf-8")
	_, err = fn.HTTPGet(clnt, fullUrl, resultPtr, ctx)
	if err != nil {
		return
	}
	ret = resultPtr
	return
}
