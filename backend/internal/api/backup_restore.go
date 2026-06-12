package api

import (
	"compress/gzip"
	"fmt"
	"io"
	"nas-dashboard/internal/database"
	"nas-dashboard/internal/middleware"
	"nas-dashboard/internal/models"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// BackupService 备份服务
type BackupService struct {
	db *gorm.DB
}

// NewBackupService 创建备份服务
func NewBackupService() *BackupService {
	return &BackupService{
		db: database.GetDB(),
	}
}

var backupService = NewBackupService()

// CreateBackupRequest 创建备份请求
type CreateBackupRequest struct {
	Name        string `json:"name" binding:"required"`
	Type        string `json:"type"` // full, incremental, differential
	Description string `json:"description"`
	IncludeDB   bool   `json:"includeDB"`
	IncludeFiles bool   `json:"includeFiles"`
	FilePaths   []string `json:"filePaths"`
}

// RestoreBackupRequest 恢复备份请求
type RestoreBackupRequest struct {
	BackupID uint `json:"backupId" binding:"required"`
}

// GetBackups 获取备份列表
func GetBackups(c *gin.Context) {
	backups, err := backupService.GetBackups()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"backups": backups,
		"total":   len(backups),
	})
}

// GetBackup 获取单个备份信息
func GetBackup(c *gin.Context) {
	id := c.Param("id")

	backupID, err := parseBackupID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid backup ID"})
		return
	}

	backup, err := backupService.GetBackup(backupID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Backup not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, backup)
}

// CreateBackup 创建备份
func CreateBackup(c *gin.Context) {
	var req CreateBackupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 获取用户信息
	userID, _ := middleware.GetUserIDAsUint(c)
	username, _ := middleware.GetUsername(c)

	// 创建备份
	backup, err := backupService.CreateBackup(req, userID, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Backup created successfully",
		"backup":  backup,
	})
}

// DeleteBackup 删除备份
func DeleteBackup(c *gin.Context) {
	id := c.Param("id")

	backupID, err := parseBackupID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid backup ID"})
		return
	}

	// 删除备份
	if err := backupService.DeleteBackup(backupID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Backup deleted successfully",
		"id":      backupID,
	})
}

// RestoreBackup 恢复备份
func RestoreBackup(c *gin.Context) {
	var req RestoreBackupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 获取用户信息
	userID, _ := middleware.GetUserIDAsUint(c)
	username, _ := middleware.GetUsername(c)

	// 恢复备份
	if err := backupService.RestoreBackup(req.BackupID, userID, username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Backup restored successfully",
		"backupId": req.BackupID,
	})
}

// DownloadBackup 下载备份文件
func DownloadBackup(c *gin.Context) {
	id := c.Param("id")

	backupID, err := parseBackupID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid backup ID"})
		return
	}

	// 获取备份信息
	backup, err := backupService.GetBackup(backupID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Backup not found"})
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(backup.FilePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Backup file not found"})
		return
	}

	// 设置响应头
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(backup.FilePath)))
	c.Header("Content-Type", "application/gzip")

	// 发送文件
	http.ServeFile(c.Writer, c.Request, backup.FilePath)
}

// GetBackups 获取备份列表
func (s *BackupService) GetBackups() ([]models.BackupRecord, error) {
	var backups []models.BackupRecord
	if err := s.db.Order("created_at DESC").Find(&backups).Error; err != nil {
		return nil, err
	}
	return backups, nil
}

// GetBackup 获取单个备份
func (s *BackupService) GetBackup(id uint) (*models.BackupRecord, error) {
	var backup models.BackupRecord
	if err := s.db.First(&backup, id).Error; err != nil {
		return nil, err
	}
	return &backup, nil
}

