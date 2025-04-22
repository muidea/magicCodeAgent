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

	packedContent string
}

func (s *Stream) OnRecv(event string, data []byte) error {
	if bytes.Equal(data, []byte(" [DONE]\n")) {
		// 这里不需要主动断开连接可以方便后续继续使用
		//s.sseClient.Close()
		return nil
	}

	respPtr := &common.ChatCompletionStreamResponse{}
	respErr := json.Unmarshal(data, respPtr)
	if respErr != nil {
		log.Errorf("OnRecv failed, error:%v", respErr)
	} else {
		if len(respPtr.Choices) > 0 {
			s.packedContent += respPtr.Choices[0].Delta.Content
		}
		if respPtr.Choices[0].FinishReason != "" {
			s.streamTracer.AppendMessage(&common.MessageInfo{
				Type:          common.RecvMsg,
				Content:       s.packedContent,
				CreateTime:    time.Now().Local().Format(fu.CSTLayoutWithMillisecond),
				ModelProvider: common.Deepseek,
			})
			s.packedContent = ""
		}
	}

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

	streamTracer.AppendMessage(&common.MessageInfo{
		Type:          common.UrlMsg,
		Content:       serverUrl,
		CreateTime:    time.Now().Local().Format(fu.CSTLayoutWithMillisecond),
		ModelProvider: common.Deepseek,
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
	defer func() {
		if err != nil {
			s.streamTracer.AppendMessage(&common.MessageInfo{
				Type:          common.ErrorMsg,
				Content:       err.Error(),
				CreateTime:    time.Now().Local().Format(fu.CSTLayoutWithMillisecond),
				ModelProvider: common.Deepseek,
			})
		}
	}()

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

	s.streamTracer.AppendMessage(&common.MessageInfo{
		Type:          common.SendMsg,
		Content:       string(byteVal),
		CreateTime:    time.Now().Local().Format(fu.CSTLayoutWithMillisecond),
		ModelProvider: common.Deepseek,
	})

	err = s.sseClient.Post(s.context, byteVal, s.header, s)
	return
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
func ChatCompletion(paramPtr *common.ChatCompletionRequest) (ret *common.ChatCompletionStreamResponse, err error) {
	resultPtr := &common.ChatCompletionStreamResponse{}
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
func BetaCompletion(paramPtr *common.BetaCompletionRequest) (ret *common.BetaCompletionStreamResponse, err error) {
	resultPtr := &common.BetaCompletionStreamResponse{}
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
