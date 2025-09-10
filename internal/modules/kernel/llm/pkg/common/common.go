package common

import (
	"github.com/muidea/magicCodeAgent/pkg/models"
	openai "github.com/sashabaranov/go-openai"
)

const LLMModule = "/internal/modules/kernel/llm"

const (
	ChatCompletionStream   = "/chat/completions"
	ChatCompletionTerminal = "/chat/completions/terminal"
	BetaCompletion         = "/beta/completions"
	ModelsList             = "/models"
	UserBalance            = "/user/balance"
)

type Usage struct {
	CompletionTokens    int `json:"completion_tokens"`
	PromptTokens        int `json:"prompt_tokens"`
	PromptCacheHitToken int `json:"prompt_cache_hit_tokens"`
	PromptMissHitToken  int `json:"prompt_cache_miss_tokens"`
	TotalTokens         int `json:"total_tokens"`
}

func (s *Usage) Update(ptr *Usage) {
	if ptr == nil {
		return
	}

	s.CompletionTokens += ptr.CompletionTokens
	s.PromptTokens += ptr.PromptTokens
	s.PromptCacheHitToken += ptr.PromptCacheHitToken
	s.PromptMissHitToken += ptr.PromptMissHitToken
	s.TotalTokens += ptr.TotalTokens
}

func (s *Usage) Reset() {
	s.CompletionTokens = 0
	s.PromptTokens = 0
	s.PromptCacheHitToken = 0
	s.PromptMissHitToken = 0
	s.TotalTokens = 0
}

func (s *Usage) Copy() *Usage {
	return &Usage{
		CompletionTokens:    s.CompletionTokens,
		PromptTokens:        s.PromptTokens,
		PromptCacheHitToken: s.PromptCacheHitToken,
		PromptMissHitToken:  s.PromptMissHitToken,
		TotalTokens:         s.TotalTokens,
	}
}

type ContentSinker interface {
	OnClose()
	OnRecv(event, content string, usagePtr *Usage)
}

type MessageSink interface {
	OnTerminate(msgPtr *models.MessageInfo)
	OnSendMessage(msgPtr *models.MessageInfo)
	OnRecvMessage(msgPtr *models.MessageInfo, usagePtr *Usage)
}

// ChatCompletionRequest 兼容OpenAI的聊天完成请求结构
type ChatCompletionRequest struct {
	openai.ChatCompletionRequest
}

func (s *ChatCompletionRequest) IsStream() bool {
	return s.Stream
}

type ChatCompletionStreamResponse struct {
	openai.ChatCompletionStreamResponse
	Usage *Usage `json:"usage,omitempty"`
}

type ChatCompletionResponse struct {
	openai.ChatCompletionResponse
	Usage *Usage `json:"usage,omitempty"`
}

// BetaCompletionRequest 测试版完成请求
type BetaCompletionRequest struct {
	openai.ChatCompletionRequest
	// 制定被补全内容的后缀。
	Suffix *string `json:"suffix,omitempty"`
}

func (s *BetaCompletionRequest) IsStream() bool {
	return s.Stream
}

type BetaCompletionStreamResponse struct {
	openai.ChatCompletionStreamResponse
}

// ModelInfo 定义模型基本信息结构
type ModelInfo struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	OwnedBy string `json:"owned_by"`
}

// ModelsListResponse 模型列表响应
/*
{
  "object": "list",
  "data": [
    {
      "id": "string",
      "object": "model",
      "owned_by": "string"
    }
  ]
}
*/
type ModelsListResponse struct {
	Object string      `json:"object"`
	Data   []ModelInfo `json:"data"`
}

// BalanceDetail 定义余额详情结构
type BalanceDetail struct {
	Currency        string `json:"currency"`
	TotalBalance    string `json:"total_balance"`
	GrantedBalance  string `json:"granted_balance"`
	ToppedUpBalance string `json:"topped_up_balance"`
}

// UserBalanceResponse 用户余额响应
/*
{
    "is_available": true,
    "balance_infos": [
        {
            "currency": "CNY",
            "total_balance": "40.05",
            "granted_balance": "0.00",
            "topped_up_balance": "40.05"
        }
    ]
}
*/
type UserBalanceResponse struct {
	IsAvailable   bool            `json:"is_available"`
	Balance_infos []BalanceDetail `json:"balance_infos"`
}
