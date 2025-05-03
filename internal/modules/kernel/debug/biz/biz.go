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
	"github.com/sashabaranov/go-openai"

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

	if config.DebugFlag() {
		ptr.SubscribeFunc(imkdc.DebugRealtimeMessage, ptr.recordRealtimeMessage)

		ptr.SubscribeFunc(imkdc.DebugHistoryMessage, ptr.recordHistoryMessage)
	}

	return ptr
}

func (s *Debug) Run() (err *cd.Error) {
	if !config.DebugFlag() {
		return
	}

	s.traceAllHistoryContent()
	return
}

func (s *Debug) checkDataDir(dataDir string) (err *cd.Error) {
	// 清空或创建data目录
	if errVal := os.RemoveAll(dataDir); errVal != nil {
		err = cd.NewError(cd.IllegalParam, fmt.Sprintf("clear data dir failed, error:%s", errVal.Error()))
		log.Errorf("clear data dir failed, err:%s", errVal)
		return
	}
	if errVal := os.MkdirAll(dataDir, 0755); errVal != nil {
		err = cd.NewError(cd.IllegalParam, fmt.Sprintf("create data dir failed, error:%s", errVal.Error()))
		log.Errorf("create data dir failed, err:%s", errVal)
		return
	}
	return
}

func (s *Debug) recordRealtimeMessage(ev event.Event, re event.Result) {
	var err *cd.Error
	defer func() {
		if re != nil {
			re.Set(nil, err)
		}
	}()

	msgPtr, msgOK := ev.Data().(*models.MessageInfo)
	if !msgOK {
		err = cd.NewError(cd.IllegalParam, "recordMessage message is nil")
		log.Warnf("recordMessage message is nil, ev:%v", ev)
		return
	}

	dataDir := filepath.Join(config.CurrentWorkPath(), "data")
	if err := s.checkDataDir(dataDir); err != nil {
		return
	}

	s.traceContent(dataDir, msgPtr)
}

func (s *Debug) recordHistoryMessage(ev event.Event, re event.Result) {
	var err *cd.Error
	defer func() {
		if re != nil {
			re.Set(nil, err)
		}
	}()

	err = s.traceAllHistoryContent()
}

func (s *Debug) traceAllHistoryContent() (err *cd.Error) {
	filter := &common.MessageFilter{
		PageIndex: 1,
		PageSize:  100,
	}

	dataDir := filepath.Join(config.CurrentWorkPath(), "data")
	if err = s.checkDataDir(dataDir); err != nil {
		return
	}

	for {
		messageList := s.filterMessage(filter)
		for _, message := range messageList {
			s.traceContent(dataDir, message)
		}
		if len(messageList) < filter.PageSize {
			break
		}
		filter.PageIndex++
	}

	return
}

func (s *Debug) traceContent(dataDir string, message *models.MessageInfo) {
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
	case common.Assistant, common.Recorder:
		// 直接保存MD内容
		fileName := filepath.Join(snDir, fmt.Sprintf("%s.md", item.Role))
		if err := os.WriteFile(fileName, []byte(item.Content), 0644); err != nil {
			log.Errorf("save assistant content failed, role:%s, err:%s", item.Role, err)
		}
	case common.User:
		// 解析ChatCompletionRequest
		chatReq := &openai.ChatCompletionRequest{}
		if err := json.Unmarshal([]byte(item.Content), chatReq); err != nil {
			log.Errorf("unmarshal chat request failed, err:%s", err)
			return
		}
		// 保存每条message
		for i, msg := range chatReq.Messages {
			content := []byte(msg.Content)
			if msg.Role == openai.ChatMessageRoleSystem {
				// 普通字符串内容
				fileName := filepath.Join(snDir, fmt.Sprintf("%s_%d_%s.md", item.Role, i, msg.Role))
				if err := os.WriteFile(fileName, content, 0644); err != nil {
					log.Errorf("save user message failed, role:%s, index:%d, err:%s",
						item.Role, i, err)
				}
				continue
			}

			if msg.Role == openai.ChatMessageRoleUser {
				for j, textItem := range msg.MultiContent {
					if textItem.Type == "text" && textItem.Text != "" {
						fileName := filepath.Join(snDir, fmt.Sprintf("%s_%d_%s_%d.md", item.Role, i, msg.Role, j))
						if err := os.WriteFile(fileName, []byte(textItem.Text), 0644); err != nil {
							log.Errorf("save user text content failed, role:%s, index:%d, textIndex:%d, err:%s",
								item.Role, i, j, err)
						}
					}
				}

				continue
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
