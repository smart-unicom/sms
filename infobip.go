// Package sms Infobip短信服务实现
package sms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// InfobipClient Infobip短信客户端
// 封装Infobip短信API调用
type InfobipClient struct {
	baseUrl  string // API基础URL
	sender   string // 发送方标识
	apiKey   string // API密钥
	template string // 短信模板
}

// InfobipConfigService Infobip配置服务
type InfobipConfigService struct {
	baseUrl string // API基础URL
	sender  string // 发送方标识
	apiKey  string // API密钥
}

// SmsService 短信服务
type SmsService struct {
	configService InfobipConfigService // 配置服务
}

// MessageData 消息数据结构体
type MessageData struct {
	Messages []Message `json:"messages"` // 消息列表
}

// Message 消息结构体
type Message struct {
	From         string        `json:"from"`         // 发送方
	Destinations []Destination `json:"destinations"` // 目标列表
	Text         string        `json:"text"`         // 消息内容
}

// Destination 目标结构体
type Destination struct {
	To string `json:"to"` // 目标号码
}

// GetInfobipClient 创建Infobip短信客户端
// 参数:
//   - sender: 发送方标识
//   - apiKey: API密钥
//   - template: 短信模板
//   - baseUrl: API基础URL列表
// 返回:
//   - *InfobipClient: Infobip短信客户端实例
//   - error: 错误信息
func GetInfobipClient(sender string, apiKey string, template string, baseUrl []string) (*InfobipClient, error) {
	if len(baseUrl) == 0 {
		return nil, fmt.Errorf("missing parameter: baseUrl")
	}

	infobipClient := &InfobipClient{
		baseUrl:  baseUrl[0],
		sender:   sender,
		apiKey:   apiKey,
		template: template,
	}

	return infobipClient, nil
}

// SendMessage 发送短信
// 参数:
//   - param: 短信模板参数
//   - targetPhoneNumber: 目标手机号码列表
// 返回:
//   - error: 错误信息
func (c *InfobipClient) SendMessage(param map[string]string, targetPhoneNumber ...string) error {
	code, ok := param["code"]
	if !ok {
		return fmt.Errorf("missing parameter: code")
	}

	if len(targetPhoneNumber) == 0 {
		return fmt.Errorf("missin parer: trgetPhoneNumber")
	}

	mobile := targetPhoneNumber[0]

	if strings.HasPrefix(mobile, "0") {
		mobile = "886" + mobile[1:]
	}
	if strings.HasPrefix(mobile, "+") {
		mobile = mobile[1:]
	}

	endpoint := fmt.Sprintf("%s/sms/2/text/advanced", c.baseUrl)
	text := fmt.Sprintf(c.template, code)

	messageData := MessageData{
		Messages: []Message{
			{
				From: c.sender,
				Destinations: []Destination{
					{
						To: mobile,
					},
				},
				Text: text,
			},
		},
	}
	headers := map[string]string{
		"Authorization": fmt.Sprintf("App %s", c.apiKey),
		"Content-Type":  "application/json",
	}

	messageDataBytes, _ := json.Marshal(messageData)
	req, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(messageDataBytes))
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
