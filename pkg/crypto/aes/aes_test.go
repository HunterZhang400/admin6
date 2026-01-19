package aes_test

import (
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"testing"
	"zzl/pkg/crypto/aes"
)

func TestAESKey(t *testing.T) {
	//key := []byte{29, 215, 184, 143, 13, 20, 23, 101, 78, 2, 241, 203, 153, 252, 72, 215, 80, 27, 213, 72, 113, 207, 162, 19, 244, 189, 82, 124, 138, 217, 219, 166}
	aesClient, err := aes.NewAESGCM(aes.DefaultCode)
	if err != nil {
		log.Fatal("Failed to create AES-GCM:", err)
	}
	plainText := "PiuV5x3058Aj44"
	enText, err := aesClient.EncryptUseStandardBase64([]byte(plainText))
	if err != nil {
		t.Error("Failed to decrypt plaintext:", err)
		return
	}
	decodeText, err := aesClient.DecryptUseStandardBase64(enText)
	fmt.Println("密文:", string(enText))

	fmt.Println("原文:", string(decodeText))
}

func TestAES(t *testing.T) {
	t.Parallel()
	// 生成随机密钥
	key, err := aes.GenerateRandomKey()
	if err != nil {
		log.Fatal("Failed to generate key:", err)
	}
	fmt.Println("Key:", key)
	// 创建加密器
	aesClient, err := aes.NewAESGCM(key)
	if err != nil {
		log.Fatal("Failed to create AES-GCM:", err)
	}

	// 加密数据
	secret := "This is a highly confidential message!"
	encrypted, err := aesClient.EncryptString([]byte(secret))
	if err != nil {
		log.Fatal("Encryption failed:", err)
	}

	fmt.Printf("Encrypted: %s\n", encrypted)

	// 解密数据
	decrypted, err := aesClient.DecryptString(encrypted)
	if err != nil {
		log.Fatal("Decryption failed:", err)
	}

	fmt.Printf("Decrypted: %s\n", decrypted)

	// 验证密钥有效性
	fmt.Println("Key valid:", aes.IsValidKey(key))

	// 使用 base64 密钥
	base64Key := aesClient.KeyToBase64()
	fmt.Println("Base64 Key:", base64Key)

	// 从 base64 密钥创建新实例
	aes2, err := aes.NewAESGCMFromBase64Key(base64Key)
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

// TestNewAESGCM 测试 NewAESGCM 函数
func TestNewAESGCM(t *testing.T) {
	// 测试有效的32字节密钥
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}

	client, err := aes.NewAESGCM(key)
	if err != nil {
		t.Fatalf("NewAESGCM failed with valid key: %v", err)
	}

	if client == nil {
		t.Fatal("NewAESGCM returned nil client")
	}
}

// TestNewAESGCMWithInvalidKey 测试 NewAESGCM 的错误处理
func TestNewAESGCMWithInvalidKey(t *testing.T) {
	// 测试无效长度的密钥
	invalidKeys := [][]byte{
		nil,              // nil key
		{},               // empty key
		{1, 2, 3},        // too short
		make([]byte, 16), // 16 bytes (too short)
		make([]byte, 24), // 24 bytes (too short)
		make([]byte, 33), // 33 bytes (too long)
	}

	for i, key := range invalidKeys {
		_, err := aes.NewAESGCM(key)
		if err == nil {
			t.Errorf("Expected error for invalid key %d, got nil", i)
		}
	}
}

// TestGenerateRandomKey 测试 GenerateRandomKey 函数
func TestGenerateRandomKey(t *testing.T) {
	key, err := aes.GenerateRandomKey()
	if err != nil {
		t.Fatalf("GenerateRandomKey failed: %v", err)
	}

	if len(key) != 32 {
		t.Errorf("GenerateRandomKey returned key of length %d, expected 32", len(key))
	}

	// 测试多次生成不同的密钥
	key2, err := aes.GenerateRandomKey()
	if err != nil {
		t.Fatalf("Second GenerateRandomKey failed: %v", err)
	}

	// 虽然理论上可能相同，但概率极低
	if string(key) == string(key2) {
		t.Log("Warning: Two random keys are identical (very unlikely)")
	}
}

