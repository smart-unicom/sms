# 📱 su-sms

一个功能强大的Go语言短信发送组件库，支持多家主流短信服务提供商的统一接口调用。

## ✨ 特性

- 🔌 **统一接口**: 提供一致的API接口，支持多家短信服务提供商
- 🌍 **多服务商支持**: 支持20+家国内外主流短信服务提供商
- 🚀 **简单易用**: 简洁的API设计，快速集成到项目中
- 🛡️ **类型安全**: 完整的Go类型定义和错误处理
- 📝 **完整文档**: 详细的中文注释和使用示例
- 🧪 **测试友好**: 内置Mock客户端，方便单元测试

## 🏢 支持的服务提供商

### 国内服务商
- 🇨🇳 **阿里云短信** (Aliyun SMS)
- 🇨🇳 **腾讯云短信** (Tencent Cloud SMS)
- 🇨🇳 **百度云短信** (Baidu Cloud SMS)
- 🇨🇳 **华为云短信** (Huawei Cloud SMS)
- 🇨🇳 **火山引擎短信** (Volc Engine SMS)
- 🇨🇳 **UCloud短信** (UCloud SMS)
- 🇨🇳 **短信宝** (SmsBao SMS)
- 🇨🇳 **互亿无线** (Huyi SMS)
- 🇨🇳 **SUBMAIL**
- 🇨🇳 **GCCPAY SMS**

### 国际服务商
- 🌍 **Twilio**
- 🌍 **Amazon SNS**
- 🌍 **Microsoft Azure ACS**
- 🌍 **Infobip**
- 🌍 **Msg91**
- 🌍 **Netgsm**
- 🌍 **OSON SMS**
- 🌍 **Uni SMS**

### 测试工具
- 🧪 **Mock SMS** (用于测试环境)

## 📦 安装

```bash
go get github.com/smart-unicom/su-sms
```

## 🚀 快速开始

### 基本用法

```go
package main

import (
    "fmt"
    "github.com/smart-unicom/su-sms"
)

func main() {
    // 创建短信客户端
    client, err := sms.NewSmsProvider(
        sms.SMS_ALIYUN,     // 服务提供商
        "your-access-id",   // 访问ID
        "your-access-key",  // 访问密钥
        "your-sign-name",   // 短信签名
        "your-template-id", // 短信模板ID
    )
    if err != nil {
        panic(err)
    }

    // 发送短信
    params := map[string]string{
        "code": "123456", // 验证码
    }
    
    err = client.SendMessage(params, "+8613800138000")
    if err != nil {
        fmt.Printf("发送失败: %v\n", err)
        return
    }
    
    fmt.Println("短信发送成功！")
}
```

## 📖 详细使用示例

### 阿里云短信

```go
package main

import (
    "github.com/smart-unicom/su-sms"
)

func main() {
    client, err := sms.NewSmsProvider(
        sms.SMS_ALIYUN,
        "LTAI5tCExxxxxxxxxxxxxx",    // 访问密钥ID
        "2alE91xxxxxxxxxxxxxxxxxxxxxx", // 访问密钥
        "测试签名",                      // 短信签名
        "SMS_123456789",              // 短信模板CODE
    )
    if err != nil {
        panic(err)
    }

    params := map[string]string{
        "code": "888888",
    }
    
    err = client.SendMessage(params, "+8613800138000")
    if err != nil {
        panic(err)
    }
}
```

### 腾讯云短信

```go
package main

import (
    "github.com/smart-unicom/su-sms"
)

func main() {
    client, err := sms.NewSmsProvider(
        sms.SMS_TENCENT,
        "AKIDxxxxxxxxxxxxxxxxxxxxxx",   // 密钥ID
        "2E7Axxxxxxxxxxxxxxxxxxxxxx",   // 密钥
        "测试签名",                        // 短信签名
        "123456",                       // 短信模板ID
        "1400123456",                   // 短信应用ID
    )
    if err != nil {
        panic(err)
    }

    params := map[string]string{
        "0": "888888", // 腾讯云使用数字索引作为参数名
    }
    
    err = client.SendMessage(params, "+8613800138000")
    if err != nil {
        panic(err)
    }
}
```

### Twilio短信

