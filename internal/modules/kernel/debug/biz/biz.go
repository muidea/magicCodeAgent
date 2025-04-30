package biz

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/log"
	"github.com/muidea/magicCommon/foundation/net"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicCodeAgent/internal/config"
	"github.com/muidea/magicCodeAgent/internal/modules/base/biz"
	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tracer/pkg/common"
	imkdc "github.com/muidea/magicCodeAgent/internal/modules/kernel/debug/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
)

type Debug struct {
	biz.Base
}

func New(
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine,
) *Debug {
	ptr := &Debug{
		Base: biz.New(imkdc.DebugModule, eventHub, backgroundRoutine),
	}

	return ptr
}

func (s *Debug) Run() (err *cd.Error) {
	filter := &common.MessageFilter{
		PageIndex: 1,
		PageSize:  100,
	}

	for {
		messageList := s.filterMessage(filter)
		for _, message := range messageList {
			s.traceContent(message)
		}
		if len(messageList) < filter.PageSize {
			break
		}
		filter.PageIndex++
	}

	return
}

func (s *Debug) traceContent(message *models.MessageInfo) {
	dataDir := filepath.Join(".", "data")

	// 清空或创建data目录
	if err := os.RemoveAll(dataDir); err != nil {
		log.Errorf("clear data dir failed, err:%s", err)
		return
	}
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Errorf("create data dir failed, err:%s", err)
		return
	}

	// 按sn创建子目录
	snStr := strconv.FormatInt(message.SN, 10)
	snDir := filepath.Join(dataDir, snStr)
	if err := os.MkdirAll(snDir, 0755); err != nil {
		log.Errorf("create sn dir failed, sn:%s, err:%s", snStr, err)
		return
	}

	if message.Item == nil {
		log.Warnf("message item is nil, sn:%s", snStr)
		return
	}

	// 处理Item
	item := message.Item
	switch item.Role {
	case "assistant", "recorder":
		// 直接保存MD内容
		fileName := filepath.Join(snDir, fmt.Sprintf("%s.md", item.Role))
		if err := os.WriteFile(fileName, []byte(item.Content), 0644); err != nil {
			log.Errorf("save assistant content failed, role:%s, err:%s", item.Role, err)
		}

	case "user":
		// 解析ChatCompletionRequest
		var chatReq struct {
			Messages []struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"messages"`
		}

		if err := json.Unmarshal([]byte(item.Content), &chatReq); err != nil {
			log.Errorf("unmarshal chat request failed, err:%s", err)
			return
		}

		// 保存每条message
		for i, msg := range chatReq.Messages {
			// 尝试解析为text content数组
			var textContents []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			}

			if err := json.Unmarshal([]byte(msg.Content), &textContents); err == nil {
				// 是text content数组，保存每个text
				for j, textItem := range textContents {
					if textItem.Type == "text" {
						fileName := filepath.Join(snDir, fmt.Sprintf("%s_%d_%s_%d.md", item.Role, i, msg.Role, j))
						if err := os.WriteFile(fileName, []byte(textItem.Text), 0644); err != nil {
							log.Errorf("save user text content failed, role:%s, index:%d, textIndex:%d, err:%s",
								item.Role, i, j, err)
						}
					}
				}
			} else {
				// 普通字符串内容
				fileName := filepath.Join(snDir, fmt.Sprintf("%s_%d_%s.md", item.Role, i, msg.Role))
				if err := os.WriteFile(fileName, []byte(msg.Content), 0644); err != nil {
					log.Errorf("save user message failed, role:%s, index:%d, err:%s",
						item.Role, i, err)
				}
			}
		}
	}
}

func (s *Debug) filterMessage(filter *common.MessageFilter) []*models.MessageInfo {
	result := &imbtc.MessageInfoListResult{}

	vals := url.Values{}
	vals = filter.Encode(vals)
	serverUrl := fmt.Sprintf("http://127.0.0.1:%s%s%s", config.BindPort(), common.ApiVersion, imbtc.FilterMessage)
	urlVal, urlErr := url.ParseRequestURI(serverUrl)
	if urlErr != nil {
		log.Errorf("filterMessage url parse failed, url:%s, err:%s", serverUrl, urlErr)
		return nil
	}
	urlVal.RawQuery = vals.Encode()
	_, httpErr := net.HTTPGet(http.DefaultClient, urlVal.String(), result)
	if httpErr != nil {
		return nil
	}

	return result.Values
}
