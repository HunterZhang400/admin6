package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

const (
	// AES-256 密钥长度 (32 字节)
	AES256KeyLength = 32
	// 标准 nonce 大小 (GCM 推荐)
	NonceSize = 12
	// 认证标签大小
	GCMTagSize = 16
)

var (
	ErrInvalidKeyLength   = errors.New("invalid key length, must be 32 bytes for AES-256")
	ErrInvalidCiphertext  = errors.New("invalid ciphertext format")
	ErrDecryptionFailed   = errors.New("decryption failed: authentication failed")
	ErrInvalidInput       = errors.New("invalid input")
	ErrInvalidNonceLength = errors.New("invalid nonce length")
)

// AESGCM 加密器结构体
type AESGCM struct {
	key []byte
}

// NewAESGCM 创建一个新的 AES-GCM 加密器
// 参数 key: 32字节的AES-256密钥
func NewAESGCM(key []byte) (*AESGCM, error) {
	if len(key) != AES256KeyLength {
		return nil, ErrInvalidKeyLength
	}
	return &AESGCM{key: key}, nil
}

// GenerateRandomKey 生成一个安全的随机密钥 (32字节)
func GenerateRandomKey() ([]byte, error) {
	key := make([]byte, AES256KeyLength)
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("failed to generate random key: %w", err)
	}
	return key, nil
}

// Encrypt 加密数据并返回 base64 编码的密文字符串
func (a *AESGCM) Encrypt(plaintext []byte) (string, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// 生成随机 nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// 加密数据
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	// 返回 base64 编码的字符串
	urlEncoded := base64.URLEncoding.EncodeToString(ciphertext)
	return urlEncoded, nil
}

// Decrypt 从 base64 编码的密文字符串解密数据
func (a *AESGCM) Decrypt(encodedCiphertext string) ([]byte, error) {
	// 解码 base64
	ciphertext, err := base64.URLEncoding.DecodeString(encodedCiphertext)
	if err != nil {
		return nil, fmt.Errorf("base64 decode failed: %w", err)
	}

	return a.DecryptBytes(ciphertext)
}

// DecryptBytes 直接从字节切片解密数据
func (a *AESGCM) DecryptBytes(ciphertext []byte) ([]byte, error) {
	if len(ciphertext) < NonceSize+GCMTagSize {
		return nil, ErrInvalidCiphertext
	}

	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// 分离 nonce 和实际密文
	nonce := ciphertext[:gcm.NonceSize()]
	ciphertext = ciphertext[gcm.NonceSize():]

	// 解密数据
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrDecryptionFailed, err)
	}

	return plaintext, nil
}

// EncryptString 加密字符串并返回 base64 编码的密文
func (a *AESGCM) EncryptString(plaintext []byte) (string, error) {
	return a.Encrypt(plaintext)
}

// DecryptString 解密 base64 编码的密文并返回字符串
func (a *AESGCM) DecryptString(encodedCiphertext string) (string, error) {

	plaintext, err := a.Decrypt(encodedCiphertext)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

// KeyToBase64 将密钥转换为 base64 字符串
func (a *AESGCM) KeyToBase64() string {
	return base64.URLEncoding.EncodeToString(a.key)
}

func (a *AESGCM) KeyToString() string {
	return base64.StdEncoding.EncodeToString(a.key)
}

// NewAESGCMFromBase64Key 从 base64 字符串创建 AESGCM 实例
func NewAESGCMFromBase64Key(base64Key string) (*AESGCM, error) {
	key, err := base64.URLEncoding.DecodeString(base64Key)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 key: %w", err)
	}
	return NewAESGCM(key)
}

// GenerateSecureKey 生成并返回 base64 编码的安全密钥
func GenerateSecureKey() (string, error) {
	key, err := GenerateRandomKey()
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(key), nil
}

// IsValidKey 检查是否是有效的 AES-256 密钥
func IsValidKey(key []byte) bool {
	return len(key) == AES256KeyLength
}

// IsValidBase64Key 检查 base64 字符串是否是有效的 AES-256 密钥
func IsValidBase64Key(base64Key string) bool {
	key, err := base64.URLEncoding.DecodeString(base64Key)
	if err != nil {
		return false
	}
	return IsValidKey(key)
}
