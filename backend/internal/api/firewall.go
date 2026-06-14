package api

import (
	"nas-dashboard/internal/database"
	"nas-dashboard/internal/models"
	"net/http"

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
	// TODO: Implement actual system call to ufw/iptables
	c.JSON(http.StatusOK, gin.H{"message": "Firewall rules applied successfully"})
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

	// TODO: Apply firewall configuration to system
	c.JSON(http.StatusOK, gin.H{
		"message": "Firewall configuration updated successfully",
		"config":  config,
	})
}
