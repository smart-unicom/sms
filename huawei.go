// Package sms 华为云短信服务实现
package sms

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
)

// 华为云短信API相关常量
const (
	WSSE_HEADER_FORMAT = "UsernameToken Username=\"%s\",PasswordDigest=\"%s\",Nonce=\"%s\",Created=\"%s\"" // WSSE头格式
	AUTH_HEADER_VALUE  = "WSSE realm=\"SDP\",profile=\"UsernameToken\",type=\"Appkey\""                    // 认证头值
)

// HuaweiClient 华为云短信客户端
// 封装华为云短信API调用
type HuaweiClient struct {
	accessId   string // 访问ID
	accessKey  string // 访问密钥
	sign       string // 短信签名
	template   string // 短信模板ID
	apiAddress string // API地址
	sender     string // 发送方号码
}

// GetHuaweiClient 创建华为云短信客户端
// 参数:
//   - accessId: 华为云访问ID
//   - accessKey: 华为云访问密钥
//   - sign: 短信签名
//   - template: 短信模板ID
//   - other: 其他参数（[0]为API地址，[1]为发送方号码）
//
// 返回:
//   - *HuaweiClient: 华为云短信客户端实例
//   - error: 错误信息
func GetHuaweiClient(accessId string, accessKey string, sign string, template string, other []string) (*HuaweiClient, error) {
	if len(other) < 2 {
		return nil, fmt.Errorf("missing parameter: apiAddress or sender")
	}

	apiAddress := fmt.Sprintf("%s/sms/batchSendSms/v1", other[0])

	huaweiClient := &HuaweiClient{
		accessId:   accessId,
		accessKey:  accessKey,
		sign:       sign,
		template:   template,
		apiAddress: apiAddress,
		sender:     other[1],
	}

	return huaweiClient, nil
}

// SendMessage 发送短信
// 参考文档: https://support.huaweicloud.com/intl/zh-cn/devg-msgsms/sms_04_0012.html
// 参数:
//   - param: 短信模板参数（需要包含"code"字段）
//   - targetPhoneNumber: 目标手机号码列表
//
// 返回:
//   - error: 错误信息
func (c *HuaweiClient) SendMessage(param map[string]string, targetPhoneNumber ...string) error {
	code, ok := param["code"]
	if !ok {
		return fmt.Errorf("missing parameter: code")
	}

	if len(targetPhoneNumber) == 0 {
		return fmt.Errorf("missing parameter: targetPhoneNumber")
	}

	phoneNumbers := strings.Join(targetPhoneNumber, ",")
	templateParas := fmt.Sprintf("[\"%s\"]", code)

	body := buildRequestBody(c.sender, phoneNumbers, c.template, templateParas, "", c.sign)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	headers["Authorization"] = AUTH_HEADER_VALUE
	headers["X-WSSE"] = buildWsseHeader(c.accessId, c.accessKey)

	_, err := post(c.apiAddress, []byte(body), headers)
	return err
}

// buildRequestBody 构建请求体
// 参数:
//   - sender: 发送方号码
//   - receiver: 接收方号码
//   - templateId: 模板ID
//   - templateParas: 模板参数
//   - statusCallBack: 状态回调地址
//   - signature: 签名
//
// 返回:
//   - string: 构建的请求体
func buildRequestBody(sender, receiver, templateId, templateParas, statusCallBack, signature string) string {
	param := "from=" + url.QueryEscape(sender) + "&to=" + url.QueryEscape(receiver) + "&templateId=" + url.QueryEscape(templateId)
	if templateParas != "" {
		param += "&templateParas=" + url.QueryEscape(templateParas)
	}
	if statusCallBack != "" {
		param += "&statusCallback=" + url.QueryEscape(statusCallBack)
	}
	if signature != "" {
		param += "&signature=" + url.QueryEscape(signature)
	}
	return param
}

// buildWsseHeader 构建WSSE认证头
// 参数:
//   - appKey: 应用密钥
//   - appSecret: 应用秘密
//
// 返回:
//   - string: WSSE认证头
func buildWsseHeader(appKey, appSecret string) string {
	cTime := time.Now().Format("2006-01-02T15:04:05Z")
	nonce := uuid.NewString()
	nonce = strings.ReplaceAll(nonce, "-", "")

	h := sha256.New()
	h.Write([]byte(nonce + cTime + appSecret))
	passwordDigestBase64Str := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return fmt.Sprintf(WSSE_HEADER_FORMAT, appKey, passwordDigestBase64Str, nonce, cTime)
}

// post 发送POST请求
// 参数:
//   - url: 请求URL
//   - param: 请求参数
//   - headers: 请求头
//
// 返回:
//   - string: 响应内容
//   - error: 错误信息
func post(url string, param []byte, headers map[string]string) (string, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(param))
	if err != nil {
		return "", err
	}

	for key, header := range headers {
		req.Header.Set(key, header)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
