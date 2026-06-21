package sso

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"sync"
)

// KeyManager 管理 OIDC id_token 签名所用的 RSA 密钥对。
//
// 为什么要存在：oauth4webapi（Immich / openid-client 等严格 OIDC 客户端使用）
// 强制要求 id_token 使用非对称签名（RS256），并通过 JWKS 端点拉取公钥来验签。
// HS256 这种对称签名会被这些客户端直接拒绝。
type KeyManager struct {
	mu        sync.RWMutex
	private   *rsa.PrivateKey
	keyID     string
	keyPath   string
}

// NewKeyManager 创建密钥管理器。如果 keyPath 指向的文件存在则加载，
// 否则生成新的 RSA-2048 密钥并以 PEM 格式写入磁盘，以便 JWKS 在重启后保持稳定。
func NewKeyManager(keyPath string) (*KeyManager, error) {
	if keyPath == "" {
		keyPath = "sso_id_token.key"
	}

	km := &KeyManager{keyPath: keyPath}

	if data, err := os.ReadFile(keyPath); err == nil {
		block, _ := pem.Decode(data)
		if block == nil {
			return nil, fmt.Errorf("invalid PEM in %s", keyPath)
		}
		key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			if anyKey, err2 := x509.ParsePKCS8PrivateKey(block.Bytes); err2 == nil {
				if rsaKey, ok := anyKey.(*rsa.PrivateKey); ok {
					key = rsaKey
				} else {
					return nil, fmt.Errorf("PKCS8 key is not RSA")
				}
			} else {
				return nil, fmt.Errorf("parse private key: %w", err)
			}
		}
		km.private = key
		km.keyID = computeKeyID(&key.PublicKey)
		return km, nil
	}

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("generate RSA key: %w", err)
	}

	der := x509.MarshalPKCS1PrivateKey(key)
	pemBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: der,
	})

	if err := os.WriteFile(keyPath, pemBytes, 0o600); err != nil {
		return nil, fmt.Errorf("persist key file %s: %w", keyPath, err)
	}

	km.private = key
	km.keyID = computeKeyID(&key.PublicKey)
	return km, nil
}

// PrivateKey 返回用于签名的 RSA 私钥。
func (km *KeyManager) PrivateKey() *rsa.PrivateKey {
	km.mu.RLock()
	defer km.mu.RUnlock()
	return km.private
}

// KeyID 返回 JWKS 中使用的 kid（公钥指纹）。
func (km *KeyManager) KeyID() string {
	km.mu.RLock()
	defer km.mu.RUnlock()
	return km.keyID
}

// JWK 将公钥序列化为 RFC 7517 / RFC 7518（JWK）格式。
// 这是 Immich / oauth4webapi 期望从 /jwks 返回的结构。
func (km *KeyManager) JWK() map[string]interface{} {
	km.mu.RLock()
	defer km.mu.RUnlock()

	pub := &km.private.PublicKey

	nBytes := pub.N.Bytes()
	eBytes := big.NewInt(int64(pub.E)).Bytes()

	return map[string]interface{}{
		"kty": "RSA",
		"use": "sig",
		"alg": "RS256",
		"kid": km.keyID,
		"n":   base64.RawURLEncoding.EncodeToString(nBytes),
		"e":   base64.RawURLEncoding.EncodeToString(eBytes),
	}
}

// computeKeyID 计算公钥的稳定标识符（RFC 7638 thumbprint 风格）。
// 使用 JWK 规范化 JSON 的 SHA-256，截取前 16 字节作 base64url。
func computeKeyID(pub *rsa.PublicKey) string {
	nBytes := pub.N.Bytes()
	eBytes := big.NewInt(int64(pub.E)).Bytes()

	canonical := map[string]interface{}{
		"e":   base64.RawURLEncoding.EncodeToString(eBytes),
		"kty": "RSA",
		"n":   base64.RawURLEncoding.EncodeToString(nBytes),
	}
	// 使用稳定的键序：Go 的 json.Marshal 对 map 是按 key 排序的
	raw, _ := json.Marshal(canonical)

	sum := sha256.Sum256(raw)
	return base64.RawURLEncoding.EncodeToString(sum[:16])
}
