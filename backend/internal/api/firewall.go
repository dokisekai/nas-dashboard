package api

import (
	"nas-dashboard/internal/database"
	"nas-dashboard/internal/models"
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FirewallAPI struct {
	db *gorm.DB
}

func NewFirewallAPI() *FirewallAPI {
	return &FirewallAPI{
		db: database.GetDB(),
	}
}

var firewallAPI = NewFirewallAPI()

func GetFirewallRules(c *gin.Context) {
	var rules []models.FirewallRule
	if err := firewallAPI.db.Order("`order` asc, id desc").Find(&rules).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rules)
}

func CreateFirewallRule(c *gin.Context) {
	var rule models.FirewallRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := firewallAPI.db.Create(&rule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, rule)
}

func UpdateFirewallRule(c *gin.Context) {
	id := c.Param("id")
	var rule models.FirewallRule
	if err := firewallAPI.db.First(&rule, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rule not found"})
		return
	}

	if err := c.ShouldBindJSON(&rule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := firewallAPI.db.Save(&rule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rule)
}

func DeleteFirewallRule(c *gin.Context) {
	id := c.Param("id")
	if err := firewallAPI.db.Delete(&models.FirewallRule{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rule deleted"})
}

func ApplyFirewallRules(c *gin.Context) {
	var rules []models.FirewallRule
	if err := firewallAPI.db.Where("enabled = ?", true).Order("`order` asc").Find(&rules).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch rules: " + err.Error()})
		return
	}

	// 异步执行防火墙应用，避免API超时
	go func() {
		// 确保 SSH 端口始终开启，防止锁定自己
		exec.Command("ufw", "allow", "22/tcp").Run()
		// 确保面板端口开启
		exec.Command("ufw", "allow", "8888/tcp").Run()
		exec.Command("ufw", "allow", "5173/tcp").Run()

		// 应用规则
		for _, rule := range rules {
			args := []string{rule.Action}
			if rule.SourceIP != "" && rule.SourceIP != "any" {
				args = append(args, "from", rule.SourceIP)
			}
			if rule.Port != "" {
				portProto := rule.Port
				if rule.Protocol != "both" && rule.Protocol != "" {
					portProto += "/" + rule.Protocol
				}
				args = append(args, "to", "any", "port", portProto)
			}
			exec.Command("ufw", args...).Run()
		}

		// 启用防火墙
		exec.Command("ufw", "--force", "enable").Run()
	}()

	c.JSON(http.StatusOK, gin.H{"message": "防火墙规则已提交并在后台应用"})
}

// FirewallConfig 防火墙配置类型
type FirewallConfig struct {
	Enabled       bool     `json:"enabled"`
	DefaultPolicy string   `json:"defaultPolicy"` // accept, drop
	AllowedPorts  []string `json:"allowedPorts"`
	Logging       bool     `json:"logging"`
	ICMP          bool     `json:"icmp"`
}

func GetFirewallConfig(c *gin.Context) {
	config := FirewallConfig{
		Enabled:       true,
		DefaultPolicy: "drop",
		AllowedPorts:  []string{"22", "80", "443"},
		Logging:       false,
		ICMP:          true,
	}
	c.JSON(http.StatusOK, config)
}

func SetFirewallConfig(c *gin.Context) {
	var config FirewallConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 应用配置到系统
	go func() {
		if config.Enabled {
			// 确保 SSH 端口始终开启
			exec.Command("ufw", "allow", "22/tcp").Run()
			exec.Command("ufw", "--force", "enable").Run()
		} else {
			exec.Command("ufw", "disable").Run()
		}

		if config.DefaultPolicy == "drop" || config.DefaultPolicy == "deny" {
			exec.Command("ufw", "default", "deny", "incoming").Run()
		} else if config.DefaultPolicy == "accept" || config.DefaultPolicy == "allow" {
			exec.Command("ufw", "default", "allow", "incoming").Run()
		}

		if config.Logging {
			exec.Command("ufw", "logging", "on").Run()
		} else {
			exec.Command("ufw", "logging", "off").Run()
		}
	}()

	c.JSON(http.StatusOK, gin.H{
		"message": "防火墙配置已更新",
		"config":  config,
	})
}
