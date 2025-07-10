// Package sms 亚马逊SNS短信服务实现
package sms

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
)

// AmazonSNSClient 亚马逊SNS短信客户端
// 封装亚马逊SNS短信API调用
type AmazonSNSClient struct {
	svc      snsiface.SNSAPI // SNS服务接口
	template string          // 短信模板
}

// GetAmazonSNSClient 创建亚马逊SNS短信客户端
// 参数:
//   - accessKeyId: AWS访问密钥ID
//   - secretAccessKey: AWS访问密钥
//   - template: 短信模板
//   - region: AWS区域列表
// 返回:
//   - *AmazonSNSClient: 亚马逊SNS短信客户端实例
//   - error: 错误信息
func GetAmazonSNSClient(accessKeyId string, secretAccessKey string, template string, region []string) (*AmazonSNSClient, error) {
	if len(region) == 0 {
		return nil, fmt.Errorf("missing parameter: region")
	}

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region[0]),
		Credentials: credentials.NewStaticCredentials(accessKeyId, secretAccessKey, ""),
	})
	if err != nil {
		return nil, err
	}

	svc := sns.New(sess)

	snsClient := &AmazonSNSClient{
		svc:      svc,
		template: template,
	}

	return snsClient, nil
}

// SendMessage 发送短信
// 参数:
//   - param: 短信模板参数（需要包含"code"字段）
//   - targetPhoneNumber: 目标手机号码列表
// 返回:
//   - error: 错误信息
func (a *AmazonSNSClient) SendMessage(param map[string]string, targetPhoneNumber ...string) error {
	code, ok := param["code"]
	if !ok {
		return fmt.Errorf("missing parameter: code")
	}

	bodyContent := fmt.Sprintf(a.template, code)

	if len(targetPhoneNumber) == 0 {
		return fmt.Errorf("missing parameter: targetPhoneNumber")
	}

	messageAttributes := make(map[string]*sns.MessageAttributeValue)
	for k, v := range param {
		messageAttributes[k] = &sns.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(v),
		}
	}

	for i := 0; i < len(targetPhoneNumber); i++ {
		_, err := a.svc.Publish(&sns.PublishInput{
			Message:           &bodyContent,
			PhoneNumber:       &targetPhoneNumber[i],
			MessageAttributes: messageAttributes,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
