// Package sms GCCPAY短信服务实现
package sms

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// GCCPAYClient GCCPAY短信客户端
// 封装GCCPAY短信API调用
type GCCPAYClient struct {
	clientname string // 客户端名称
	secret     string // 客户端密钥
	template   string // 短信模板
}

// params 短信发送参数结构体
type params struct {
	Mobile         string            `json:"mobile"`          // 手机号码
	TemplateCode   string            `json:"template_code"`   // 模板代码
	TemplateParams map[string]string `json:"template_params"` // 模板参数
}

// GetGCCPAYClient 创建GCCPAY短信客户端
// 参数:
//   - clientname: 客户端名称
//   - secret: 客户端密钥
//   - template: 短信模板
// 返回:
//   - *GCCPAYClient: GCCPAY短信客户端实例
//   - error: 错误信息
func GetGCCPAYClient(clientname string, secret string, template string) (*GCCPAYClient, error) {
	gccPayClient := &GCCPAYClient{
		clientname: clientname,
		secret:     secret,
		template:   template,
	}

	return gccPayClient, nil
}

// RandStringBytesCrypto 生成指定长度的随机字符串
// 参数:
//   - n: 字节长度
// 返回:
//   - string: Base64编码的随机字符串
//   - error: 错误信息
func RandStringBytesCrypto(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

// Md5 计算字符串的MD5哈希值
// 参数:
//   - str: 待计算的字符串
// 返回:
//   - string: MD5哈希值的十六进制字符串
func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// SendMessage 发送短信
// 参数:
//   - param: 短信模板参数
//   - targetPhoneNumber: 目标手机号码列表
// 返回:
//   - error: 错误信息
func (c *GCCPAYClient) SendMessage(param map[string]string, targetPhoneNumber ...string) error {
	_, ok := param["code"]
	if !ok {
		return fmt.Errorf("missing parameter: code")
	}

	if len(targetPhoneNumber) == 0 {
		return fmt.Errorf("missing parameter: targetPhoneNumber")
	}

	reqParams := make(map[string]params)

	for _, mobile := range targetPhoneNumber {
		if strings.HasPrefix(mobile, "+") {
			mobile = mobile[1:]
		}
		randomString, err := RandStringBytesCrypto(16)
		if err != nil {
			return fmt.Errorf("SMS key generation failed")
		}

		reqParams[randomString] = params{
			Mobile:         mobile,
			TemplateCode:   c.template,
			TemplateParams: param,
		}
	}

	requestBody := new(bytes.Buffer)
	err := json.NewEncoder(requestBody).Encode(reqParams)
	if err != nil {
		return fmt.Errorf("SMS sending failed")
	}

	// 生成签名
	timestamp := time.Now().Unix()

	sign := Md5(fmt.Sprintf("%s%d%s", c.clientname, timestamp, c.secret))

	reqUrl := "https://smscenter.sgate.sa/api/v1/client/sendSms"

	// 发送请求
	req, _ := http.NewRequest("POST", reqUrl, requestBody)
	req.Header.Set("clientname", c.clientname)
	req.Header.Set("timestamp", fmt.Sprintf("%d", timestamp))
	req.Header.Set("sign", sign)
	req.Header.Set("content-type", "application/json;")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	err = resp.Body.Close()
	if err != nil {
		return err
	}

	return nil
}
