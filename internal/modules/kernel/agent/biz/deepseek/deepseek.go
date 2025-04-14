package deepseek

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

	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tracer/pkg/common"
	"github.com/muidea/magicCodeAgent/internal/modules/kernel/agent/config"
	"github.com/muidea/magicCodeAgent/internal/modules/kernel/agent/pkg/common"
)

type Stream struct {
	sseClient    *sse.Client
	context      context.Context
	header       url.Values
	streamSink   sse.Sinker
	streamTracer imbtc.Tracer
}

func (s *Stream) OnRecv(event string, data []byte) error {
	if bytes.Equal(data, []byte(" [DONE]\n")) {
		// 这里不需要主动断开连接可以方便后续继续使用
		//s.sseClient.Close()
		return nil
	}

	s.streamTracer.AppendMessage(&imbtc.MessageInfo{
		Type:       common.RecvMsg,
		Content:    string(data),
		CreateTime: time.Now().UTC().Format(fu.CSTLayoutWithMillisecond),
		Owner:      common.Deepseek,
	})

	s.streamSink.OnRecv(event, data)
	return nil
}

func (s *Stream) OnClose() {
	s.sseClient.Close()
}

func NewStream(serverUrl string,
	context context.Context,
	header url.Values,
	streamSink sse.Sinker,
	streamTracer imbtc.Tracer,
) *Stream {

	streamTracer.AppendMessage(&imbtc.MessageInfo{
		Type:       common.UrlMsg,
		Content:    serverUrl,
		CreateTime: time.Now().UTC().Format(fu.CSTLayoutWithMillisecond),
		Owner:      common.Deepseek,
	})

	return &Stream{
		sseClient: sse.NewClient(
			serverUrl,
			config.GetDeepSeekRetryInterval(),
			config.GetDeepSeekMaxRetries(),
		),
		context:      context,
		header:       header,
		streamSink:   streamSink,
		streamTracer: streamTracer,
	}
}

