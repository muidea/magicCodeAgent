package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/muidea/magicCommon/foundation/log"
	fn "github.com/muidea/magicCommon/foundation/net"
	fu "github.com/muidea/magicCommon/foundation/util"

	"github.com/muidea/magicEngine/sse"

	"github.com/muidea/magicCodeAgent/internal/modules/kernel/llm/config"
	imklc "github.com/muidea/magicCodeAgent/internal/modules/kernel/llm/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
)

type StreamTunnel struct {
	sseClient            *sse.Client
	context              context.Context
	header               url.Values
	contextStreamSinker  sse.StreamSinker
	contextContentSinker common.ContentSinker

	localContentSinker common.ContentSinker
	packedContent      string
}

func (s *StreamTunnel) OnClose() {
	s.sseClient.Close()
	if s.contextStreamSinker != nil {
		s.contextStreamSinker.OnClose()
	}
	if s.contextContentSinker != nil {
		s.contextContentSinker.OnClose()
	}
	if s.localContentSinker != nil {
		s.localContentSinker.OnClose()
	}
}

func (s *StreamTunnel) OnRecv(event string, data []byte) {
	if bytes.Equal(data, []byte(" [DONE]\n")) {
		// 这里不需要主动断开连接可以方便后续继续使用
		//s.sseClient.Close()
		return
	}

	if s.contextStreamSinker != nil {
		s.contextStreamSinker.OnRecv(event, data)
	}

	respPtr := &imklc.ChatCompletionStreamResponse{}
	respErr := json.Unmarshal(data, respPtr)
	if respErr != nil {
		log.Errorf("OnRecv failed, error:%v", respErr)
	} else {
		if len(respPtr.Choices) > 0 {
			s.packedContent += respPtr.Choices[0].Delta.Content
		}
		if respPtr.Choices[0].FinishReason != "" {
			if s.contextContentSinker != nil {
				s.contextContentSinker.OnRecv(event, s.packedContent)
			}
			if s.localContentSinker != nil {
				s.localContentSinker.OnRecv(event, s.packedContent)
			}
			s.packedContent = ""
		}
	}
}

func (s *StreamTunnel) Run(value any) (err error) {
	if value == nil {
		err = s.sseClient.Get(s.context, s.header, s)
		return
	}

	byteVal, byteErr := json.Marshal(value)
	if byteErr != nil {
		err = byteErr
		log.Errorf("Run failed, error:%v", err)
		return
	}

	err = s.sseClient.Post(s.context, byteVal, s.header, s)
	if err != nil {
		log.Errorf("Run failed, error:%v", err)
	}
	return
}

func NewStream(
	serverUrl string,
	retryWait time.Duration,
	maxRetries int,
	context context.Context,
	localContentSinker common.ContentSinker,
	header url.Values,
) *StreamTunnel {
	var streamSinker sse.StreamSinker
	var cotentSinker common.ContentSinker
	streamVal := context.Value(sse.StreamSink{})
	streamSinkVal, streamSinkOK := streamVal.(sse.StreamSinker)
	if streamSinkOK {
		streamSinker = streamSinkVal
	}
	contentVal := context.Value(common.ContentSink{})
	cotentSinkVal, cotentSinkOK := contentVal.(common.ContentSinker)
	if cotentSinkOK {
		cotentSinker = cotentSinkVal
	}

	return &StreamTunnel{
		sseClient: sse.NewClient(
			serverUrl,
			retryWait,
			maxRetries,
		),
		context:              context,
		header:               header,
		contextStreamSinker:  streamSinker,
		contextContentSinker: cotentSinker,
		localContentSinker:   localContentSinker,
	}
}

type Provider struct {
	workspaceName string
	messageSink   common.MessageSink
}

func NewProvider(
	workspaceName string,
	messageSink common.MessageSink,
) *Provider {
	return &Provider{
		workspaceName: workspaceName,
		messageSink:   messageSink,
	}
}

func (s *Provider) OnClose() {
	closeMsg := &models.MessageInfo{
		Item:       &models.MessageItem{Role: common.Recorder, Content: "terminate"},
		CreateTime: time.Now().Local().Format(fu.CSTLayoutWithMillisecond),
		Workspace:  s.workspaceName,
	}
	if s.messageSink != nil {
		s.messageSink.OnMessage(closeMsg)
	}
}

