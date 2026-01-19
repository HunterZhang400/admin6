package shaid

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"mini_web_v2/pkg/crypto"
	"strconv"
	"strings"
	"sync"
)

// map[int]string
var idToStringMap = sync.Map{}

// map[string]int
var stringToIDMap = sync.Map{}

func DecodeDynamicID(key string, id2 string) (int, error) {
	if v, exist := stringToIDMap.Load(id2); exist {
		return v.(int), nil
	}
	aes, err := crypto.NewAESGCMFromBase64Key(key)
	if err != nil {
		return 0, err
	}
	decodeID, err := base64.URLEncoding.DecodeString(id2)
	if err != nil {
		return 0, err
	}
	v, err := aes.DecryptString(string(decodeID))
	if err != nil {
		return 0, err
	}
	id1, err := strconv.Atoi(v)
	if err != nil {
		return 0, err
	}
	stringToIDMap.Store(id2, id1)
	return id1, nil
}

func Encode2DynamicID(key string, id int) (string, error) {
	if v, exist := idToStringMap.Load(id); exist {
		return v.(string), nil
	}
	aes, err := crypto.NewAESGCMFromBase64Key(key)
	if err != nil {
		return "", err
	}
	v, err := aes.EncryptString([]byte(strconv.Itoa(id)))
	if err != nil {
		return "", err
	}
	v = base64.URLEncoding.EncodeToString([]byte(v))
	idToStringMap.Store(id, v)
	return v, nil
}

var ID_Salt = []byte{2, 158, 40, 199, 222, 3, 13, 47, 55, 200, 39, 178, 8}

func GetSHAID(id int) string {
	s := sha256.New()
	buf := make([]byte, 8) // int64 固定 8 字节
	binary.BigEndian.PutUint64(buf, uint64(id))
	s.Write(ID_Salt)
	s.Write(buf)
	bs := s.Sum(nil)
	return SafeFilenameFromBase64(base64.URLEncoding.EncodeToString(bs))
}

// SafeFilenameFromBase64 生成安全的文件名
func SafeFilenameFromBase64(b64 string) string {
	// 1. 移除填充字符
	clean := []byte(strings.TrimRight(b64, "="))

	for i, b := range clean {
		switch b {
		case '-', '.', ' ':
			clean[i] = '_'
		}
	}
	return string(clean)
}
