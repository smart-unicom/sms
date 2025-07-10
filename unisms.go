// Package sms UniSMS短信服务实现
package sms

import (
	"errors"
	"fmt"
	"strings"

	unisms "github.com/apistd/uni-go-sdk/sms"
)

// UnismsClient UniSMS短信客户端
// 封装UniSMS短信API调用
type UnismsClient struct {
	core     *unisms.UniSMSClient // UniSMS核心客户端
	sign     string               // 短信签名
	template string               // 短信模板
}

// GetUnismsClient 创建UniSMS短信客户端
// 参数:
//   - accessId: 访问ID
//   - accessKey: 访问密钥
//   - signature: 短信签名
//   - templateId: 短信模板ID
// 返回:
//   - *UnismsClient: UniSMS短信客户端实例
//   - error: 错误信息
func GetUnismsClient(accessId string, accessKey string, signature string, templateId string) (*UnismsClient, error) {
	client := unisms.NewClient(accessId, accessKey)

	// 检查accessId和accessKey的正确性
	msg := unisms.BuildMessage()
	msg.SetTo("test")
	msg.SetTemplateId("pub_verif_register") // 免费模板
	_, err := client.Send(msg)
	if strings.Contains(err.Error(), "[104111] InvalidAccessKeyId") {
		return nil, err
	}

	unismsClient := &UnismsClient{
		core:     client,
		sign:     signature,
		template: templateId,
	}

	return unismsClient, nil
}

// SendMessage 发送短信
// 参数:
//   - param: 短信模板参数
//   - targetPhoneNumber: 目标手机号码列表
// 返回:
//   - error: 错误信息
func (c *UnismsClient) SendMessage(param map[string]string, targetPhoneNumber ...string) error {
	if len(targetPhoneNumber) == 0 {
		return fmt.Errorf("missing parameter: targetPhoneNumber")
	}

	msg := unisms.BuildMessage()
	msg.SetTo(targetPhoneNumber...)
	msg.SetSignature(c.sign)
	msg.SetTemplateId(c.template)

	resp, err := c.core.Send(msg)
	if err != nil {
		return err
	}

	if resp.Code != "0" {
		return errors.New(resp.Message)
	}

	return nil
}
