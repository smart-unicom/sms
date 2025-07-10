// Package sms 短信宝短信服务实现
package sms

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// SmsBaoClient 短信宝客户端
// 封装短信宝API调用
type SmsBaoClient struct {
	username string // 用户名
	apikey   string // API密钥
	sign     string // 短信签名
	template string // 短信模板
	goodsid  string // 商品ID
}

// GetSmsbaoClient 创建短信宝客户端
// 参数:
//   - username: 短信宝用户名
//   - apikey: 短信宝API密钥
//   - sign: 短信签名
//   - template: 短信模板
//   - other: 其他参数（[0]为商品ID，可选）
// 返回:
//   - *SmsBaoClient: 短信宝客户端实例
//   - error: 错误信息
func GetSmsbaoClient(username string, apikey string, sign string, template string, other []string) (*SmsBaoClient, error) {
	var goodsid string
	if len(other) == 0 {
		goodsid = ""
	} else {
		goodsid = other[0]
	}
	return &SmsBaoClient{
		username: username,
		apikey:   apikey,
		sign:     sign,
		template: template,
		goodsid:  goodsid,
	}, nil
}

// SendMessage 发送短信
// 参数:
//   - param: 短信模板参数（需要包含"code"字段）
//   - targetPhoneNumber: 目标手机号码列表（仅支持中国大陆号码）
// 返回:
//   - error: 错误信息
func (c *SmsBaoClient) SendMessage(param map[string]string, targetPhoneNumber ...string) error {
	code, ok := param["code"]
	if !ok {
		return fmt.Errorf("missing parameter: code")
	}

	if len(targetPhoneNumber) == 0 {
		return fmt.Errorf("missing parameter: targetPhoneNumber")
	}

	smsContent := url.QueryEscape("【" + c.sign + "】" + fmt.Sprintf(c.template, code))
	for _, mobile := range targetPhoneNumber {
		if strings.HasPrefix(mobile, "+86") {
			mobile = mobile[3:]
		} else if strings.HasPrefix(mobile, "+") {
			return fmt.Errorf("unsupported country code")
		}
		// 短信宝API接口地址
		url := fmt.Sprintf("https://api.smsbao.com/sms?u=%s&p=%s&g=%s&m=%s&c=%s", c.username, c.apikey, c.goodsid, mobile, smsContent)

		client := &http.Client{}
		req, _ := http.NewRequest("GET", url, nil)
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		switch string(body) {
		case "30":
			return fmt.Errorf("password error")
		case "40":
			return fmt.Errorf("account not exist")
		case "41":
			return fmt.Errorf("overdue account")
		case "43":
			return fmt.Errorf("IP address limit")
		case "50":
			return fmt.Errorf("content contain forbidden words")
		case "51":
			return fmt.Errorf("phone number incorrect")
		}
	}

	return nil
}