// TestEncryptDecryptRoundTrip 测试加密解密的往返
func TestEncryptDecryptRoundTrip(t *testing.T) {
	key, err := aes.GenerateRandomKey()
	if err != nil {
		t.Fatalf("GenerateRandomKey failed: %v", err)
	}

	client, err := aes.NewAESGCM(key)
	if err != nil {
		t.Fatalf("NewAESGCM failed: %v", err)
	}

	testCases := []string{
		"",
		"Hello, World!",
		"This is a test message with special characters: !@#$%^&*()",
		"中文测试",
		strings.Repeat("A", 1000), // 长字符串
		"Multi-line\nmessage\twith\ttabs",
	}

	for _, plaintext := range testCases {
		// 测试 Encrypt/Decrypt
		encrypted, err := client.Encrypt([]byte(plaintext))
		if err != nil {
			t.Errorf("Encrypt failed for %q: %v", plaintext, err)
			continue
		}

		decrypted, err := client.Decrypt(encrypted)
		if err != nil {
			t.Errorf("Decrypt failed for %q: %v", encrypted, err)
			continue
		}

		if string(decrypted) != plaintext {
			t.Errorf("Round trip failed for %q: got %q", plaintext, string(decrypted))
		}
	}
}

// TestEncryptStringDecryptString 测试字符串加密解密
func TestEncryptStringDecryptString(t *testing.T) {
	key, err := aes.GenerateRandomKey()
	if err != nil {
		t.Fatalf("GenerateRandomKey failed: %v", err)
	}

	client, err := aes.NewAESGCM(key)
	if err != nil {
		t.Fatalf("NewAESGCM failed: %v", err)
	}

	testCases := []string{
		"",
		"Hello, World!",
		"This is a test message",
		"中文测试",
	}

	for _, plaintext := range testCases {
		// 测试 EncryptString/DecryptString
		encrypted, err := client.EncryptString([]byte(plaintext))
		if err != nil {
			t.Errorf("EncryptString failed for %q: %v", plaintext, err)
			continue
		}

		decrypted, err := client.DecryptString(encrypted)
		if err != nil {
			t.Errorf("DecryptString failed for %q: %v", encrypted, err)
			continue
		}

		if decrypted != plaintext {
			t.Errorf("Round trip failed for %q: got %q", plaintext, decrypted)
		}
	}
}

// TestEncryptUseStandardBase64DecryptUseStandardBase64 测试标准Base64加密解密
func TestEncryptUseStandardBase64DecryptUseStandardBase64(t *testing.T) {
	key, err := aes.GenerateRandomKey()
	if err != nil {
		t.Fatalf("GenerateRandomKey failed: %v", err)
	}

	client, err := aes.NewAESGCM(key)
	if err != nil {
		t.Fatalf("NewAESGCM failed: %v", err)
	}

	plaintext := "Test message for standard base64"

	// 测试标准Base64加密
	encrypted, err := client.EncryptUseStandardBase64([]byte(plaintext))
	if err != nil {
		t.Fatalf("EncryptUseStandardBase64 failed: %v", err)
	}

	// 验证是标准Base64格式（可能包含填充）
	// 注意：标准Base64不一定总是包含填充字符，取决于数据长度
	if len(encrypted) == 0 {
		t.Error("Standard base64 should not be empty")
	}

	// 测试标准Base64解密
	decrypted, err := client.DecryptUseStandardBase64(encrypted)
	if err != nil {
		t.Fatalf("DecryptUseStandardBase64 failed: %v", err)
	}

	if string(decrypted) != plaintext {
		t.Errorf("Round trip failed: got %q, expected %q", string(decrypted), plaintext)
	}
}

// TestDecryptWithInvalidInput 测试解密无效输入
func TestDecryptWithInvalidInput(t *testing.T) {
	key, err := aes.GenerateRandomKey()
	if err != nil {
		t.Fatalf("GenerateRandomKey failed: %v", err)
	}

	client, err := aes.NewAESGCM(key)
	if err != nil {
		t.Fatalf("NewAESGCM failed: %v", err)
	}

	// 测试无效的base64字符串
	_, err = client.Decrypt("invalid-base64!")
	if err == nil {
		t.Error("Expected error for invalid base64 input")
	}

	// 测试空字符串
	_, err = client.Decrypt("")
	if err == nil {
		t.Error("Expected error for empty input")
	}

	// 测试太短的密文
	shortCiphertext := base64.RawURLEncoding.EncodeToString([]byte("short"))
	_, err = client.Decrypt(shortCiphertext)
	if err == nil {
		t.Error("Expected error for too short ciphertext")
	}
}

