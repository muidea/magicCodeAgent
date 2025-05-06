package util

import (
	"regexp"
	"strings"
)

// SpliteBlocks 从content中分解出blocks
// 单个block的格式如下
// <blockTag>blockContent</blockTag>
// blockTag:是一个完整的单词，不包含空格换行符等特殊符号
// blockContent可以是任意内容, 可以包含空格和换行符
// 多个blocks用空行分隔开
func SpliteBlocks(content string) (ret []string, err error) {
	// 初始化返回值
	ret = []string{}

	// 按空行分割内容
	lines := strings.Split(content, "\n")

	// 当前正在处理的块
	var currentBlock strings.Builder
	// 当前块的标签名
	var currentTag string
	// 是否在块内
	inBlock := false
	// 是否已经找到开始标签
	foundStartTag := false

	// 用于匹配开始标签的正则表达式
	startTagRegex := regexp.MustCompile(`^\s*<([a-zA-Z0-9_]+)>\s*$`)
	// 用于匹配结束标签的正则表达式
	endTagRegex := regexp.MustCompile(`^\s*</([a-zA-Z0-9_]+)>\s*$`)

	for _, line := range lines {
		// 检查是否是开始标签
		if !inBlock {
			startMatches := startTagRegex.FindStringSubmatch(line)
			if len(startMatches) == 2 {
				// 找到开始标签
				currentTag = startMatches[1]
				inBlock = true
				foundStartTag = true
				currentBlock.WriteString(line)
				continue
			}
		}

		// 检查是否是结束标签
		if inBlock {
			endMatches := endTagRegex.FindStringSubmatch(line)
			if len(endMatches) == 2 && endMatches[1] == currentTag {
				// 找到匹配的结束标签
				currentBlock.WriteString("\n" + line)
				// 添加完整的块到结果中
				ret = append(ret, currentBlock.String())
				// 重置状态
				currentBlock.Reset()
				currentTag = ""
				inBlock = false
				foundStartTag = false
				continue
			}
		}

		// 如果在块内，添加当前行到块中
		if inBlock {
			if foundStartTag {
				// 如果已经找到开始标签，添加换行符
				currentBlock.WriteString("\n")
				foundStartTag = false
			}
			currentBlock.WriteString(line + "\n")
		}
	}

	return ret, nil
}
