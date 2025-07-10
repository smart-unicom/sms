// Package sms 火山引擎短信服务实现
package sms

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/volcengine/volc-sdk-golang/service/sms"
)

// VolcClient 火山引擎短信客户端
// 封装火山引擎短信API调用
type VolcClient struct {
	core       *sms.SMS // 火山引擎SMS客户端
	sign       string   // 短信签名
	template   string   // 短信模板ID
	smsAccount string   // 短信账户
}

// GetVolcClient 创建火山引擎短信客户端
// 参数:
//   - accessId: 火山引擎访问ID
//   - accessKey: 火山引擎访问密钥
//   - sign: 短信签名
//   - templateId: 短信模板ID
//   - smsAccount: 短信账户列表
// 返回:
//   - *VolcClient: 火山引擎短信客户端实例
//   - error: 错误信息
func GetVolcClient(accessId, accessKey, sign, templateId string, smsAccount []string) (*VolcClient, error) {
	if len(smsAccount) == 0 {
		return nil, fmt.Errorf("missing parameter: smsAccount")
	}

	client := sms.NewInstance()
	client.Client.SetAccessKey(accessId)
	client.Client.SetSecretKey(accessKey)

	volcClient := &VolcClient{
		core:       client,
		sign:       sign,
		template:   templateId,
		smsAccount: smsAccount[0],
	}

	return volcClient, nil
}

// SendMessage 发送短信
// 参数:
//   - param: 短信模板参数
//   - targetPhoneNumber: 目标手机号码列表
// 返回:
//   - error: 错误信息
func (c *VolcClient) SendMessage(param map[string]string, targetPhoneNumber ...string) error {
	if len(targetPhoneNumber) == 0 {
		return fmt.Errorf("missing parameter: targetPhoneNumber")
	}

	requestParam, err := json.Marshal(param)
	if err != nil {
		return err
	}

	req := &sms.SmsRequest{
		SmsAccount:    c.smsAccount,
		Sign:          c.sign,
		TemplateID:    c.template,
		TemplateParam: string(requestParam),
		PhoneNumbers:  strings.Join(targetPhoneNumber, ","),
	}

	resp, statusCode, err := c.core.Send(req)
	if err != nil {
		return fmt.Errorf("send message failed, error: %q", err.Error())
	}
	if statusCode < 200 || statusCode > 299 {
		return fmt.Errorf("send message failed, statusCode: %d", statusCode)
	}
	if resp.ResponseMetadata.Error != nil {
		return fmt.Errorf("send message failed, code: %q, message: %q", resp.ResponseMetadata.Error.Code, resp.ResponseMetadata.Error.Message)
	}

	return nil
}