// CreateBackup 创建备份
func (s *BackupService) CreateBackup(req CreateBackupRequest, userID uint, username string) (*models.BackupRecord, error) {
	// 创建备份目录
	backupDir := "/var/backups/nas-dashboard"
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create backup directory: %w", err)
	}

	// 生成备份文件名
	timestamp := time.Now().Format("20060102_150405")
	backupFileName := fmt.Sprintf("%s_%s.tar.gz", req.Name, timestamp)
	backupFilePath := filepath.Join(backupDir, backupFileName)

	// 创建备份记录
	backup := &models.BackupRecord{
		Name:        req.Name,
		Type:        req.Type,
		FilePath:    backupFilePath,
		Status:      "in_progress",
		CreatedBy:   username,
		Description: req.Description,
	}

	if err := s.db.Create(backup).Error; err != nil {
		return nil, fmt.Errorf("failed to create backup record: %w", err)
	}

	// 创建备份文件
	if err := s.createBackupFile(backup, req); err != nil {
		backup.Status = "failed"
		s.db.Save(backup)
		return nil, fmt.Errorf("failed to create backup: %w", err)
	}

	// 获取文件大小
	if info, err := os.Stat(backupFilePath); err == nil {
		backup.Size = info.Size()
	}

	// 更新备份记录
	now := time.Now()
	backup.Status = "completed"
	backup.CompletedAt = &now
	if err := s.db.Save(backup).Error; err != nil {
		return nil, fmt.Errorf("failed to update backup record: %w", err)
	}

	return backup, nil
}

// createBackupFile 创建备份文件
func (s *BackupService) createBackupFile(backup *models.BackupRecord, req CreateBackupRequest) error {
	// 创建临时目录
	tempDir := filepath.Join(os.TempDir(), fmt.Sprintf("backup_%d", backup.ID))
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	// 备份数据库
	if req.IncludeDB {
		if err := s.backupDatabase(tempDir); err != nil {
			return fmt.Errorf("failed to backup database: %w", err)
		}
	}

	// 备份文件
	if req.IncludeFiles && len(req.FilePaths) > 0 {
		if err := s.backupFiles(tempDir, req.FilePaths); err != nil {
			return fmt.Errorf("failed to backup files: %w", err)
		}
	}

	// 创建压缩文件
	if err := s.createCompressedBackup(tempDir, backup.FilePath); err != nil {
		return fmt.Errorf("failed to create compressed backup: %w", err)
	}

	return nil
}

// backupDatabase 备份数据库
func (s *BackupService) backupDatabase(destDir string) error {
	// 使用 pg_dump 备份 PostgreSQL 数据库
	dbBackupFile := filepath.Join(destDir, "database.sql.gz")

	// 获取数据库连接信息
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")

	// 执行 pg_dump
	cmd := exec.Command("pg_dump", "-h", dbHost, "-U", dbUser, dbName)

	// 创建 gzip 压缩文件
	outputFile, err := os.Create(dbBackupFile)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	gzipWriter := gzip.NewWriter(outputFile)
	defer gzipWriter.Close()

	// 设置命令输出
	cmd.Stdout = gzipWriter

	// 设置环境变量
	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", os.Getenv("DB_PASSWORD")))

	// 执行命令
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("pg_dump failed: %w", err)
	}

	return nil
}

// backupFiles 备份文件
func (s *BackupService) backupFiles(destDir string, filePaths []string) error {
	// 创建文件列表
	for _, filePath := range filePaths {
		if err := s.copyPathToBackup(filePath, destDir); err != nil {
			return fmt.Errorf("failed to backup %s: %w", filePath, err)
		}
	}
	return nil
}

// copyPathToBackup 复制路径到备份目录
func (s *BackupService) copyPathToBackup(source, destDir string) error {
	// 获取路径信息
	info, err := os.Stat(source)
	if err != nil {
		return err
	}

	// 创建目标路径
	destPath := filepath.Join(destDir, source)
	if info.IsDir() {
		if err := os.MkdirAll(destPath, 0755); err != nil {
			return err
		}
		// 递归复制目录
		return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			relPath, err := filepath.Rel(source, path)
			if err != nil {
				return err
			}

			destPath := filepath.Join(destDir, source, relPath)

			if info.IsDir() {
				return os.MkdirAll(destPath, info.Mode())
			}

			return copyFile(path, destPath)
		})
	}

	// 复制单个文件
	return copyFile(source, destPath)
}