// TestKeyToBase64 测试密钥转Base64
func TestKeyToBase64(t *testing.T) {
	key, err := aes.GenerateRandomKey()
	if err != nil {
		t.Fatalf("GenerateRandomKey failed: %v", err)
	}

	client, err := aes.NewAESGCM(key)
	if err != nil {
		t.Fatalf("NewAESGCM failed: %v", err)
	}

	base64Key := client.KeyToBase64()

	// 验证是有效的base64
	decodedKey, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		t.Errorf("KeyToBase64 returned invalid base64: %v", err)
	}

	if len(decodedKey) != 32 {
		t.Errorf("Decoded key length is %d, expected 32", len(decodedKey))
	}

	// 验证解码后的密钥与原始密钥相同
	if string(decodedKey) != string(key) {
		t.Error("Decoded key does not match original key")
	}
}

// TestNewAESGCMFromBase64Key 测试从Base64密钥创建实例
func TestNewAESGCMFromBase64Key(t *testing.T) {
	key, err := aes.GenerateRandomKey()
	if err != nil {
		t.Fatalf("GenerateRandomKey failed: %v", err)
	}

	base64Key := base64.StdEncoding.EncodeToString(key)

	client, err := aes.NewAESGCMFromBase64Key(base64Key)
	if err != nil {
		t.Fatalf("NewAESGCMFromBase64Key failed: %v", err)
	}

	if client == nil {
		t.Fatal("NewAESGCMFromBase64Key returned nil client")
	}

	// 测试加密解密
	plaintext := "Test message"
	encrypted, err := client.EncryptString([]byte(plaintext))
	if err != nil {
		t.Fatalf("EncryptString failed: %v", err)
	}

	decrypted, err := client.DecryptString(encrypted)
	if err != nil {
		t.Fatalf("DecryptString failed: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("Round trip failed: got %q, expected %q", decrypted, plaintext)
	}
}

// TestNewAESGCMFromBase64KeyWithInvalidInput 测试从无效Base64密钥创建实例
func TestNewAESGCMFromBase64KeyWithInvalidInput(t *testing.T) {
	// 测试无效的base64字符串
	_, err := aes.NewAESGCMFromBase64Key("invalid-base64!")
	if err == nil {
		t.Error("Expected error for invalid base64 key")
	}

	// 测试空字符串
	_, err = aes.NewAESGCMFromBase64Key("")
	if err == nil {
		t.Error("Expected error for empty base64 key")
	}

	// 测试错误长度的base64密钥
	shortKey := base64.StdEncoding.EncodeToString([]byte("short"))
	_, err = aes.NewAESGCMFromBase64Key(shortKey)
	if err == nil {
		t.Error("Expected error for short base64 key")
	}
}

// TestGenerateSecureKey 测试生成安全密钥
func TestGenerateSecureKey(t *testing.T) {
	base64Key, err := aes.GenerateSecureKey()
	if err != nil {
		t.Fatalf("GenerateSecureKey failed: %v", err)
	}

	// 验证是有效的base64
	key, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		t.Errorf("GenerateSecureKey returned invalid base64: %v", err)
	}

	if len(key) != 32 {
		t.Errorf("Generated key length is %d, expected 32", len(key))
	}

	// 测试多次生成不同的密钥
	base64Key2, err := aes.GenerateSecureKey()
	if err != nil {
		t.Fatalf("Second GenerateSecureKey failed: %v", err)
	}

	if base64Key == base64Key2 {
		t.Log("Warning: Two secure keys are identical (very unlikely)")
	}
}

// TestIsValidKey 测试密钥验证
func TestIsValidKey(t *testing.T) {
	// 测试有效密钥
	validKey := make([]byte, 32)
	if !aes.IsValidKey(validKey) {
		t.Error("Valid key should return true")
	}

	// 测试无效密钥
	invalidKeys := [][]byte{
		nil,
		{},
		{1, 2, 3},
		make([]byte, 16),
		make([]byte, 24),
		make([]byte, 33),
	}

	for i, key := range invalidKeys {
		if aes.IsValidKey(key) {
			t.Errorf("Invalid key %d should return false", i)
		}
	}
}

// TestIsValidBase64Key 测试Base64密钥验证
func TestIsValidBase64Key(t *testing.T) {
	// 测试有效的base64密钥
	validKey := base64.StdEncoding.EncodeToString(make([]byte, 32))
	if !aes.IsValidBase64Key(validKey) {
		t.Error("Valid base64 key should return true")
	}

	// 测试无效的base64密钥
	invalidKeys := []string{
		"",
		"invalid-base64!",
		base64.StdEncoding.EncodeToString([]byte("short")),
		"not-base64",
	}

	for i, key := range invalidKeys {
		if aes.IsValidBase64Key(key) {
			t.Errorf("Invalid base64 key %d should return false", i)
		}
	}
}