```go
package main

import (
    "github.com/smart-unicom/su-sms"
)

func main() {
    client, err := sms.NewSmsProvider(
        sms.SMS_TWILIO,
        "ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", // 账户SID
        "your_auth_token",                    // 认证令牌
        "",                                    // 签名（Twilio不需要）
        "您的验证码是: %s",                      // 短信模板
    )
    if err != nil {
        panic(err)
    }

    params := map[string]string{
        "code": "888888",
    }
    
    // 注意：Twilio需要发送方号码作为第一个参数
    err = client.SendMessage(params, "+1234567890", "+8613800138000")
    if err != nil {
        panic(err)
    }
}
```

### 华为云短信

```go
package main

import (
    "github.com/smart-unicom/su-sms"
)

func main() {
    client, err := sms.NewSmsProvider(
        sms.SMS_HUAWEI,
        "your-app-key",                     // 应用密钥
        "your-app-secret",                  // 应用秘密
        "测试签名",                           // 短信签名
        "8ff55eac1d0b478ab3c06c3c6a492300", // 短信模板ID
        "https://smsapi.cn-north-4.myhuaweicloud.com", // API地址
        "+8610690000",                      // 发送方号码
    )
    if err != nil {
        panic(err)
    }

    params := map[string]string{
        "code": "888888",
    }
    
    err = client.SendMessage(params, "+8613800138000")
    if err != nil {
        panic(err)
    }
}
```

### Mock客户端（测试用）

```go
package main

import (
    "github.com/smart-unicom/su-sms"
)

func main() {
    // Mock客户端不会实际发送短信，适用于测试环境
    client, err := sms.NewSmsProvider(
        sms.SMS_MOCK,
        "mock-access-id",
        "mock-access-key",
        "mock-sign",
        "mock-template",
    )
    if err != nil {
        panic(err)
    }

    params := map[string]string{
        "code": "888888",
    }
    
    // 总是返回成功，不会实际发送短信
    err = client.SendMessage(params, "+8613800138000")
    if err != nil {
        panic(err)
    }
    
    fmt.Println("Mock短信发送成功（未实际发送）")
}
```

## 🔧 API参考

### 创建客户端

```go
func NewSmsProvider(provider string, accessId string, accessKey string, sign string, template string, other ...string) (SmsProvider, error)
```

**参数说明：**
- `provider`: 服务提供商类型（使用预定义常量）
- `accessId`: 访问ID/用户名
- `accessKey`: 访问密钥/密码
- `sign`: 短信签名
- `template`: 短信模板ID或内容
- `other`: 其他参数（根据不同服务商而定）

### 发送短信

```go
SendMessage(param map[string]string, targetPhoneNumber ...string) error
```

**参数说明：**
- `param`: 短信模板参数
- `targetPhoneNumber`: 目标手机号码列表

### 服务提供商常量

```go
const (
    SMS_ALIYUN   = "Aliyun_SMS"        // 阿里云短信
    SMS_TENCENT  = "Tencent_Cloud_SMS" // 腾讯云短信
    SMS_TWILIO   = "Twilio_SMS"        // Twilio短信
    SMS_AMAZON   = "Amazon_SNS"        // 亚马逊SNS
    SMS_AZURE    = "Azure_ACS"         // 微软Azure
    SMS_HUAWEI   = "Huawei_Cloud_SMS"  // 华为云短信
    SMS_MOCK     = "Mock SMS"          // 模拟短信
    // ... 更多常量
)
```

## 🧪 测试

运行所有测试：

```bash
go test -v ./...
```

运行特定测试：

```bash
go test -v -run TestMockSMS
```

## 🤝 贡献

欢迎提交Issue和Pull Request来帮助改进这个项目！

### 贡献指南

1. Fork 这个仓库
2. 创建你的特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交你的修改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开一个 Pull Request

### 代码规范

请遵循项目中的[代码规范](代码规范.md)，确保：
- 所有注释使用中文
- 遵循Go语言命名约定
- 添加适当的错误处理
- 编写相应的测试用例

## 📄 许可证

本项目基于 [Apache 2.0 许可证](LICENSE) 开源。

## 🙏 致谢

感谢所有为这个项目做出贡献的开发者！

---

如果这个项目对你有帮助，请给我们一个 ⭐️！