func (s *Provider) OnRecv(event string, content string) {
	recvMsg := &models.MessageInfo{
		Item:       &models.MessageItem{Role: common.Assistant, Content: content},
		CreateTime: time.Now().Local().Format(fu.CSTLayoutWithMillisecond),
		Workspace:  s.workspaceName,
	}
	if s.messageSink != nil {
		s.messageSink.OnMessage(recvMsg)
	}
}

// ChatCompletion 调用DeepSeek聊天补全API
// 参数:
//
//	paramPtr: 聊天补全请求参数指针
//
// 返回值:
//
//	ret: 聊天补全响应结果指针
//	err: 错误信息，调用成功时为nil
//
// 注意:
//
//	需要配置有效的DeepSeek API Key才能正常调用
func (s *Provider) ChatCompletion(
	workspaceName string,
	chatContext context.Context,
	paramPtr *imklc.ChatCompletionStreamRequest) (ret *imklc.ChatCompletionStreamResponse, err error) {
	providerType := imklc.RemoteDeepseek
	if val := chatContext.Value(imklc.DeepseekProvider{}); val != nil {
		providerType = val.(string)
	}
	resultPtr := &imklc.ChatCompletionStreamResponse{}
	clnt := http.DefaultClient
	urlVal, _ := url.ParseRequestURI(config.GetDeepSeekUrl(providerType))
	urlVal.Path = path.Join(urlVal.Path, imklc.ChatCompletion)
	header := url.Values{}
	header.Set("Content-Type", "application/json;charset=utf-8")
	header.Set("Accept", "application/json;charset=utf-8")
	header.Set("Authorization", fmt.Sprintf("Bearer %s", config.GetDeepSeekAPIKey(providerType)))
	_, err = fn.HTTPPost(clnt, urlVal.String(), paramPtr, resultPtr, header)
	if err != nil {
		log.Errorf("ChatCompletion failed, url:%v, error:%v", urlVal.String(), err)
		return
	}

	ret = resultPtr
	return
}

func (s *Provider) ChatCompletionStream(
	workspaceName string,
	chatContext context.Context,
	paramPtr *imklc.ChatCompletionStreamRequest) error {
	providerType := imklc.RemoteDeepseek
	if val := chatContext.Value(imklc.DeepseekProvider{}); val != nil {
		providerType = val.(string)
	}

	urlVal, _ := url.ParseRequestURI(config.GetDeepSeekUrl(providerType))
	urlVal.Path = path.Join(urlVal.Path, imklc.ChatCompletion)

	header := url.Values{}
	header.Set("Content-Type", "application/json;charset=utf-8")
	header.Set("Authorization", fmt.Sprintf("Bearer %s", config.GetDeepSeekAPIKey(providerType)))

	paramPtr.Model = config.GetDeepSeekModelName(providerType)

	byteVal, byteErr := json.Marshal(paramPtr)
	if byteErr == nil && s.messageSink != nil {
		s.messageSink.OnMessage(&models.MessageInfo{
			Item:       &models.MessageItem{Role: common.User, Content: string(byteVal)},
			CreateTime: time.Now().Local().Format(fu.CSTLayoutWithMillisecond),
			Workspace:  workspaceName,
		})
	}

	stream := NewStream(
		urlVal.String(),
		config.GetDeepSeekRetryInterval(),
		config.GetDeepSeekMaxRetries(),
		chatContext,
		s,
		header)
	return stream.Run(paramPtr)
}

// BetaCompletion 调用DeepSeek测试版补全API
// 参数:
//
//	paramPtr: 测试版补全请求参数指针
//
// 返回值:
//
//	ret: 测试版补全响应结果指针
//	err: 错误信息，调用成功时为nil
//
// 注意:
//
//	该API可能不稳定，建议仅用于测试目的
func (s *Provider) BetaCompletion(
	workspaceName string,
	chatContext context.Context,
	paramPtr *imklc.BetaCompletionRequest) (ret *imklc.BetaCompletionStreamResponse, err error) {

	providerType := imklc.RemoteDeepseek
	if val := chatContext.Value(imklc.DeepseekProvider{}); val != nil {
		providerType = val.(string)
	}

	paramPtr.Model = config.GetDeepSeekModelName(providerType)

	resultPtr := &imklc.BetaCompletionStreamResponse{}
	clnt := http.DefaultClient
	urlVal, _ := url.ParseRequestURI(config.GetDeepSeekUrl(providerType))
	urlVal.Path = path.Join(urlVal.Path, imklc.BetaCompletion)
	header := url.Values{}
	header.Set("Content-Type", "application/json;charset=utf-8")
	header.Set("Accept", "application/json;charset=utf-8")
	header.Set("Authorization", fmt.Sprintf("Bearer %s", config.GetDeepSeekAPIKey(providerType)))
	_, err = fn.HTTPPost(clnt, urlVal.String(), paramPtr, resultPtr, header)
	if err != nil {
		log.Errorf("BetaCompletion failed, url:%v, error:%v", urlVal.String(), err)
		return
	}

	ret = resultPtr
	return
}

