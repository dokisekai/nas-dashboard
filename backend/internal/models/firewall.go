package models

// FirewallRule 防火墙规则模型
type FirewallRule struct {
	BaseModel
	Name        string `gorm:"not null" json:"name"`
	Action      string `gorm:"default:'allow'" json:"action"` // allow, deny
	Protocol    string `gorm:"default:'tcp'" json:"protocol"` // tcp, udp, both
	Port        string `json:"port"`                          // 端口号，可以是范围 如 "80", "20-21"
	SourceIP    string `gorm:"default:'any'" json:"sourceIp"` // 来源IP，如 "192.168.1.1", "192.168.1.0/24", "any"
	Description string `json:"description"`
	Enabled     bool   `gorm:"default:true" json:"enabled"`
	Order       int    `gorm:"default:0" json:"order"`
}
