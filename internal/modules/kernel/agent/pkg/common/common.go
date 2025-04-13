package common

const AgentModule = "/internal/modules/kernel/agent"

const (
	ChatCompletion = "/chat/completions"
	BetaCompletion = "/beta/completions"
	ModelsList     = "/models"
	UserBalance    = "/user/balance"
)

// 以下是参照DeepSeek api定义的各接口的Request和Response定义
// https://api-docs.deepseek.com/zh-cn/api/deepseek-api

// 公共消息结构
type Message struct {
	Role             string  `json:"role"` // "system", "user", "assistant", "tool"
	Content          any     `json:"content"`
	Name             string  `json:"name"`
	Prefix           *bool   `json:"prefix"`            // only for "assistant"
	ReasoningContent *string `json:"reasoning_content"` // only for "assistant"
	ToolCallID       *string `json:"tool_call_id"`      // only for "tool"
}

// 公共选项结构
type Choice struct {
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	Index        int `json:"index"`
	Logprobs     any `json:"logprobs"`
	FinishReason any `json:"finish_reason"`
}

type CompletionTokenDetails struct {
	ReasoningTokens int `json:"reasoning_tokens"`
}

// 公共用量统计结构
type Usage struct {
	CompletionTokens       int                    `json:"completion_tokens"`
	PromptTokens           int                    `json:"prompt_tokens"`
	PromptCacheHitTokens   int                    `json:"prompt_cache_hit_tokens"`
	PromptCacheMissTokens  int                    `json:"prompt_cache_miss_tokens"`
	TotalTokens            int                    `json:"total_tokens"`
	CompletionTokenDetails CompletionTokenDetails `json:"completion_token_details"`
}

/*
一个 object，指定模型必须输出的格式。

设置为 { "type": "json_object" } 以启用 JSON 模式，该模式保证模型生成的消息是有效的 JSON。

注意: 使用 JSON 模式时，你还必须通过系统或用户消息指示模型生成 JSON。否则，模型可能会生成不断的空白字符，直到生成达到令牌限制，从而导致请求长时间运行并显得“卡住”。此外，如果 finish_reason="length"，这表示生成超过了 max_tokens 或对话超过了最大上下文长度，消息内容可能会被部分截断。
*/
type ResponseFormat struct {
	Type string `json:"type"`
}

/*
流式输出相关选项。只有在 stream 参数为 true 时，才可设置此参数。
IncludeUseage: 如果设置为 true，在流式消息最后的 data: [DONE] 之前将会传输一个额外的块。此块上的 usage 字段显示整个请求的 token 使用统计信息，而 choices 字段将始终是一个空数组。所有其他块也将包含一个 usage 字段，但其值为 null。
*/
type StreamOptions struct {
	IncludeUseage bool `json:"include_usage"`
}

type Tool struct{}

type ToolChoice struct {
	ChatCompletionToolChoice      string `json:"chat_completion"`
	ChatCompletionNamedToolChoice string `json:"chat_completion_named"`
}

// ChatCompletionRequest 兼容OpenAI的聊天完成请求结构
type ChatCompletionRequest struct {
	// 对话的消息列表。
	Messages []Message `json:"messages"`
	// 使用的模型的 ID。您可以使用 deepseek-chat。
	Model string `json:"model"` // "deepseek-chat", "deepseek-reasoner"
	// 介于 -2.0 和 2.0 之间的数字。如果该值为正，那么新 token 会根据其在已有文本中的出现频率受到相应的惩罚，降低模型重复相同内容的可能性。
	FrequencyPenalty *float64 `json:"frequency_penalty,omitempty"` // -2<= and <=2 ,default 0
	// 介于 1 到 8192 间的整数，限制一次请求中模型生成 completion 的最大 token 数。输入 token 和输出 token 的总长度受模型的上下文长度的限制。
	MaxTokens *int `json:"max_tokens"` // >1 default 4096
	// 介于 -2.0 和 2.0 之间的数字。如果该值为正，那么新 token 会根据其是否已在已有文本中出现受到相应的惩罚，从而增加模型谈论新主题的可能性。
	PresencePenalty *float64        `json:"presence_penalty,omitempty"` // -2<= and <=2 ,default 0
	ResponseFormat  *ResponseFormat `json:"response_format,omitempty"`
	// 一个 string 或最多包含 16 个 string 的 list，在遇到这些词时，API 将停止生成更多的 token。
	Stop *[]string `json:"stop,omitempty"`
	// 如果设置为 True，将会以 SSE（server-sent events）的形式以流式发送消息增量。消息流以 data: [DONE] 结尾。
	Stream        *bool          `json:"stream"`
	StreamOptions *StreamOptions `json:"stream_options,omitempty"`
	// 采样温度，介于 0 和 2 之间。更高的值，如 0.8，会使输出更随机，而更低的值，如 0.2，会使其更加集中和确定。 我们通常建议可以更改这个值或者更改 top_p，但不建议同时对两者进行修改。
	// default 1
	Temperature *float64 `json:"temperature"`
	// 作为调节采样温度的替代方案，模型会考虑前 top_p 概率的 token 的结果。所以 0.1 就意味着只有包括在最高 10% 概率中的 token 会被考虑。 我们通常建议修改这个值或者更改 temperature，但不建议同时对两者进行修改。
	// <= 1 default 1
	TopP *float64 `json:"top_p"`
	// 模型可能会调用的 tool 的列表。目前，仅支持 function 作为工具。使用此参数来提供以 JSON 作为输入参数的 function 列表。最多支持 128 个 function。
	Tools      *[]Tool     `json:"tools,omitempty"`
	ToolChoice *ToolChoice `json:"tool_choice,omitempty"`
	// 是否返回所输出 token 的对数概率。如果为 true，则在 message 的 content 中返回每个输出 token 的对数概率。
	Logprobs *bool `json:"logprobs,omitempty"`
	// 一个介于 0 到 20 之间的整数 N，指定每个输出位置返回输出概率 top N 的 token，且返回这些 token 的对数概率。指定此参数时，logprobs 必须为 true。
	TopLogprobs *int `json:"top_logprobs"`
}