func (s *Provider) BetaCompletionStream(
	workspaceName string,
	chatContext context.Context,
	paramPtr *imklc.BetaCompletionRequest) error {

	providerType := imklc.RemoteDeepseek
	if val := chatContext.Value(imklc.DeepseekProvider{}); val != nil {
		providerType = val.(string)
	}

	urlVal, _ := url.ParseRequestURI(config.GetDeepSeekUrl(providerType))
	urlVal.Path = path.Join(urlVal.Path, imklc.BetaCompletion)
	header := url.Values{}
	header.Set("Content-Type", "application/json;charset=utf-8")
	header.Set("Authorization", fmt.Sprintf("Bearer %s", config.GetDeepSeekAPIKey(providerType)))

	byteVal, byteErr := json.Marshal(paramPtr)
	if byteErr == nil && s.messageSink != nil {
		s.messageSink.OnMessage(&models.MessageInfo{
			Item:       &models.MessageItem{Role: common.User, Content: string(byteVal)},
			CreateTime: time.Now().Local().Format(fu.CSTLayoutWithMillisecond),
			Workspace:  workspaceName,
		})
	}

	stream := NewStream(urlVal.String(),
		config.GetDeepSeekRetryInterval(),
		config.GetDeepSeekMaxRetries(),
		chatContext,
		s,
		header)

	return stream.Run(paramPtr)
}

// ModelsList 获取DeepSeek可用模型列表
// 返回值:
//
//	ret: 模型列表响应结果指针
//	err: 错误信息，调用成功时为nil
//
// 注意:
//
//	需要配置有效的DeepSeek API Key才能正常调用
func (s *Provider) ModelsList() (ret *imklc.ModelsListResponse, err error) {
	resultPtr := &imklc.ModelsListResponse{}
	clnt := http.DefaultClient
	urlVal, _ := url.ParseRequestURI(config.GetDeepSeekUrl(imklc.RemoteDeepseek))
	urlVal.Path = path.Join(urlVal.Path, imklc.ModelsList)
	ctx := url.Values{}
	ctx.Set("Content-Type", "application/json;charset=utf-8")
	ctx.Set("Accept", "application/json;charset=utf-8")
	ctx.Set("Authorization", fmt.Sprintf("Bearer %s", config.GetDeepSeekAPIKey(imklc.RemoteDeepseek)))
	_, err = fn.HTTPGet(clnt, urlVal.String(), resultPtr, ctx)
	if err != nil {
		log.Errorf("ModelsList failed, url:%v, error:%v", urlVal.String(), err)
		return
	}
	ret = resultPtr
	return
}

// UserBalance 查询DeepSeek账户余额
// 返回值:
//
//	ret: 账户余额响应结果指针
//	err: 错误信息，调用成功时为nil
//
// 注意:
//
//	该接口不需要API Key认证
func (s *Provider) UserBalance() (ret *imklc.UserBalanceResponse, err error) {
	resultPtr := &imklc.UserBalanceResponse{}
	clnt := http.DefaultClient
	urlVal, _ := url.ParseRequestURI(config.GetDeepSeekUrl(imklc.RemoteDeepseek))
	urlVal.Path = path.Join(urlVal.Path, imklc.UserBalance)
	ctx := url.Values{}
	ctx.Set("Content-Type", "application/json;charset=utf-8")
	ctx.Set("Accept", "application/json;charset=utf-8")
	ctx.Set("Authorization", fmt.Sprintf("Bearer %s", config.GetDeepSeekAPIKey(imklc.RemoteDeepseek)))
	_, err = fn.HTTPGet(clnt, urlVal.String(), resultPtr, ctx)
	if err != nil {
		log.Errorf("UserBalance failed, url:%v, error:%v", urlVal.String(), err)
		return
	}
	ret = resultPtr
	return
}