// copyFile 复制文件
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return err
	}

	// 复制文件权限
	info, err := os.Stat(src)
	if err != nil {
		return err
	}

	return os.Chmod(dst, info.Mode())
}

// createCompressedBackup 创建压缩备份
func (s *BackupService) createCompressedBackup(sourceDir, destFile string) error {
	cmd := exec.Command("tar", "-czf", destFile, "-C", sourceDir, ".")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("tar command failed: %w, output: %s", err, string(output))
	}
	return nil
}

// DeleteBackup 删除备份
func (s *BackupService) DeleteBackup(id uint) error {
	// 获取备份信息
	backup, err := s.GetBackup(id)
	if err != nil {
		return err
	}

	// 删除备份文件
	if backup.FilePath != "" {
		if err := os.Remove(backup.FilePath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to delete backup file: %w", err)
		}
	}

	// 删除数据库记录
	return s.db.Delete(backup).Error
}

// RestoreBackup 恢复备份
func (s *BackupService) RestoreBackup(backupID uint, userID uint, username string) error {
	// 获取备份信息
	backup, err := s.GetBackup(backupID)
	if err != nil {
		return err
	}

	// 检查备份文件是否存在
	if _, err := os.Stat(backup.FilePath); os.IsNotExist(err) {
		return fmt.Errorf("backup file not found: %s", backup.FilePath)
	}

	// 创建临时目录
	tempDir := filepath.Join(os.TempDir(), fmt.Sprintf("restore_%d", backupID))
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	// 解压备份文件
	if err := s.extractBackup(backup.FilePath, tempDir); err != nil {
		return fmt.Errorf("failed to extract backup: %w", err)
	}

	// 恢复数据库
	dbBackupFile := filepath.Join(tempDir, "database.sql.gz")
	if _, err := os.Stat(dbBackupFile); err == nil {
		if err := s.restoreDatabase(dbBackupFile); err != nil {
			return fmt.Errorf("failed to restore database: %w", err)
		}
	}

	// 恢复文件
	if err := s.restoreFiles(tempDir); err != nil {
		return fmt.Errorf("failed to restore files: %w", err)
	}

	return nil
}

// extractBackup 解压备份文件
func (s *BackupService) extractBackup(backupFile, destDir string) error {
	cmd := exec.Command("tar", "-xzf", backupFile, "-C", destDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("tar command failed: %w, output: %s", err, string(output))
	}
	return nil
}

// restoreDatabase 恢复数据库
func (s *BackupService) restoreDatabase(backupFile string) error {
	// 打开压缩文件
	file, err := os.Open(backupFile)
	if err != nil {
		return err
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	// 获取数据库连接信息
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")

	// 执行 psql 恢复
	cmd := exec.Command("psql", "-h", dbHost, "-U", dbUser, dbName)
	cmd.Stdin = gzipReader

	// 设置环境变量
	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", os.Getenv("DB_PASSWORD")))

	// 执行命令
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("psql failed: %w", err)
	}

	return nil
}

// restoreFiles 恢复文件
func (s *BackupService) restoreFiles(sourceDir string) error {
	// 遍历备份目录
	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过数据库备份文件
		if strings.HasSuffix(path, "database.sql.gz") {
			return nil
		}

		// 计算相对路径
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		// 跳过根目录
		if relPath == "." {
			return nil
		}

		// 目标路径
		destPath := filepath.Join("/", relPath)

		if info.IsDir() {
			// 创建目录
			if err := os.MkdirAll(destPath, info.Mode()); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", destPath, err)
			}
		} else {
			// 复制文件
			if err := copyFile(path, destPath); err != nil {
				return fmt.Errorf("failed to copy file to %s: %w", destPath, err)
			}
		}

		return nil
	})
}

// parseBackupID 解析备份ID
func parseBackupID(id string) (uint, error) {
	var backupID uint
	if _, err := fmt.Sscanf(id, "%d", &backupID); err != nil {
		return 0, err
	}
	return backupID, nil
}