func (s *Stream) Run(value any) (err error) {
	if value == nil {
		return s.sseClient.Get(s.context, s.header, s)
	}

	byteVal, byteErr := json.Marshal(value)
	if byteErr != nil {
		return
	}

	s.streamTracer.AppendMessage(&imbtc.MessageInfo{
		Type:       common.SendMsg,
		Content:    string(byteVal),
		CreateTime: time.Now().UTC().Format(fu.CSTLayoutWithMillisecond),
		Owner:      common.Deepseek,
	})

	return s.sseClient.Post(s.context, byteVal, s.header, s)
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
func ChatCompletion(paramPtr *common.ChatCompletionRequest) (ret *common.ChatCompletionResponse, err error) {
	resultPtr := &common.ChatCompletionResponse{}
	clnt := http.DefaultClient
	urlVal, _ := url.ParseRequestURI(config.GetDeepSeekUrl())
	urlVal.Path = path.Join(urlVal.Path, common.ChatCompletion)
	ctx := url.Values{}
	ctx.Set("Content-Type", "application/json;charset=utf-8")
	ctx.Set("Accept", "application/json;charset=utf-8")
	ctx.Set("Authorization", fmt.Sprintf("Bearer %s", config.GetDeepSeekAPIKey()))
	_, err = fn.HTTPPost(clnt, urlVal.String(), paramPtr, resultPtr, ctx)
	if err != nil {
		log.Errorf("ChatCompletion failed, url:%v, error:%v", urlVal.String(), err)
		return
	}

	ret = resultPtr
	return
}

func ChatCompletionStream(ctx context.Context, paramPtr *common.ChatCompletionRequest, sink sse.Sinker, tracer imbtc.Tracer) error {
	urlVal, _ := url.ParseRequestURI(config.GetDeepSeekUrl())
	urlVal.Path = path.Join(urlVal.Path, common.ChatCompletion)
	header := url.Values{}
	header.Set("Content-Type", "application/json;charset=utf-8")
	header.Set("Authorization", fmt.Sprintf("Bearer %s", config.GetDeepSeekAPIKey()))

	stream := NewStream(urlVal.String(), ctx, header, sink, tracer)
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
func BetaCompletion(paramPtr *common.BetaCompletionRequest) (ret *common.BetaCompletionResponse, err error) {
	resultPtr := &common.BetaCompletionResponse{}
	clnt := http.DefaultClient
	urlVal, _ := url.ParseRequestURI(config.GetDeepSeekUrl())
	urlVal.Path = path.Join(urlVal.Path, common.BetaCompletion)
	ctx := url.Values{}
	ctx.Set("Content-Type", "application/json;charset=utf-8")
	ctx.Set("Accept", "application/json;charset=utf-8")
	ctx.Set("Authorization", fmt.Sprintf("Bearer %s", config.GetDeepSeekAPIKey()))
	_, err = fn.HTTPPost(clnt, urlVal.String(), paramPtr, resultPtr, ctx)
	if err != nil {
		log.Errorf("BetaCompletion failed, url:%v, error:%v", urlVal.String(), err)
		return
	}

	ret = resultPtr
	return
}

func BetaCompletionStream(ctx context.Context, paramPtr *common.BetaCompletionRequest, sink sse.Sinker, tracer imbtc.Tracer) error {
	urlVal, _ := url.ParseRequestURI(config.GetDeepSeekUrl())
	urlVal.Path = path.Join(urlVal.Path, common.BetaCompletion)
	header := url.Values{}
	header.Set("Content-Type", "application/json;charset=utf-8")
	header.Set("Authorization", fmt.Sprintf("Bearer %s", config.GetDeepSeekAPIKey()))

	stream := NewStream(urlVal.String(), ctx, header, sink, tracer)
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
func ModelsList() (ret *common.ModelsListResponse, err error) {
	resultPtr := &common.ModelsListResponse{}
	clnt := http.DefaultClient
	urlVal, _ := url.ParseRequestURI(config.GetDeepSeekUrl())
	urlVal.Path = path.Join(urlVal.Path, common.ModelsList)
	ctx := url.Values{}
	ctx.Set("Content-Type", "application/json;charset=utf-8")
	ctx.Set("Accept", "application/json;charset=utf-8")
	ctx.Set("Authorization", fmt.Sprintf("Bearer %s", config.GetDeepSeekAPIKey()))
	_, err = fn.HTTPGet(clnt, urlVal.String(), resultPtr, ctx)
	if err != nil {
		log.Errorf("ModelsList failed, url:%v, error:%v", urlVal.String(), err)
		return
	}
	ret = resultPtr
	return
}

func ModelsListStream(ctx context.Context, sink sse.Sinker, tracer imbtc.Tracer) error {
	urlVal, _ := url.ParseRequestURI(config.GetDeepSeekUrl())
	urlVal.Path = path.Join(urlVal.Path, common.ModelsList)
	header := url.Values{}
	header.Set("Content-Type", "application/json;charset=utf-8")
	header.Set("Accept", "text/event-stream")
	header.Set("Authorization", fmt.Sprintf("Bearer %s", config.GetDeepSeekAPIKey()))

	stream := NewStream(urlVal.String(), ctx, header, sink, tracer)
	return stream.Run(nil)
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
func UserBalance() (ret *common.UserBalanceResponse, err error) {
	resultPtr := &common.UserBalanceResponse{}
	clnt := http.DefaultClient
	urlVal, _ := url.ParseRequestURI(config.GetDeepSeekUrl())
	urlVal.Path = path.Join(urlVal.Path, common.UserBalance)
	//fullUrl := path.Join(config.GetDeepSeekUrl(), common.UserBalance)
	ctx := url.Values{}
	ctx.Set("Content-Type", "application/json;charset=utf-8")
	ctx.Set("Accept", "application/json;charset=utf-8")
	ctx.Set("Authorization", fmt.Sprintf("Bearer %s", config.GetDeepSeekAPIKey()))
	_, err = fn.HTTPGet(clnt, urlVal.String(), resultPtr, ctx)
	if err != nil {
		log.Errorf("UserBalance failed, url:%v, error:%v", urlVal.String(), err)
		return
	}
	ret = resultPtr
	return
}

func UserBalanceStream(ctx context.Context, sink sse.Sinker, tracer imbtc.Tracer) error {
	urlVal, _ := url.ParseRequestURI(config.GetDeepSeekUrl())
	urlVal.Path = path.Join(urlVal.Path, common.UserBalance)
	header := url.Values{}
	header.Set("Content-Type", "application/json;charset=utf-8")
	header.Set("Accept", "text/event-stream")
	header.Set("Authorization", fmt.Sprintf("Bearer %s", config.GetDeepSeekAPIKey()))

	stream := NewStream(urlVal.String(), ctx, header, sink, tracer)
	return stream.Run(nil)
}
