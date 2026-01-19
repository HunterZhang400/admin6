package crypto

import (
	"encoding/base64"
	"fmt"
	"log"
	"testing"
)

func TestAES(t *testing.T) {
	// 生成随机密钥
	key, err := GenerateRandomKey()
	if err != nil {
		log.Fatal("Failed to generate key:", err)
	}
	aes, err := NewAESGCM(key)
	fmt.Println(aes.KeyToBase64())
	keyStr := "9SiGMgQyUft5GAOAHW9yRJkj4UP-xVQYvNq1wpDbbRg="
	keyBytes, err := base64.URLEncoding.DecodeString(keyStr)
	if err != nil {
		log.Println("base64 url decode error:", err)
		return
	}
	// 创建加密器
	aes, err = NewAESGCM(keyBytes)
	if err != nil {
		log.Fatal("Failed to create AES-GCM:", err)
	}

	// 加密数据
	secret := "45789545"
	encrypted, err := aes.EncryptString([]byte(secret))
	if err != nil {
		log.Fatal("Encryption failed:", err)
	}

	fmt.Printf("Encrypted: %s\n", encrypted)

	// 解密数据
	decrypted, err := aes.DecryptString(encrypted)
	if err != nil {
		log.Fatal("Decryption failed:", err)
	}

	fmt.Printf("Decrypted: %s\n", decrypted)

	// 验证密钥有效性
	fmt.Println("Key valid:", IsValidKey(keyBytes))

	// 使用 base64 密钥
	base64Key := aes.KeyToBase64()
	fmt.Println("Base64 Key:", base64Key)

	// 从 base64 密钥创建新实例
	aes2, err := NewAESGCMFromBase64Key(base64Key)
	if err != nil {
		log.Fatal("Failed to create from base64 key:", err)
	}

	// 使用新实例解密
	decrypted2, err := aes2.DecryptString(encrypted)
	if err != nil {
		log.Fatal("Decryption with new instance failed:", err)
	}

	fmt.Printf("Decrypted with new instance: %s\n", decrypted2)
}
