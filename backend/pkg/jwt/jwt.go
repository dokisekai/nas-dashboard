package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrTokenExpired = errors.New("token expired")
	secretKey       = []byte("nas-dashboard-secret-key-change-in-production")
)

// Claims JWT 声明
type Claims struct {
	UserID   uint   `json:"userId"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT token (默认 24 小时有效期)
func GenerateToken(username string) (string, error) {
	return GenerateTokenWithExpiry(username, 24*time.Hour)
}

// GenerateTokenWithExpiry 生成指定有效期的 JWT token
func GenerateTokenWithExpiry(username string, expiry time.Duration) (string, error) {
	// 为admin用户设置角色
	role := "user"
	if username == "admin" {
		role = "admin"
	}

	claims := Claims{
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// ParseToken 解析 JWT token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// ValidateToken 验证 token 是否有效
func ValidateToken(tokenString string) bool {
	_, err := ParseToken(tokenString)
	return err == nil
}

// GetUsernameFromToken 从 token 中提取用户名
func GetUsernameFromToken(tokenString string) (string, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.Username, nil
}
