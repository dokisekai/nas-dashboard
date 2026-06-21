package api

import (
	"nas-dashboard/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// StorageService 全局存储服务实例
var storageService *service.StorageService

// InitStorageService 初始化存储服务
func InitStorageService(db *database.DBConfig) {
	storageService = service.NewStorageService(db.DB)
}

// GetStoragePoolHealth 获取存储池健康状态
func GetStoragePoolHealth(c *gin.Context) {
	poolID := c.Query("poolId")
	if poolID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少存储池ID"})
		return
	}

	health, err := storageService.GetStoragePoolHealth(poolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, health)
}

// GetAllStoragePoolsHealth 获取所有存储池健康状态
func GetAllStoragePoolsHealth(c *gin.Context) {
	var pools []models.StoragePool
	if err := database.DB.Find(&pools).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取存储池失败"})
		return
	}

	var healthReports []service.StoragePoolHealth
	for _, pool := range pools {
		health, err := storageService.GetStoragePoolHealth(pool.ID)
		if err != nil {
			// 跳过错误的存储池
			continue
		}
		healthReports = append(healthReports, *health)
	}

	c.JSON(http.StatusOK, healthReports)
}

// CreateSnapshot 创建快照
func CreateSnapshot(c *gin.Context) {
	var req struct {
		PoolID      string `json:"poolId" binding:"required"`
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	snapshot, err := storageService.CreateSnapshot(req.PoolID, req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, snapshot)
}

// ListSnapshots 列出快照
func ListSnapshots(c *gin.Context) {
	poolID := c.Query("poolId")
	if poolID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少存储池ID"})
		return
	}

	snapshots, err := storageService.ListSnapshots(poolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, snapshots)
}

// DeleteSnapshot 删除快照
func DeleteSnapshot(c *gin.Context) {
	snapshotID := c.Param("id")
	if snapshotID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少快照ID"})
		return
	}

	if err := storageService.DeleteSnapshot(snapshotID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "快照已删除"})
}

// RestoreSnapshot 恢复快照
func RestoreSnapshot(c *gin.Context) {
	snapshotID := c.Param("id")
	if snapshotID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少快照ID"})
		return
	}

	if err := storageService.RestoreSnapshot(snapshotID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "快照已恢复"})
}

// CreateBackupJob 创建备份任务
func CreateBackupJob(c *gin.Context) {
	var job service.BackupJob
	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdJob, err := storageService.CreateBackupJob(job)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdJob)
}

// ListBackupJobs 列出备份任务
func ListBackupJobs(c *gin.Context) {
	jobs, err := storageService.ListBackupJobs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, jobs)
}

// RunBackupJob 执行备份任务
func RunBackupJob(c *gin.Context) {
	jobID := c.Param("id")
	if jobID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少任务ID"})
		return
	}

	// 异步执行备份
	go func() {
		if err := storageService.RunBackupJob(jobID); err != nil {
			log.Printf("备份任务执行失败: %v", err)
		}
	}()

	c.JSON(http.StatusOK, gin.H{"message": "备份任务已启动"})
}

// GetBackupJobStatus 获取备份任务状态
func GetBackupJobStatus(c *gin.Context) {
	jobID := c.Param("id")
	if jobID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少任务ID"})
		return
	}

	var job models.BackupJob
	if err := database.DB.Where("id = ?", jobID).First(&job).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "备份任务不存在"})
		return
	}

	c.JSON(http.StatusOK, job)
}

// DeleteBackupJob 删除备份任务
func DeleteBackupJob(c *gin.Context) {
	jobID := c.Param("id")
	if jobID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少任务ID"})
		return
	}

	if err := database.DB.Delete(&models.BackupJob{}, jobID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除备份任务失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "备份任务已删除"})
}

// CleanOldBackups 清理旧备份
func CleanOldBackups(c *gin.Context) {
	jobID := c.Param("id")
	if jobID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少任务ID"})
		return
	}

	if err := storageService.CleanOldBackups(jobID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "旧备份已清理"})
}

// GetDiskPerformance 获取磁盘性能指标
func GetDiskPerformance(c *gin.Context) {
	device := c.Query("device")
	if device == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少设备参数"})
		return
	}

	// 获取磁盘性能指标
	// 这里需要调用系统命令或读取 /proc/diskstats
	stats := map[string]interface{}{
		"device": device,
		"reads":  0,
		"writes": 0,
		"readBytesPerSec":  0,
		"writeBytesPerSec": 0,
		"readIOPS":  0,
		"writeIOPS": 0,
		"avgLatency": 0,
	}

	c.JSON(http.StatusOK, stats)
}

// GetStorageRecommendations 获取存储优化建议
func GetStorageRecommendations(c *gin.Context) {
	poolID := c.Query("poolId")
	if poolID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少存储池ID"})
		return
	}

	health, err := storageService.GetStoragePoolHealth(poolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	recommendations := []string{}

	// 添加推荐建议
	for _, rec := range health.Recommendations {
		recommendations = append(recommendations, rec)
	}

	// 根据健康状态添加更多建议
	switch health.Status {
	case "healthy":
		recommendations = append(recommendations, "存储池状态良好，建议定期创建快照")
		if health.UsagePercent < 50 {
			recommendations = append(recommendations, "磁盘使用率较低，可以考虑添加更多数据")
		}
	case "degraded":
		recommendations = append(recommendations, "存储池状态异常，请立即检查问题磁盘")
		recommendations = append(recommendations, "建议检查 SMART 数据并考虑更换故障磁盘")
	case "failed":
		recommendations = append(recommendations, "存储池严重故障，数据可能有风险")
		recommendations = append(recommendations, "立即备份重要数据并检查硬件")
	}

	c.JSON(http.StatusOK, gin.H{
		"recommendations": recommendations,
		"health": health.Status,
	})
}