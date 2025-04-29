package common

import (
	openai "github.com/sashabaranov/go-openai"
)

const AgentModule = "/internal/modules/kernel/agent"

const (
	ChatCompletion = "/chat/completions"
	BetaCompletion = "/beta/completions"
	ModelsList     = "/models"
	UserBalance    = "/user/balance"
)

// ChatCompletionStreamRequest 兼容OpenAI的聊天完成请求结构
type ChatCompletionStreamRequest struct {
	openai.ChatCompletionRequest
}

func (s *ChatCompletionStreamRequest) IsStream() bool {
	return s.Stream
}

type ChatCompletionStreamResponse struct {
	openai.ChatCompletionStreamResponse
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
      "total_balance": "110.00",
      "granted_balance": "10.00",
      "topped_up_balance": "100.00"
    }
  ]
}
*/
type UserBalanceResponse struct {
	IsAvailable   bool          `json:"is_available"`
	Balance_infos BalanceDetail `json:"balance_infos"`
}
