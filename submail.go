// Package sms SUBMAIL短信服务实现
package sms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

// SubmailClient SUBMAIL短信客户端
// 封装SUBMAIL短信API调用
type SubmailClient struct {
	api       string // API地址
	appid     string // 应用ID
	signature string // 签名
	project   string // 项目标识
}

// SubmailResult SUBMAIL响应结果结构体
type SubmailResult struct {
	Status string `json:"status"` // 状态
	Code   int    `json:"code"`   // 状态码
	Msg    string `json:"msg"`    // 消息
}

// buildSubmailPostdata 构建SUBMAIL POST数据
// 参数:
//   - param: 短信模板参数
//   - appid: 应用ID
//   - signature: 签名
//   - project: 项目标识
//   - targetPhoneNumber: 目标手机号码列表
// 返回:
//   - map[string]string: POST数据
//   - error: 错误信息
func buildSubmailPostdata(param map[string]string, appid string, signature string, project string, targetPhoneNumber []string) (map[string]string, error) {
	multi := make([]map[string]interface{}, 0, 32)

	for _, phoneNumber := range targetPhoneNumber[0:] {
		multi = append(multi, map[string]interface{}{
			"to":   phoneNumber,
			"vars": param,
		})
	}

	m, err := json.Marshal(multi)
	if err != nil {
		return nil, err
	}

	postdata := make(map[string]string)
	postdata["appid"] = appid
	postdata["signature"] = signature
	postdata["project"] = project
	postdata["multi"] = string(m)
	return postdata, nil
}

// GetSubmailClient 创建SUBMAIL短信客户端
// 参数:
//   - appid: 应用ID
//   - signature: 签名
//   - project: 项目标识
// 返回:
//   - *SubmailClient: SUBMAIL短信客户端实例
//   - error: 错误信息
func GetSubmailClient(appid string, signature string, project string) (*SubmailClient, error) {
	submailClient := &SubmailClient{
		api:       "https://api-v4.mysubmail.com/sms/multixsend",
		appid:     appid,
		signature: signature,
		project:   project,
	}
	return submailClient, nil
}

// SendMessage 发送短信
// 参数:
//   - param: 短信模板参数
//   - targetPhoneNumber: 目标手机号码列表
// 返回:
//   - error: 错误信息
func (c *SubmailClient) SendMessage(param map[string]string, targetPhoneNumber ...string) error {
	postdata, err := buildSubmailPostdata(param, c.appid, c.signature, c.project, targetPhoneNumber)
	if err != nil {
		return err
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for key, val := range postdata {
		err = writer.WriteField(key, val)
		if err != nil {
			return err
		}
	}

	contentType := writer.FormDataContentType()
	err = writer.Close()
	if err != nil {
		return err
	}

	resp, err := http.Post(c.api, contentType, body)
	if err != nil {
		return err
	}

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return handleSubmailResult(result)
}

// handleSubmailResult 处理SUBMAIL响应结果
// 参数:
//   - result: 响应结果字节数组
// 返回:
//   - error: 错误信息
func handleSubmailResult(result []byte) error {
	var submailSuccessResult []SubmailResult
	err := json.Unmarshal(result, &submailSuccessResult)
	if err != nil {
		var submailErrorResult SubmailResult
		err = json.Unmarshal(result, &submailErrorResult)
		if err != nil {
			return err
		}

		if submailErrorResult.Msg != "" {
			return fmt.Errorf(submailErrorResult.Msg)
		}
	}

	errMsgs := []string{}
	for _, submailResult := range submailSuccessResult {
		if submailResult.Status != "success" {
			errMsg := fmt.Sprintf("%s, %d, %s", submailResult.Status, submailResult.Code, submailResult.Msg)
			errMsgs = append(errMsgs, errMsg)
		}
	}

	if len(errMsgs) > 0 {
		return fmt.Errorf(strings.Join(errMsgs, "|"))
	}

	return nil
}
