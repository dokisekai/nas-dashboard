package models

import (
	"time"
	"gorm.io/gorm"
)

// OAuthClient OAuth客户端模型
type OAuthClient struct {
	ID           uint      `json:"id" gorm:"primarykey"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Name         string    `json:"name" gorm:"size:100;not null"`
	ClientID     string    `json:"client_id" gorm:"size:100;uniqueIndex;not null"`
	ClientSecret string    `json:"client_secret" gorm:"size:255;not null"`
	RedirectURIs string    `json:"redirect_uris" gorm:"type:text;not null"` // JSON数组
	GrantTypes   string    `json:"grant_types" gorm:"type:text;not null"`    // JSON数组
	Scopes       string    `json:"scopes" gorm:"type:text;not null"`         // JSON数组
	Status       string    `json:"status" gorm:"size:20;default:active"`
}

// TableName 指定表名
func (OAuthClient) TableName() string {
	return "oauth_clients"
}

// BeforeCreate 创建前钩子
func (c *OAuthClient) BeforeCreate(tx *gorm.DB) error {
	if c.Status == "" {
		c.Status = "active"
	}
	return nil
}

// OAuthToken OAuth令牌模型
type OAuthToken struct {
	ID           uint      `json:"id" gorm:"primarykey"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	UserID       uint      `json:"user_id" gorm:"not null;index"`
	ClientID     string    `json:"client_id" gorm:"size:100;not null;index"`
	AccessToken  string    `json:"access_token" gorm:"type:text;not null"`
	RefreshToken string    `json:"refresh_token" gorm:"type:text"`
	TokenType    string    `json:"token_type" gorm:"size:20;default:Bearer"`
	ExpiresAt    time.Time `json:"expires_at"`
	RevokedAt    *time.Time `json:"revoked_at"`
}

// TableName 指定表名
func (OAuthToken) TableName() string {
	return "oauth_tokens"
}

// OAuthAuthorization OAuth授权记录
type OAuthAuthorization struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	ClientID  string    `json:"client_id" gorm:"size:100;not null;index"`
	Scopes    string    `json:"scopes" gorm:"type:text;not null"` // JSON数组
	Code      string    `json:"code" gorm:"size:100"`
	ExpiresAt time.Time `json:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at"`
}

// TableName 指定表名
func (OAuthAuthorization) TableName() string {
	return "oauth_authorizations"
}

// IsExpired 检查授权是否过期
func (a *OAuthAuthorization) IsExpired() bool {
	return time.Now().After(a.ExpiresAt)
}

// IsRevoked 检查授权是否被撤销
func (a *OAuthAuthorization) IsRevoked() bool {
	return a.RevokedAt != nil && a.RevokedAt.Before(time.Now())
}

// IsValid 检查授权是否有效
func (a *OAuthAuthorization) IsValid() bool {
	return !a.IsExpired() && !a.IsRevoked()
}