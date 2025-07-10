// Package sms 阿里云短信服务实现
package sms

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

// AliyunClient 阿里云短信客户端
// 封装阿里云短信API调用
type AliyunClient struct {
	template string             // 短信模板ID
	sign     string             // 短信签名
	core     *dysmsapi.Client   // 阿里云SDK客户端
}

// AliyunResult 阿里云短信发送结果
type AliyunResult struct {
	RequestId string // 请求ID
	Message   string // 响应消息
}

// GetAliyunClient 创建阿里云短信客户端
// 参数:
//   - accessId: 阿里云访问ID
//   - accessKey: 阿里云访问密钥
//   - sign: 短信签名
//   - template: 短信模板ID
// 返回:
//   - *AliyunClient: 阿里云短信客户端实例
//   - error: 错误信息
func GetAliyunClient(accessId string, accessKey string, sign string, template string) (*AliyunClient, error) {
	region := "cn-hangzhou"
	client, err := dysmsapi.NewClientWithAccessKey(region, accessId, accessKey)
	if err != nil {
		return nil, err
	}

	aliyunClient := &AliyunClient{
		template: template,
		core:     client,
		sign:     sign,
	}

	return aliyunClient, nil
}

// SendMessage 发送短信
// 参数:
//   - param: 短信模板参数
//   - targetPhoneNumber: 目标手机号码列表
// 返回:
//   - error: 错误信息
func (c *AliyunClient) SendMessage(param map[string]string, targetPhoneNumber ...string) error {
	requestParam, err := json.Marshal(param)
	if err != nil {
		return err
	}

	if len(targetPhoneNumber) == 0 {
		return fmt.Errorf("missing parameter: targetPhoneNumber")
	}

	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = strings.Join(targetPhoneNumber, ",")
	request.TemplateCode = c.template
	request.TemplateParam = string(requestParam)
	request.SignName = c.sign

	response, err := c.core.SendSms(request)
	if err != nil {
		return err
	}

	if response.Code != "OK" {
		aliyunResult := AliyunResult{}
		err = json.Unmarshal(response.GetHttpContentBytes(), &aliyunResult)
		if err != nil {
			return err
		}

		if aliyunResult.Message != "" {
			return fmt.Errorf(aliyunResult.Message)
		}
	}

	return nil
}
