// Package sms Twilio短信服务实现
package sms

import (
	"fmt"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

// TwilioClient Twilio短信客户端
// 封装Twilio短信API调用
type TwilioClient struct {
	template string              // 短信模板
	core     *twilio.RestClient  // Twilio REST客户端
}

// GetTwilioClient 创建Twilio短信客户端
// 参数:
//   - accessId: Twilio账户SID
//   - accessKey: Twilio认证令牌
//   - template: 短信模板
// 返回:
//   - *TwilioClient: Twilio短信客户端实例
//   - error: 错误信息
func GetTwilioClient(accessId string, accessKey string, template string) (*TwilioClient, error) {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accessId,
		Password: accessKey,
	})

	twilioClient := &TwilioClient{
		core:     client,
		template: template,
	}

	return twilioClient, nil
}

// SendMessage 发送短信
// 注意: targetPhoneNumber[0]是发送方号码，因此targetPhoneNumber至少需要两个参数
// 参数:
//   - param: 短信模板参数（需要包含"code"字段）
//   - targetPhoneNumber: 手机号码列表（[0]为发送方，[1:]为接收方）
// 返回:
//   - error: 错误信息
func (c *TwilioClient) SendMessage(param map[string]string, targetPhoneNumber ...string) error {
	code, ok := param["code"]
	if !ok {
		return fmt.Errorf("missing parameter: code")
	}

	bodyContent := fmt.Sprintf(c.template, code)

	if len(targetPhoneNumber) < 2 {
		return fmt.Errorf("bad parameter: targetPhoneNumber")
	}

	params := &openapi.CreateMessageParams{}
	params.SetFrom(targetPhoneNumber[0])
	params.SetBody(bodyContent)

	for i := 1; i < len(targetPhoneNumber); i++ {
		params.SetTo(targetPhoneNumber[i])
		_, err := c.core.Api.CreateMessage(params)
		if err != nil {
			return err
		}
	}

	return nil
}
