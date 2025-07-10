// Package sms 互亿无线短信服务实现
package sms

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// HuyiClient 互亿无线短信客户端
// 封装互亿无线短信API调用
type HuyiClient struct {
	appId    string // 应用ID
	appKey   string // 应用密钥
	template string // 短信模板
}

// GetHuyiClient 创建互亿无线短信客户端
// 参数:
//   - appId: 应用ID
//   - appKey: 应用密钥
//   - template: 短信模板
// 返回:
//   - *HuyiClient: 互亿无线短信客户端实例
//   - error: 错误信息
func GetHuyiClient(appId string, appKey string, template string) (*HuyiClient, error) {
	return &HuyiClient{
		appId:    appId,
		appKey:   appKey,
		template: template,
	}, nil
}

// GetMd5String 计算字符串的MD5哈希值
// 参数:
//   - s: 待计算的字符串
// 返回:
//   - string: MD5哈希值的十六进制字符串
func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// SendMessage 发送短信
// 参数:
//   - param: 短信模板参数
//   - targetPhoneNumber: 目标手机号码列表
// 返回:
//   - error: 错误信息
func (hc *HuyiClient) SendMessage(param map[string]string, targetPhoneNumber ...string) error {
	code, ok := param["code"]
	if !ok {
		return fmt.Errorf("missing parameter: code")
	}

	if len(targetPhoneNumber) == 0 {
		return fmt.Errorf("missin parer: trgetPhoneNumber")
	}

	_now := strconv.FormatInt(time.Now().Unix(), 10)
	smsContent := fmt.Sprintf(hc.template, code)
	v := url.Values{}
	v.Set("account", hc.appId)
	v.Set("content", smsContent)
	v.Set("time", _now)
	passwordStr := hc.appId + hc.appKey + "%s" + smsContent + _now
	for _, mobile := range targetPhoneNumber {
		password := fmt.Sprintf(passwordStr, mobile)
		v.Set("password", GetMd5String(password))
		v.Set("mobile", mobile)

		body := strings.NewReader(v.Encode()) // 编码表单数据
		client := &http.Client{}
		req, _ := http.NewRequest("POST", "http://106.ihuyi.com/webservice/sms.php?method=Submit&format=json", body)

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

		resp, err := client.Do(req) // 发送远程请求
		if err != nil {
			return err
		}
		defer resp.Body.Close() // 关闭ReadCloser
		_, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
	}

	return nil
}
