// Package sms 短信服务提供商统一接口
package sms

import "fmt"

// 短信服务提供商常量定义
const (
	SMS_TWILIO  string = "Twilio_SMS"        // Twilio短信服务
	SMS_AMAZON  string = "Amazon_SNS"        // 亚马逊SNS短信服务
	SMS_AZURE   string = "Azure_ACS"         // 微软Azure通信服务
	SMS_MSG91   string = "Msg91_SMS"         // Msg91短信服务
	SMS_GCCPAY  string = "GCCPAY_SMS"        // GCCPAY短信服务
	SMS_INFOBIP string = "Infobip_SMS"       // Infobip短信服务
	SMS_SUBMAIL string = "SUBMAIL_SMS"       // SUBMAIL短信服务
	SMS_SMSBAO  string = "SmsBao_SMS"        // 短信宝服务
	SMS_ALIYUN  string = "Aliyun_SMS"        // 阿里云短信服务
	SMS_TENCENT string = "Tencent_Cloud_SMS" // 腾讯云短信服务
	SMS_BAIdU   string = "Baidu_Cloud_SMS"   // 百度云短信服务
	SMS_VOCL    string = "Volc_Engine_SMS"   // 火山引擎短信服务
	SMS_HUAWEI  string = "Huawei_Cloud_SMS"  // 华为云短信服务
	SMS_UCloud  string = "UCloud_SMS"        // UCloud短信服务
	SMS_HUYI    string = "Huyi_SMS"          // 互亿无线短信服务
	SMS_MOCK    string = "Mock SMS"          // 模拟短信服务
	SMS_NETGSM  string = "Netgsm_SMS"        // Netgsm短信服务
	SMS_OSONI   string = "OSON_SMS"          // OSON短信服务
	SMS_UNI     string = "Uni_SMS"           // Uni短信服务
)

// SmsProvider 短信服务提供商接口
// 定义了发送短信的统一接口
type SmsProvider interface {
	// SendMessage 发送短信
	// 参数:
	//   - param: 短信模板参数
	//   - targetPhoneNumber: 目标手机号码列表
	// 返回:
	//   - error: 错误信息
	SendMessage(param map[string]string, targetPhoneNumber ...string) error
}

// NewSmsProvider 创建短信服务提供商实例
// 参数:
//   - provider: 服务提供商类型
//   - accessId: 访问ID
//   - accessKey: 访问密钥
//   - sign: 短信签名
//   - template: 短信模板
//   - other: 其他参数
// 返回:
//   - SmsProvider: 短信服务提供商实例
//   - error: 错误信息
func NewSmsProvider(provider string, accessId string, accessKey string, sign string, template string, other ...string) (SmsProvider, error) {
	switch provider {
	case SMS_TWILIO:
		return GetTwilioClient(accessId, accessKey, template)
	case SMS_AMAZON:
		return GetAmazonSNSClient(accessId, accessKey, template, other)
	case SMS_AZURE:
		return GetACSClient(accessKey, template, other)
	case SMS_MSG91:
		return GetMsg91Client(accessId, accessKey, template)
	case SMS_GCCPAY:
		return GetGCCPAYClient(accessId, accessKey, template)
	case SMS_INFOBIP:
		return GetInfobipClient(accessId, accessKey, template, other)
	case SMS_SUBMAIL:
		return GetSubmailClient(accessId, accessKey, template)
	case SMS_SMSBAO:
		return GetSmsbaoClient(accessId, accessKey, sign, template, other)
	case SMS_ALIYUN:
		return GetAliyunClient(accessId, accessKey, sign, template)
	case SMS_TENCENT:
		return GetTencentClient(accessId, accessKey, sign, template, other)
	case SMS_BAIdU:
		return GetBceClient(accessId, accessKey, sign, template, other)
	case SMS_VOCL:
		return GetVolcClient(accessId, accessKey, sign, template, other)
	case SMS_HUAWEI:
		return GetHuaweiClient(accessId, accessKey, sign, template, other)
	case SMS_UCloud:
		return GetUcloudClient(accessId, accessKey, sign, template, other)
	case SMS_HUYI:
		return GetHuyiClient(accessId, accessKey, template)
	case SMS_NETGSM:
		return GetNetgsmClient(accessId, accessKey, sign, template)
	case SMS_MOCK:
		return NewMocker(accessId, accessKey, sign, template, other)
	case SMS_OSONI:
		return GetOsonClient(accessId, accessKey, sign, template)
	case SMS_UNI:
		return GetUnismsClient(accessId, accessKey, sign, template)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}
