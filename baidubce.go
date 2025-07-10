// Package sms 百度云短信服务实现
package sms

import (
	"fmt"
	"strings"

	"github.com/baidubce/bce-sdk-go/services/sms"
	"github.com/baidubce/bce-sdk-go/services/sms/api"
)

// BaiduClient 百度云短信客户端
// 封装百度云短信API调用
type BaiduClient struct {
	sign     string      // 短信签名ID
	template string      // 短信模板ID
	core     *sms.Client // 百度云SMS客户端
}

// GetBceClient 创建百度云短信客户端
// 参数:
//   - accessId: 百度云访问ID
//   - accessKey: 百度云访问密钥
//   - sign: 短信签名ID
//   - template: 短信模板ID
//   - endpoint: 服务端点列表
// 返回:
//   - *BaiduClient: 百度云短信客户端实例
//   - error: 错误信息
func GetBceClient(accessId, accessKey, sign, template string, endpoint []string) (*BaiduClient, error) {
	if len(endpoint) == 0 {
		return nil, fmt.Errorf("missing parameter: endpoint")
	}

	client, err := sms.NewClient(accessId, accessKey, endpoint[0])
	if err != nil {
		return nil, err
	}

	bceClient := &BaiduClient{
		sign:     sign,
		template: template,
		core:     client,
	}

	return bceClient, nil
}

// SendMessage 发送短信
// 参数:
//   - param: 短信模板参数（需要包含"code"字段）
//   - targetPhoneNumber: 目标手机号码列表
// 返回:
//   - error: 错误信息
func (c *BaiduClient) SendMessage(param map[string]string, targetPhoneNumber ...string) error {
	code, ok := param["code"]
	if !ok {
		return fmt.Errorf("missing parameter: code")
	}

	if len(targetPhoneNumber) == 0 {
		return fmt.Errorf("missing parameter: targetPhoneNumber")
	}

	contentMap := make(map[string]interface{})
	contentMap["code"] = code

	sendSmsArgs := &api.SendSmsArgs{
		Mobile:      strings.Join(targetPhoneNumber, ","),
		SignatureId: c.sign,
		Template:    c.template,
		ContentVar:  contentMap,
	}

	_, err := c.core.SendSms(sendSmsArgs)
	if err != nil {
		return err
	}

	return nil
}
