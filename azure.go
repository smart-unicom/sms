// Package sms 微软Azure通信服务短信实现
package sms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// ACSClient Azure通信服务短信客户端
// 封装Azure通信服务短信API调用
type ACSClient struct {
	AccessToken string // 访问令牌
	Endpoint    string // 服务端点
	Message     string // 短信内容
	Sender      string // 发送方号码
}

// reqBody 短信发送请求体
type reqBody struct {
	From          string         `json:"from"`          // 发送方
	Message       string         `json:"message"`       // 短信内容
	SMSRecipients []smsRecipient `json:"smsRecipients"` // 接收方列表
}

// smsRecipient 短信接收方
type smsRecipient struct {
	To string `json:"to"` // 接收方号码
}

// GetACSClient 创建Azure通信服务短信客户端
// 参数:
//   - accessToken: Azure访问令牌
//   - message: 短信内容
//   - other: 其他参数（[0]为端点，[1]为发送方号码）
// 返回:
//   - *ACSClient: Azure通信服务短信客户端实例
//   - error: 错误信息
func GetACSClient(accessToken string, message string, other []string) (*ACSClient, error) {
	if len(other) < 2 {
		return nil, fmt.Errorf("missing parameter: endpoint or sender")
	}

	acsClient := &ACSClient{
		AccessToken: accessToken,
		Endpoint:    other[0],
		Message:     message,
		Sender:      other[1],
	}

	return acsClient, nil
}

// SendMessage 发送短信
// 参数:
//   - param: 短信模板参数（当前未使用）
//   - targetPhoneNumber: 目标手机号码列表
// 返回:
//   - error: 错误信息
func (a *ACSClient) SendMessage(param map[string]string, targetPhoneNumber ...string) error {
	if len(targetPhoneNumber) == 0 {
		return fmt.Errorf("missing parameter: targetPhoneNumber")
	}

	reqBody := &reqBody{
		From:          a.Sender,
		Message:       a.Message,
		SMSRecipients: make([]smsRecipient, 0),
	}
	for _, mobile := range targetPhoneNumber {
		reqBody.SMSRecipients = append(reqBody.SMSRecipients, smsRecipient{To: mobile})
	}

	url := fmt.Sprintf("%s/sms?api-version=2021-03-07", a.Endpoint)

	client := &http.Client{}

	requestBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("error creating request body: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+a.AccessToken)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}

	resp.Body.Close()

	return nil
}
