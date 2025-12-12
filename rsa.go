package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
)

// LoadPrivateKey 从PEM文件加载私钥，只支持Java PKCS#8格式
func LoadPrivateKey(filename string) (*rsa.PrivateKey, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(file)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	// 只支持PKCS#8格式解析（Java格式）
	pkcs8Key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("无法解析Java私钥: %w", err)
	}

	// 验证是RSA私钥
	rsaKey, ok := pkcs8Key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("不是RSA私钥")
	}

	return rsaKey, nil
}

// RSA_PSS_Sign 使用RSA-PSS算法对数据进行签名
func RSA_PSS_Sign(privateKey *rsa.PrivateKey, data []byte) (string, error) {
	hashed := sha256.Sum256(data)

	// 设置盐长度，使用与哈希长度相同的值（32字节，SHA256）
	pssOptions := &rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthEqualsHash,
	}

	signature, err := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, hashed[:], pssOptions)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}
