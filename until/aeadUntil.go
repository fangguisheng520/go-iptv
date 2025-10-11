package until

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"golang.org/x/crypto/chacha20poly1305"
)

type ChaCha20 struct {
	Key []byte
}

// 创建新的实例，key长度必须为 32 字节
func NewChaCha20(key string) *ChaCha20 {
	return &ChaCha20{
		Key: []byte(key),
	}
}

// 加密，返回 URL-safe Base64 字符串
func (c *ChaCha20) Encrypt(plainText string) (string, error) {
	aead, err := chacha20poly1305.New(c.Key)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, chacha20poly1305.NonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherText := aead.Seal(nil, nonce, []byte(plainText), nil)

	// 将 nonce + cipherText 一起编码为 Base64
	full := append(nonce, cipherText...)
	return base64.URLEncoding.EncodeToString(full), nil
}

// 解密，输入 URL-safe Base64 字符串
func (c *ChaCha20) Decrypt(encoded string) (string, error) {
	data, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("base64 decode failed: %v", err)
	}

	aead, err := chacha20poly1305.New(c.Key)
	if err != nil {
		return "", err
	}

	if len(data) < chacha20poly1305.NonceSize {
		return "", fmt.Errorf("data too short")
	}

	nonce := data[:chacha20poly1305.NonceSize]
	cipherText := data[chacha20poly1305.NonceSize:]

	plainText, err := aead.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", fmt.Errorf("decrypt failed: %v", err)
	}

	return string(plainText), nil
}
