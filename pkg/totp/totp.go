package totp

import (
	"encoding/base32"
	"fmt"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

const (
	Issuer = "zzl"
)

// TOTPService TOTP服务
type TOTPService struct{}

// NewTOTPService 创建TOTP服务实例
func NewTOTPService() *TOTPService {
	return &TOTPService{}
}

// GenerateSecret 为用户生成TOTP密钥
func (s *TOTPService) GenerateSecret(userID uint, accountName string) (*otp.Key, error) {
	// 生成TOTP密钥
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      Issuer,
		AccountName: accountName,
		Period:      30,
		Digits:      otp.DigitsSix,
		Algorithm:   otp.AlgorithmSHA1,
	})
	if err != nil {
		return nil, fmt.Errorf("生成TOTP密钥失败: %v", err)
	}

	return key, nil
}

// GenerateQRCodeURL 生成二维码URL
func (s *TOTPService) GenerateQRCodeURL(secret string, accountName string) string {
	// 构建TOTP URL，使用简短的格式
	url := fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s",
		Issuer,
		accountName,
		secret,
		Issuer,
	)
	return url
}

// ValidateCode 验证TOTP验证码
func (s *TOTPService) ValidateCode(code string, secret string) bool {
	// 验证验证码
	valid := totp.Validate(code, secret)
	return valid
}

// ValidateCodeWithWindow 验证TOTP验证码（允许时间窗口）
func (s *TOTPService) ValidateCodeWithWindow(code string, secret string, window int) bool {
	// 验证验证码，允许时间窗口
	valid, err := totp.ValidateCustom(code, secret, time.Now(), totp.ValidateOpts{
		Period:    30,
		Skew:      uint(window), // 允许前后几个周期
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
	if err != nil {
		fmt.Println("ValidateCodeWithWindow", err)
		return false
	}
	return valid
}

// GenerateQRCodeDataURL 生成二维码数据URL（用于前端显示）
func (s *TOTPService) GenerateQRCodeDataURL(secret string, accountName string) (string, error) {
	// 生成二维码URL
	qrURL := s.GenerateQRCodeURL(secret, accountName)

	// 这里可以使用第三方库生成二维码图片
	// 为了简化，我们返回一个简单的数据URL
	// 实际项目中可以使用 github.com/skip2/go-qrcode 等库

	// 简单的二维码数据URL（实际应该生成真正的二维码图片）
	qrDataURL := fmt.Sprintf("data:image/svg+xml;base64,%s",
		base32.StdEncoding.EncodeToString([]byte(qrURL)))

	return qrDataURL, nil
}

// GetCurrentCode 获取当前时间窗口的验证码（用于测试）
func (s *TOTPService) GetCurrentCode(secret string) (string, error) {
	code, err := totp.GenerateCode(secret, time.Now())
	if err != nil {
		return "", fmt.Errorf("生成当前验证码失败: %v", err)
	}
	return code, nil
}
