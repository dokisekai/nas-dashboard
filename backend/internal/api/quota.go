package api

import (
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"nas-dashboard/internal/models"
	"gorm.io/gorm"
)

// QuotaAPI 配额API处理器
type QuotaAPI struct {
	DB *gorm.DB
}

// NewQuotaAPI 创建配额API
func NewQuotaAPI(db *gorm.DB) *QuotaAPI {
	return &QuotaAPI{DB: db}
}

// GetUserQuota 获取用户配额
func (api *QuotaAPI) GetUserQuota(c *gin.Context) {
	username := c.Param("username")

	var user models.User
	if result := api.DB.Where("username = ?", username).First(&user); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var quotas []models.UserQuota
	if result := api.DB.Where("user_id = ?", user.ID).Find(&quotas); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"quotas": quotas})
}

// SetUserQuota 设置用户配额
func (api *QuotaAPI) SetUserQuota(c *gin.Context) {
	username := c.Param("username")

	var req struct {
		Path        string `json:"path" binding:"required"`
		SoftLimit   uint64 `json:"softLimit"`
		HardLimit   uint64 `json:"hardLimit"`
		GracePeriod int    `json:"gracePeriod"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if result := api.DB.Where("username = ?", username).First(&user); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 检查是否已存在配额
	var existingQuota models.UserQuota
	result := api.DB.Where("user_id = ? AND path = ?", user.ID, req.Path).First(&existingQuota)

	if result.Error == nil {
		// 更新现有配额
		existingQuota.SoftLimit = req.SoftLimit
		existingQuota.HardLimit = req.HardLimit
		existingQuota.GracePeriod = req.GracePeriod

		if result := api.DB.Save(&existingQuota); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
	} else {
		// 创建新配额
		quota := models.UserQuota{
			UserID:      user.ID,
			Path:        req.Path,
			SoftLimit:   req.SoftLimit,
			HardLimit:   req.HardLimit,
			GracePeriod: req.GracePeriod,
		}

		if result := api.DB.Create(&quota); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
	}

	// 应用系统配额
	if err := api.setSystemQuota(username, req.Path, req.SoftLimit, req.HardLimit); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to set system quota: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Quota set successfully"})
}

// GetAllQuotas 获取所有用户配额
func (api *QuotaAPI) GetAllQuotas(c *gin.Context) {
	quotaType := c.Query("type") // user, group

	if quotaType == "group" {
		api.GetAllGroupQuotas(c)
		return
	}

	// 默认获取用户配额
	var quotas []models.UserQuota
	if result := api.DB.Preload("User").Find(&quotas); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"quotas": quotas})
}

// getAllGroupQuotas 获取组配额
func (api *QuotaAPI) GetAllGroupQuotas(c *gin.Context) {
	var quotas []models.GroupQuota
	if result := api.DB.Find(&quotas); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"quotas": quotas})
}

// GetQuotaReport 生成配额报告
func (api *QuotaAPI) GetQuotaReport(c *gin.Context) {
	// 获取所有用户配额
	var userQuotas []models.UserQuota
	if result := api.DB.Preload("User").Find(&userQuotas); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// 获取所有组配额
	var groupQuotas []models.GroupQuota
	if result := api.DB.Find(&groupQuotas); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	report := make([]models.QuotaReport, 0)

	// 处理用户配额
	for _, quota := range userQuotas {
		status := "ok"
		usedPercent := 0.0

		if quota.HardLimit > 0 {
			usedPercent = float64(quota.UsedSpace) / float64(quota.HardLimit) * 100
		}

		if quota.SoftLimit > 0 && quota.UsedSpace > quota.SoftLimit {
			status = "warning"
		}

		if quota.HardLimit > 0 && quota.UsedSpace > quota.HardLimit {
			status = "exceeded"
		}

		username := "unknown"
		if quota.User != nil {
			username = quota.User.Username
		}

		reportItem := models.QuotaReport{
			Name:        username,
			Type:        "user",
			Path:        quota.Path,
			UsedSpace:   quota.UsedSpace,
			SoftLimit:   quota.SoftLimit,
			HardLimit:   quota.HardLimit,
			UsedPercent: usedPercent,
			Status:      status,
			FilesUsed:   quota.FilesUsed,
			FilesSoft:   quota.FilesSoft,
			FilesHard:   quota.FilesHard,
			GeneratedAt: time.Now(),
		}

		report = append(report, reportItem)
	}

	// 处理组配额
	for _, quota := range groupQuotas {
		status := "ok"
		usedPercent := 0.0

		if quota.HardLimit > 0 {
			usedPercent = float64(quota.UsedSpace) / float64(quota.HardLimit) * 100
		}

		if quota.SoftLimit > 0 && quota.UsedSpace > quota.SoftLimit {
			status = "warning"
		}

		if quota.HardLimit > 0 && quota.UsedSpace > quota.HardLimit {
			status = "exceeded"
		}

		reportItem := models.QuotaReport{
			Name:        fmt.Sprintf("group_%d", quota.GroupID),
			Type:        "group",
			Path:        quota.Path,
			UsedSpace:   quota.UsedSpace,
			SoftLimit:   quota.SoftLimit,
			HardLimit:   quota.HardLimit,
			UsedPercent: usedPercent,
			Status:      status,
			FilesUsed:   quota.FilesUsed,
			FilesSoft:   quota.FilesSoft,
			FilesHard:   quota.FilesHard,
			GeneratedAt: time.Now(),
		}

		report = append(report, reportItem)
	}

	c.JSON(http.StatusOK, gin.H{"report": report})
}

// setSystemQuota 设置系统配额
func (api *QuotaAPI) setSystemQuota(username, path string, softLimit, hardLimit uint64) error {
	// 使用setquota命令设置配额
	// 格式: setquota -u username soft hard directory

	// 检查文件系统是否支持配额
	cmd := exec.Command("quotaon", "-p", path)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("quota not supported on filesystem: %w", err)
	}

	if !strings.Contains(string(output), "quotaon") {
		// 启用配额
		cmd = exec.Command("quotaon", "-u", path)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to enable quota: %w", err)
		}
	}

	// 设置用户配额
	softMB := softLimit / (1024 * 1024)
	hardMB := hardLimit / (1024 * 1024)

	cmd = exec.Command("setquota", "-u", username,
		strconv.FormatUint(softMB, 10),
		strconv.FormatUint(hardMB, 10),
		"0", "0", path)

	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to set quota: %w, output: %s", err, string(output))
	}

	return nil
}

// getQuotaUsage 获取配额使用情况
func (api *QuotaAPI) getQuotaUsage(username, path string) (usedSpace uint64, err error) {
	// 使用quota命令获取配额使用情况
	cmd := exec.Command("quota", "-u", "-v", username)
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("failed to get quota usage: %w", err)
	}

	// 解析quota输出
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, path) {
			fields := strings.Fields(line)
			if len(fields) >= 6 {
				// quota输出格式: username used soft hard grace
				if used, err := strconv.ParseUint(fields[2], 10, 64); err == nil {
					return used * 1024, nil // 转换为字节
				}
			}
		}
	}

	return 0, fmt.Errorf("quota usage not found")
}

// CheckQuotaAlerts 检查配额告警（应该在定时任务中调用）
func (api *QuotaAPI) CheckQuotaAlerts() {
	// 获取所有用户配额
	var quotas []models.UserQuota
	if result := api.DB.Preload("User").Find(&quotas); result.Error != nil {
		return
	}

	for _, quota := range quotas {
		alertType := ""
		severity := ""
		message := ""

		// 检查软限制
		if quota.SoftLimit > 0 && quota.UsedSpace > quota.SoftLimit {
			alertType = "soft_limit"
			severity = "warning"
			message = fmt.Sprintf("User %s has exceeded soft quota limit on %s", quota.User.Username, quota.Path)
		}

		// 检查硬限制
		if quota.HardLimit > 0 && quota.UsedSpace > quota.HardLimit {
			alertType = "hard_limit"
			severity = "critical"
			message = fmt.Sprintf("User %s has exceeded hard quota limit on %s", quota.User.Username, quota.Path)
		}

		// 检查宽限期
		if quota.GracePeriod > 0 {
			// 这里需要更复杂的逻辑来处理宽限期
		}

		// 如果有告警，记录到数据库
		if alertType != "" {
			alert := models.QuotaAlert{
				UserID:    quota.UserID,
				Type:      "user",
				Path:      quota.Path,
				AlertType: alertType,
				Severity:  severity,
				Message:   message,
				Resolved:  false,
			}

			api.DB.Create(&alert)
		}
	}
}