func (s *ChatCompletionRequest) IsStream() bool {
	return s.Stream != nil && *s.Stream
}

// ChatCompletionResponse 标准聊天完成响应
/*
{
  "id": "string",
  "choices": [
    {
      "finish_reason": "stop",
      "index": 0,
      "message": {
        "content": "string",
        "reasoning_content": "string",
        "tool_calls": [
          {
            "id": "string",
            "type": "function",
            "function": {
              "name": "string",
              "arguments": "string"
            }
          }
        ],
        "role": "assistant"
      },
      "logprobs": {
        "content": [
          {
            "token": "string",
            "logprob": 0,
            "bytes": [
              0
            ],
            "top_logprobs": [
              {
                "token": "string",
                "logprob": 0,
                "bytes": [
                  0
                ]
              }
            ]
          }
        ]
      }
    }
  ],
  "created": 0,
  "model": "string",
  "system_fingerprint": "string",
  "object": "chat.completion",
  "usage": {
    "completion_tokens": 0,
    "prompt_tokens": 0,
    "prompt_cache_hit_tokens": 0,
    "prompt_cache_miss_tokens": 0,
    "total_tokens": 0,
    "completion_tokens_details": {
      "reasoning_tokens": 0
    }
  }
}
*/
type ChatCompletionResponse struct {
	ID                string   `json:"id"`
	Choices           []Choice `json:"choices"`
	Created           int      `json:"created"`
	Model             string   `json:"model"`
	SystemFingerprint string   `json:"system_fingerprint"`
	Object            string   `json:"object"`
	Usage             Usage    `json:"usage"`
}

// BetaCompletionRequest 测试版完成请求
type BetaCompletionRequest struct {
	// 模型的 ID
	Model string `json:"model"`
	// 用于生成完成内容的提示
	// default "Once upon a time"
	Prompt string `json:"prompt"`
	// 在输出中，把 prompt 的内容也输出出来
	Echo *bool `json:"echo"`
	// 介于 -2.0 和 2.0 之间的数字。如果该值为正，那么新 token 会根据其在已有文本中的出现频率受到相应的惩罚，降低模型重复相同内容的可能性。
	FrequencyPenalty *float64 `json:"frequency_penalty,omitempty"` // -2<= and <=2 ,default 0
	// 是否返回所输出 token 的对数概率。如果为 true，则在 message 的 content 中返回每个输出 token 的对数概率。
	// <= 20
	Logprobs *int `json:"logprobs,omitempty"`
	// 最大生成 token 数量。
	MaxTokens int `json:"max_tokens"`
	// 介于 -2.0 和 2.0 之间的数字。如果该值为正，那么新 token 会根据其是否已在已有文本中出现受到相应的惩罚，从而增加模型谈论新主题的可能性。
	PresencePenalty *float64 `json:"presence_penalty,omitempty"` // -2<= and <=2 ,default 0
	// 一个 string 或最多包含 16 个 string 的 list，在遇到这些词时，API 将停止生成更多的 token。
	Stop *[]string `json:"stop,omitempty"`
	// 如果设置为 True，将会以 SSE（server-sent events）的形式以流式发送消息增量。消息流以 data: [DONE] 结尾。
	Stream        *bool          `json:"stream"`
	StreamOptions *StreamOptions `json:"stream_options,omitempty"`
	// 制定被补全内容的后缀。
	Suffix *string `json:"suffix,omitempty"`
	// 采样温度，介于 0 和 2 之间。更高的值，如 0.8，会使输出更随机，而更低的值，如 0.2，会使其更加集中和确定。 我们通常建议可以更改这个值或者更改 top_p，但不建议同时对两者进行修改。
	// default 1
	Temperature *float64 `json:"temperature"`
	// 作为调节采样温度的替代方案，模型会考虑前 top_p 概率的 token 的结果。所以 0.1 就意味着只有包括在最高 10% 概率中的 token 会被考虑。 我们通常建议修改这个值或者更改 temperature，但不建议同时对两者进行修改。
	// <= 1 default 1
	TopP *float64 `json:"top_p"`
}

func (s *BetaCompletionRequest) IsStream() bool {
	return s.Stream != nil && *s.Stream
}

// BetaCompletionResponse 测试版完成响应
/*
{
  "id": "string",
  "choices": [
    {
      "finish_reason": "stop",
      "index": 0,
      "logprobs": {
        "text_offset": [
          0
        ],
        "token_logprobs": [
          0
        ],
        "tokens": [
          "string"
        ],
        "top_logprobs": [
          {}
        ]
      },
      "text": "string"
    }
  ],
  "created": 0,
  "model": "string",
  "system_fingerprint": "string",
  "object": "text_completion",
  "usage": {
    "completion_tokens": 0,
    "prompt_tokens": 0,
    "prompt_cache_hit_tokens": 0,
    "prompt_cache_miss_tokens": 0,
    "total_tokens": 0,
    "completion_tokens_details": {
      "reasoning_tokens": 0
    }
  }
}
*/
type BetaCompletionResponse ChatCompletionResponse

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
