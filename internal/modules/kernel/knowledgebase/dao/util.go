package dao

import "strings"

// ContextType 表示当前查询的上下文类型
type ContextType int

const (
	ContextFAQ          ContextType = iota // 精确问答
	ContextOpenDialogue                    // 开放式对话
	ContextMixed                           // 混合场景
)

// AutoAdjustWeights 根据文本和上下文类型自动推荐 vectorWeight 和 ftsWeight
func AutoAdjustWeights(text string, context ContextType, keywords []string) (vectorWeight, ftsWeight float32) {
	textLen := len(strings.TrimSpace(text))
	numKeywords := len(keywords)
	keywordDensity := float32(0)
	if textLen > 0 {
		keywordDensity = float32(numKeywords) / float32(textLen)
	}

	// 默认权重
	vectorWeight = 0.5
	ftsWeight = 0.5

	switch context {
	case ContextFAQ:
		// 精确问答，关键词匹配重要
		vectorWeight = 0.3
		ftsWeight = 0.7

		// 如果文本非常短，关键词权重进一步增加
		if textLen < 20 {
			vectorWeight = 0.2
			ftsWeight = 0.8
		}
	case ContextOpenDialogue:
		// 开放式对话，语义匹配重要
		vectorWeight = 0.8
		ftsWeight = 0.2

		// 如果关键词密度高，适当增加 FTS 权重
		if keywordDensity > 0.05 {
			vectorWeight = 0.7
			ftsWeight = 0.3
		}
	case ContextMixed:
		// 混合场景，语义和关键词权重平衡
		vectorWeight = 0.6
		ftsWeight = 0.4

		// 如果文本非常短且关键词密度高，增加 FTS 权重
		if textLen < 50 && keywordDensity > 0.05 {
			vectorWeight = 0.5
			ftsWeight = 0.5
		}
	}

	// 确保总和为 1
	total := vectorWeight + ftsWeight
	vectorWeight /= total
	ftsWeight /= total

	return vectorWeight, ftsWeight
}
