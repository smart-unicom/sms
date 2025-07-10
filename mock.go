// Package sms 模拟短信服务实现
package sms

// Mocker 模拟短信客户端
// 用于测试环境，不实际发送短信
type Mocker struct{}

// 确保Mocker实现了SmsProvider接口
var _ SmsProvider = &Mocker{}

// NewMocker 创建模拟短信客户端
// 参数:
//   - accessId: 访问ID（模拟用，不实际使用）
//   - accessKey: 访问密钥（模拟用，不实际使用）
//   - sign: 短信签名（模拟用，不实际使用）
//   - templateId: 短信模板ID（模拟用，不实际使用）
//   - smsAccount: 短信账户（模拟用，不实际使用）
// 返回:
//   - *Mocker: 模拟短信客户端实例
//   - error: 错误信息
func NewMocker(accessId, accessKey, sign, templateId string, smsAccount []string) (*Mocker, error) {
	return &Mocker{}, nil
}

// SendMessage 模拟发送短信
// 参数:
//   - param: 短信模板参数（不实际使用）
//   - targetPhoneNumber: 目标手机号码列表（不实际使用）
// 返回:
//   - error: 始终返回nil（模拟成功）
func (m *Mocker) SendMessage(param map[string]string, targetPhoneNumber ...string) error {
	return nil
}
