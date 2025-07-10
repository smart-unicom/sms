// Package sms Msg91短信服务实现
package sms

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Msg91Client Msg91短信客户端
// 封装Msg91短信API调用
type Msg91Client struct {
	authKey    string // 认证密钥
	senderId   string // 发送方ID
	templateId string // 模板ID
}

// GetMsg91Client 创建Msg91短信客户端
// 参数:
//   - senderId: 发送方ID
//   - authKey: 认证密钥
//   - templateId: 模板ID
// 返回:
//   - *Msg91Client: Msg91短信客户端实例
//   - error: 错误信息
func GetMsg91Client(senderId string, authKey string, templateId string) (*Msg91Client, error) {
	msg91Client := &Msg91Client{
		authKey:    authKey,
		senderId:   senderId,
		templateId: templateId,
	}

	return msg91Client, nil
}

// SendMessage 发送短信
// 参数:
//   - param: 短信模板参数
//   - targetPhoneNumber: 目标手机号码列表
// 返回:
//   - error: 错误信息
func (m *Msg91Client) SendMessage(param map[string]string, targetPhoneNumber ...string) error {
	if len(targetPhoneNumber) == 0 {
		return fmt.Errorf("missing parameter: targetPhoneNumber")
	}

	url := "https://control.msg91.com/api/v5/flow/"

	for _, mobile := range targetPhoneNumber {
		if strings.HasPrefix(mobile, "+") {
			mobile = mobile[1:]
		}

		payload, err := buildPayload(m.templateId, m.senderId, "0", mobile, param)
		if err != nil {
			return fmt.Errorf("SMS build payload failed: %v", err)
		}

		err = postMsg91SendRequest(url, strings.NewReader(payload), m.authKey)
		if err != nil {
			return fmt.Errorf("send message failed: %v", err)
		}
	}

	return nil
}

// buildPayload 构建请求负载
// 参数:
//   - templateId: 模板ID
//   - senderId: 发送方ID
//   - shortURL: 短链接
//   - mobiles: 手机号码
//   - variables: 变量参数
// 返回:
//   - string: JSON格式的请求负载
//   - error: 错误信息
func buildPayload(templateId, senderId, shortURL, mobiles string, variables map[string]string) (string, error) {
	payload := make(map[string]interface{})

	payload["template_id"] = templateId
	payload["sender"] = senderId
	payload["short_url"] = shortURL
	payload["mobiles"] = mobiles

	for k, v := range variables {
		payload[k] = v
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

// postMsg91SendRequest 发送Msg91请求
// 参数:
//   - url: 请求URL
//   - payload: 请求负载
//   - authKey: 认证密钥
// 返回:
//   - error: 错误信息
func postMsg91SendRequest(url string, payload io.Reader, authKey string) error {
	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authkey", authKey)

	res, _ := http.DefaultClient.Do(req)

	err := res.Body.Close()
	if err != nil {
		return err
	}

	return nil
}
