// Package sms OSON短信服务实现
package sms

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
)

// OsonClient OSON短信客户端
// 封装OSON短信API调用
type OsonClient struct {
	Endpoint         string // API端点
	SenderId         string // 发送方ID
	SecretAccessHash string // 访问密钥哈希
	Sign             string // 短信签名
	Message          string // 短信内容
}

// OsonResponse OSON响应结构体
type OsonResponse struct {
	Status        string    `json:"status"`          // 状态 (ok)
	Timestamp     time.Time `json:"timestamp"`       // 时间戳 (2017-07-07 16:58:12)
	TxnId         string    `json:"txn_id"`          // 交易ID
	MsgId         uint      `json:"msg_id"`          // 消息ID (40127)
	SmscMsgId     string    `json:"smsc_msg_id"`     // SMSC消息ID
	SmscMsgStatus string    `json:"smsc_msg_status"` // SMSC消息状态
	SmscMsgParts  string    `json:"smsc_msg_parts"`  // SMSC消息部分
}

// GetOsonClient 创建OSON短信客户端
// 参数:
//   - senderId: 发送方ID
//   - secretAccessHash: 访问密钥哈希
//   - sign: 短信签名
//   - message: 短信内容
//
// 返回:
//   - *OsonClient: OSON短信客户端实例
//   - error: 错误信息
func GetOsonClient(senderId, secretAccessHash, sign, message string) (*OsonClient, error) {
	return &OsonClient{
		Endpoint:         "https://api.osonsms.com/sendsms_v1.php",
		SenderId:         senderId,
		SecretAccessHash: secretAccessHash,
		Sign:             sign,
		Message:          message,
	}, nil
}

// SendMessage 发送短信
// 参数:
//   - param: 短信模板参数
//   - targetPhoneNumber: 目标手机号码列表
//
// 返回:
//   - error: 错误信息
func (c *OsonClient) SendMessage(param map[string]string, targetPhoneNumber ...string) (err error) {
	// 初始化HTTP客户端用于向短信中心发送请求
	// 设置25+秒的超时时间以确保短信中心的响应已被处理
	client := &http.Client{
		Timeout: 20 * time.Second,
	}

	if c.Message == "" {
		c.Message = fmt.Sprintf("Hello. Your authorization code: %s", param["code"])
	} else {
		c.Message += param["code"]
	}

	txnId := uuid.NewString()
	buildStrHash := strings.Join([]string{txnId, c.SenderId, c.Sign, targetPhoneNumber[0], c.SecretAccessHash}, ";")

	hash := sha256.New()
	hash.Write([]byte(buildStrHash))
	bs := hash.Sum(nil)
	strHash := fmt.Sprintf("%x", bs)

	urlLink, err := url.Parse(c.Endpoint)
	if err != nil {
		return
	}

	urlParams := url.Values{}
	urlParams.Add("from", c.Sign)
	urlParams.Add("phone_number", targetPhoneNumber[0])
	urlParams.Add("msg", c.Message)
	urlParams.Add("str_hash", strHash)
	urlParams.Add("txn_id", txnId)
	urlParams.Add("login", c.SenderId)

	urlLink.RawQuery = urlParams.Encode()

	request, err := http.NewRequest(http.MethodGet, urlLink.String(), nil)
	if err != nil {
		return
	}

	resp, err := client.Do(request)
	if err != nil {
		return
	}

	resultBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result OsonResponse
	if err = json.Unmarshal(resultBytes, &result); err != nil {
		return
	}

	if result.Status != "ok" {
		return fmt.Errorf("sms service returned error status not 200: Status Code: %d Error: %s", resp.StatusCode, string(resultBytes))
	}

	return
}
