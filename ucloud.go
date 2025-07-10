// Package sms UCloud短信服务实现
package sms

import (
	"fmt"

	"github.com/ucloud/ucloud-sdk-go/services/usms"
	"github.com/ucloud/ucloud-sdk-go/ucloud"
	"github.com/ucloud/ucloud-sdk-go/ucloud/auth"
	"github.com/ucloud/ucloud-sdk-go/ucloud/config"
)

// UcloudClient UCloud短信客户端
// 封装UCloud短信API调用
type UcloudClient struct {
	core       *usms.USMSClient // UCloud USMS客户端
	ProjectId  string           // 项目ID
	PrivateKey string           // 私钥
	PublicKey  string           // 公钥
	Sign       string           // 短信签名
	Template   string           // 短信模板ID
}

// GetUcloudClient 创建UCloud短信客户端
// 参数:
//   - publicKey: UCloud公钥
//   - privateKey: UCloud私钥
//   - sign: 短信签名
//   - template: 短信模板ID
//   - projectId: 项目ID列表
// 返回:
//   - *UcloudClient: UCloud短信客户端实例
//   - error: 错误信息
func GetUcloudClient(publicKey string, privateKey string, sign string, template string, projectId []string) (*UcloudClient, error) {
	if len(projectId) == 0 {
		return nil, fmt.Errorf("missing parameter: projectId")
	}

	cfg := config.NewConfig()
	cfg.ProjectId = projectId[0]
	credential := auth.NewCredential()
	credential.PublicKey = publicKey
	credential.PrivateKey = privateKey

	client := usms.NewClient(&cfg, &credential)

	ucloudClient := &UcloudClient{
		core:       client,
		ProjectId:  projectId[0],
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Sign:       sign,
		Template:   template,
	}

	return ucloudClient, nil
}

// SendMessage 发送短信
// 参数:
//   - param: 短信模板参数（需要包含"code"字段）
//   - targetPhoneNumber: 目标手机号码列表
// 返回:
//   - error: 错误信息
func (c *UcloudClient) SendMessage(param map[string]string, targetPhoneNumber ...string) error {
	code, ok := param["code"]
	if !ok {
		return fmt.Errorf("missing parameter: code")
	}

	if len(targetPhoneNumber) == 0 {
		return fmt.Errorf("missing parameter: targetPhoneNumber")
	}

	req := c.core.NewSendUSMSMessageRequest()
	req.SigContent = ucloud.String(c.Sign)
	req.TemplateId = ucloud.String(c.Template)
	req.PhoneNumbers = targetPhoneNumber
	req.TemplateParams = []string{code}
	response, err := c.core.SendUSMSMessage(req)
	if err != nil {
		return err
	}
	if response.RetCode != 0 {
		return fmt.Errorf(response.Message)
	}
	return nil
}