// TestDecryptBytes 测试直接字节解密
func TestDecryptBytes(t *testing.T) {
	key, err := aes.GenerateRandomKey()
	if err != nil {
		t.Fatalf("GenerateRandomKey failed: %v", err)
	}

	client, err := aes.NewAESGCM(key)
	if err != nil {
		t.Fatalf("NewAESGCM failed: %v", err)
	}

	plaintext := "Test message for byte decryption"

	// 先加密获取密文
	encrypted, err := client.Encrypt([]byte(plaintext))
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	// 解码base64获取原始字节
	ciphertext, err := base64.RawURLEncoding.DecodeString(encrypted)
	if err != nil {
		t.Fatalf("Base64 decode failed: %v", err)
	}

	// 直接解密字节
	decrypted, err := client.DecryptBytes(ciphertext)
	if err != nil {
		t.Fatalf("DecryptBytes failed: %v", err)
	}

	if string(decrypted) != plaintext {
		t.Errorf("DecryptBytes failed: got %q, expected %q", string(decrypted), plaintext)
	}
}

// TestDecryptBytesWithInvalidInput 测试字节解密的无效输入
func TestDecryptBytesWithInvalidInput(t *testing.T) {
	key, err := aes.GenerateRandomKey()
	if err != nil {
		t.Fatalf("GenerateRandomKey failed: %v", err)
	}

	client, err := aes.NewAESGCM(key)
	if err != nil {
		t.Fatalf("NewAESGCM failed: %v", err)
	}

	// 测试太短的密文
	shortCiphertext := []byte("short")
	_, err = client.DecryptBytes(shortCiphertext)
	if err == nil {
		t.Error("Expected error for too short ciphertext")
	}

	// 测试空密文
	_, err = client.DecryptBytes([]byte{})
	if err == nil {
		t.Error("Expected error for empty ciphertext")
	}
}

// TestEncryptionConsistency 测试加密一致性
func TestEncryptionConsistency(t *testing.T) {
	key, err := aes.GenerateRandomKey()
	if err != nil {
		t.Fatalf("GenerateRandomKey failed: %v", err)
	}

	client, err := aes.NewAESGCM(key)
	if err != nil {
		t.Fatalf("NewAESGCM failed: %v", err)
	}

	plaintext := "Consistency test message"

	// 多次加密应该产生不同的结果（因为nonce是随机的）
	encrypted1, err := client.Encrypt([]byte(plaintext))
	if err != nil {
		t.Fatalf("First encrypt failed: %v", err)
	}

	encrypted2, err := client.Encrypt([]byte(plaintext))
	if err != nil {
		t.Fatalf("Second encrypt failed: %v", err)
	}

	// 加密结果应该不同（因为nonce不同）
	if encrypted1 == encrypted2 {
		t.Error("Encryption should produce different results due to random nonce")
	}

	// 但解密结果应该相同
	decrypted1, err := client.Decrypt(encrypted1)
	if err != nil {
		t.Fatalf("First decrypt failed: %v", err)
	}

	decrypted2, err := client.Decrypt(encrypted2)
	if err != nil {
		t.Fatalf("Second decrypt failed: %v", err)
	}

	if string(decrypted1) != plaintext || string(decrypted2) != plaintext {
		t.Error("Both decryptions should produce the original plaintext")
	}
}

// TestLargeDataEncryption 测试大数据加密
func TestLargeDataEncryption(t *testing.T) {
	key, err := aes.GenerateRandomKey()
	if err != nil {
		t.Fatalf("GenerateRandomKey failed: %v", err)
	}

	client, err := aes.NewAESGCM(key)
	if err != nil {
		t.Fatalf("NewAESGCM failed: %v", err)
	}

	// 创建1MB的数据
	largeData := make([]byte, 1024*1024)
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}

	// 加密大数据
	encrypted, err := client.Encrypt(largeData)
	if err != nil {
		t.Fatalf("Encrypt large data failed: %v", err)
	}

	// 解密大数据
	decrypted, err := client.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Decrypt large data failed: %v", err)
	}

	if len(decrypted) != len(largeData) {
		t.Errorf("Decrypted data length %d != original length %d", len(decrypted), len(largeData))
	}

	if string(decrypted) != string(largeData) {
		t.Error("Decrypted large data does not match original")
	}
}
