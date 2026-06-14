package api

import (
	"net/http"
	"nas-dashboard/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"os/exec"
	"time"
)

type SyncAPI struct {
	DB *gorm.DB
}

func NewSyncAPI(db *gorm.DB) *SyncAPI {
	return &SyncAPI{DB: db}
}

// GetSyncJobs 获取所有同步任务
func (api *SyncAPI) GetSyncJobs(c *gin.Context) {
	var jobs []models.SyncJob
	if err := api.DB.Find(&jobs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, jobs)
}

// CreateSyncJob 创建同步任务
func (api *SyncAPI) CreateSyncJob(c *gin.Context) {
	var job models.SyncJob
	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	job.Status = "idle"
	if err := api.DB.Create(&job).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, job)
}

// RunSyncJob 手动运行同步任务
func (api *SyncAPI) RunSyncJob(c *gin.Context) {
	id := c.Param("id")
	var job models.SyncJob
	if err := api.DB.First(&job, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	// 异步运行，避免阻塞请求
	go api.executeSync(job)

	c.JSON(http.StatusOK, gin.H{"message": "Sync job started"})
}

func (api *SyncAPI) executeSync(job models.SyncJob) {
	// 更新状态为运行中
	api.DB.Model(&job).Update("status", "running")

	// 根据类型构建命令
	var cmd *exec.Cmd
	switch job.Type {
	case "rsync":
		args := []string{"-avz"}
		if job.DeleteExtra {
			args = append(args, "--delete")
		}
		if job.Checksum {
			args = append(args, "--checksum")
		}
		args = append(args, job.SourcePath, job.DestPath)
		cmd = exec.Command("rsync", args...)
	default:
		// 默认使用 cp 或其它
		cmd = exec.Command("cp", "-r", job.SourcePath, job.DestPath)
	}

	output, err := cmd.CombinedOutput()
	
	now := time.Now()
	status := "completed"
	lastError := ""
	if err != nil {
		status = "failed"
		lastError = string(output)
	}

	api.DB.Model(&job).Updates(map[string]interface{}{
		"status":     status,
		"lastRun":    &now,
		"lastError":  lastError,
	})
}

// Backup API
type BackupAPI struct {
	DB *gorm.DB
}

func NewBackupAPI(db *gorm.DB) *BackupAPI {
	return &BackupAPI{DB: db}
}

func (api *BackupAPI) GetRepos(c *gin.Context) {
	var repos []models.BackupRepo
	api.DB.Find(&repos)
	c.JSON(http.StatusOK, repos)
}

func (api *BackupAPI) CreateRepo(c *gin.Context) {
	var repo models.BackupRepo
	if err := c.ShouldBindJSON(&repo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 初始化 restic 仓库
	// cmd := exec.Command("restic", "-r", repo.URL, "init")
	// env := os.Environ()
	// env = append(env, "RESTIC_PASSWORD="+repo.Password)
	// ... 实际环境中需要更复杂的处理

	api.DB.Create(&repo)
	c.JSON(http.StatusCreated, repo)
}

func (api *BackupAPI) GetTasks(c *gin.Context) {
	var tasks []models.BackupTask
	api.DB.Preload("Repo").Find(&tasks)
	c.JSON(http.StatusOK, tasks)
}

func (api *BackupAPI) CreateTask(c *gin.Context) {
	var task models.BackupTask
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	api.DB.Create(&task)
	c.JSON(http.StatusCreated, task)
}
