// Package sms Netgsm短信服务实现
package sms

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// NetgsmClient Netgsm短信客户端
// 封装Netgsm短信API调用
type NetgsmClient struct {
	accessId   string       // 访问ID
	accessKey  string       // 访问密钥
	sign       string       // 短信签名
	template   string       // 短信模板
	httpClient *http.Client // HTTP客户端
}

// NetgsmResponse Netgsm响应结构体
type NetgsmResponse struct {
	Code  string `xml:"main>code"`  // 响应代码
	JobId string `xml:"main>jobId"` // 任务ID
	Error string `xml:"main>error"` // 错误信息
}

// GetNetgsmClient 创建Netgsm短信客户端
// 参数:
//   - accessId: 访问ID
//   - accessKey: 访问密钥
//   - sign: 短信签名
//   - template: 短信模板
// 返回:
//   - *NetgsmClient: Netgsm短信客户端实例
//   - error: 错误信息
func GetNetgsmClient(accessId, accessKey, sign, template string) (*NetgsmClient, error) {
	return &NetgsmClient{
		accessId:   accessId,
		accessKey:  accessKey,
		sign:       sign,
		template:   template,
		httpClient: &http.Client{},
	}, nil
}

// SendMessage 发送短信
// 参数:
//   - param: 短信模板参数
//   - targetPhoneNumber: 目标手机号码列表
// 返回:
//   - error: 错误信息
func (c *NetgsmClient) SendMessage(param map[string]string, targetPhoneNumber ...string) error {
	if len(targetPhoneNumber) == 0 {
		return fmt.Errorf("missing parameter: targetPhoneNumber")
	}

	for _, phoneNumber := range targetPhoneNumber {
		data := fmt.Sprintf(`
<mainbody>
   <header>
       <usercode>%s</usercode>
       <password>%s</password>
       <msgheader>%s</msgheader>
   </header>
   <body>
       <msg>
           <![CDATA[%s]]>
       </msg>
       <no>%s</no>
   </body>
</mainbody>`, c.accessId, c.accessKey, c.sign, c.template, phoneNumber)

		headers := map[string]string{
			"Content-Type": "application/xml",
		}

		respBody, err := c.postXML("https://api.netgsm.com.tr/sms/send/otp", data, headers)
		if err != nil {
			return err
		}

		var netgsmResponse NetgsmResponse
		if err := xml.Unmarshal([]byte(respBody), &netgsmResponse); err != nil {
			return err
		}

		if netgsmResponse.Code != "0" {
			return errors.New(netgsmResponse.Error)
		}
	}
	return nil
}

// postXML 发送XML格式的POST请求
// 参数:
//   - url: 请求URL
//   - xmlData: XML数据
//   - headers: 请求头
// 返回:
//   - string: 响应内容
//   - error: 错误信息
func (c *NetgsmClient) postXML(url, xmlData string, headers map[string]string) (string, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(xmlData)))
	if err != nil {
		return "", err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respBody), nil
}
