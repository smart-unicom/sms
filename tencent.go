// Package sms 腾讯云短信服务实现
package sms

import (
	"fmt"
	"strconv"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

// TencentClient 腾讯云短信客户端
// 封装腾讯云短信API调用
type TencentClient struct {
	core     *sms.Client // 腾讯云SDK客户端
	appId    string      // 应用ID
	sign     string      // 短信签名
	template string      // 短信模板ID
}

// GetTencentClient 创建腾讯云短信客户端
// 参数:
//   - accessId: 腾讯云访问ID
//   - accessKey: 腾讯云访问密钥
//   - sign: 短信签名
//   - templateId: 短信模板ID
//   - appId: 应用ID列表
// 返回:
//   - *TencentClient: 腾讯云短信客户端实例
//   - error: 错误信息
func GetTencentClient(accessId string, accessKey string, sign string, templateId string, appId []string) (*TencentClient, error) {
	if len(appId) == 0 {
		return nil, fmt.Errorf("missing parameter: appId")
	}

	credential := common.NewCredential(accessId, accessKey)
	config := profile.NewClientProfile()
	config.HttpProfile.ReqMethod = "POST"

	region := "ap-guangzhou"
	client, err := sms.NewClient(credential, region, config)
	if err != nil {
		return nil, err
	}

	tencentClient := &TencentClient{
		core:     client,
		appId:    appId[0],
		sign:     sign,
		template: templateId,
	}

	return tencentClient, nil
}

// SendMessage 发送短信
// 参数:
//   - param: 短信模板参数（按索引顺序："0", "1", "2"...）
//   - targetPhoneNumber: 目标手机号码列表
// 返回:
//   - error: 错误信息
func (c *TencentClient) SendMessage(param map[string]string, targetPhoneNumber ...string) error {
	if len(targetPhoneNumber) == 0 {
		return fmt.Errorf("missing parameter: targetPhoneNumber")
	}

	var paramArray []string
	index := 0
	for {
		value := param[strconv.Itoa(index)]
		if len(value) == 0 {
			break
		}
		paramArray = append(paramArray, value)
		index++
	}

	request := sms.NewSendSmsRequest()
	request.SmsSdkAppId = common.StringPtr(c.appId)
	request.SignName = common.StringPtr(c.sign)
	request.TemplateParamSet = common.StringPtrs(paramArray)
	request.TemplateId = common.StringPtr(c.template)
	request.PhoneNumberSet = common.StringPtrs(targetPhoneNumber)

	_, err := c.core.SendSms(request)
	return err
}
