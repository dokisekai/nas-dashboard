package api

import (
	"fmt"
	"io"
	"mime"
	"nas-dashboard/internal/database"
	"nas-dashboard/internal/middleware"
	"nas-dashboard/internal/models"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// FileService 文件服务
type FileService struct {
	db *gorm.DB
}

// NewFileService 创建文件服务
func NewFileService() *FileService {
	return &FileService{
		db: database.GetDB(),
	}
}

var fileService = NewFileService()

// FileInfo 文件信息
type FileInfo struct {
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	Size        int64     `json:"size"`
	IsDir       bool      `json:"isDir"`
	ModTime     time.Time `json:"modTime"`
	Permissions string    `json:"permissions"`
	Owner       string    `json:"owner"`
	Group       string    `json:"group"`
	MimeType    string    `json:"mimeType,omitempty"`
}

// ListFilesRequest 列出文件请求
type ListFilesRequest struct {
	Path string `json:"path" binding:"required"`
}

// UploadRequest 上传文件请求
type UploadRequest struct {
	Path string `form:"path" binding:"required"`
}

// CreateDirectoryRequest 创建目录请求
type CreateDirectoryRequest struct {
	Path        string `json:"path" binding:"required"`
	Permissions string `json:"permissions"`
}

// MoveRequest 移动/重命名请求
type MoveRequest struct {
	OldPath string `json:"oldPath" binding:"required"`
	NewPath string `json:"newPath" binding:"required"`
}

// DeleteRequest 删除请求
type DeleteRequest struct {
	Path  string `json:"path" binding:"required"`
	Force bool   `json:"force"` // 强制删除，不回收站
}

// ListFiles 列出目录中的文件
func ListFiles(c *gin.Context) {
	var req ListFilesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 验证路径
	cleanPath, err := fileService.ValidatePath(req.Path)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查权限
	if err := fileService.CheckPermission(c, cleanPath, "read"); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	// 获取用户ID
	userID, _ := middleware.GetUserIDAsUint(c)

	// 读取目录
	files, err := fileService.ListDirectory(cleanPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 记录访问日志
	fileService.LogAccess(userID, cleanPath, "read", true, "")

	c.JSON(http.StatusOK, gin.H{
		"path":  cleanPath,
		"files": files,
	})
}

// GetFileInfo 获取文件信息
func GetFileInfo(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Path parameter is required"})
		return
	}

	// 验证路径
	cleanPath, err := fileService.ValidatePath(path)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查权限
	if err := fileService.CheckPermission(c, cleanPath, "read"); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	// 获取用户ID
	userID, _ := middleware.GetUserIDAsUint(c)

	// 获取文件信息
	fileInfo, err := fileService.GetFileInfo(cleanPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// 记录访问日志
	fileService.LogAccess(userID, cleanPath, "read", true, "")

	c.JSON(http.StatusOK, fileInfo)
}

// DownloadFile 下载文件
func DownloadFile(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Path parameter is required"})
		return
	}

	// 验证路径
	cleanPath, err := fileService.ValidatePath(path)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查权限
	if err := fileService.CheckPermission(c, cleanPath, "read"); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	// 获取用户ID
	userID, _ := middleware.GetUserIDAsUint(c)

	// 打开文件
	file, err := os.Open(cleanPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	defer file.Close()

	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file info"})
		return
	}

	// 设置响应头
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileInfo.Name()))
	c.Header("Content-Type", getMimeType(cleanPath))
	c.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	// 发送文件
	http.ServeContent(c.Writer, c.Request, fileInfo.Name(), fileInfo.ModTime(), file)

	// 记录访问日志
	fileService.LogAccess(userID, cleanPath, "download", true, "")
}

// UploadFile 上传文件
func UploadFile(c *gin.Context) {
	var req UploadRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 验证路径
	cleanPath, err := fileService.ValidatePath(req.Path)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查权限
	if err := fileService.CheckPermission(c, cleanPath, "write"); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// 获取用户ID
	userID, _ := middleware.GetUserIDAsUint(c)

	// 创建目标文件
	destPath := filepath.Join(cleanPath, file.Filename)
	destFile, err := os.Create(destPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file"})
		return
	}
	defer destFile.Close()

	// 打开上传的文件
	srcFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open uploaded file"})
		return
	}
	defer srcFile.Close()

	// 复制文件内容
	if _, err := io.Copy(destFile, srcFile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// 记录访问日志
	fileService.LogAccess(userID, destPath, "upload", true, "")

	c.JSON(http.StatusOK, gin.H{
		"message":  "File uploaded successfully",
		"filename": file.Filename,
		"path":     destPath,
	})
}

// CreateDirectory 创建目录
func CreateDirectory(c *gin.Context) {
	var req CreateDirectoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 验证路径
	cleanPath, err := fileService.ValidatePath(req.Path)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查权限
	if err := fileService.CheckPermission(c, cleanPath, "write"); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	// 获取用户ID
	userID, _ := middleware.GetUserIDAsUint(c)

	// 创建目录
	permissions := os.FileMode(0755)
	if req.Permissions != "" {
		_, err := fmt.Sscanf(req.Permissions, "%o", &permissions)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid permissions format"})
			return
		}
	}

	if err := os.MkdirAll(cleanPath, permissions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory"})
		return
	}

	// 记录访问日志
	fileService.LogAccess(userID, cleanPath, "create", true, "")

	c.JSON(http.StatusOK, gin.H{
		"message": "Directory created successfully",
		"path":    cleanPath,
	})
}

// MoveFile 移动/重命名文件
func MoveFile(c *gin.Context) {
	var req MoveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 验证路径
	oldPath, err := fileService.ValidatePath(req.OldPath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newPath, err := fileService.ValidatePath(req.NewPath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查权限
	if err := fileService.CheckPermission(c, oldPath, "write"); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	if err := fileService.CheckPermission(c, newPath, "write"); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	// 获取用户ID
	userID, _ := middleware.GetUserIDAsUint(c)

	// 移动文件
	if err := os.Rename(oldPath, newPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to move file"})
		return
	}

	// 记录访问日志
	fileService.LogAccess(userID, oldPath, "move", true, "to: "+newPath)

	c.JSON(http.StatusOK, gin.H{
		"message": "File moved successfully",
		"oldPath": oldPath,
		"newPath": newPath,
	})
}

// DeleteFile 删除文件
func DeleteFile(c *gin.Context) {
	var req DeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 验证路径
	cleanPath, err := fileService.ValidatePath(req.Path)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查权限
	if err := fileService.CheckPermission(c, cleanPath, "delete"); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	// 获取用户ID
	userID, _ := middleware.GetUserIDAsUint(c)

	// 删除文件
	if req.Force {
		// 强制删除
		err = os.RemoveAll(cleanPath)
	} else {
		// 软删除（移动到回收站）
		err = fileService.MoveToTrash(cleanPath)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
		return
	}

	// 记录访问日志
	fileService.LogAccess(userID, cleanPath, "delete", true, "")

	c.JSON(http.StatusOK, gin.H{
		"message": "File deleted successfully",
		"path":    cleanPath,
	})
}

// ValidatePath 验证并清理路径
func (s *FileService) ValidatePath(path string) (string, error) {
	// 清理路径
	cleanPath := filepath.Clean(path)

	// 检查路径是否为绝对路径
	if !filepath.IsAbs(cleanPath) {
		return "", fmt.Errorf("path must be absolute")
	}

	// 检查路径是否包含危险字符
	if strings.Contains(cleanPath, "..") {
		return "", fmt.Errorf("path cannot contain '..'")
	}

	// 检查路径是否存在
	if _, err := os.Stat(cleanPath); os.IsNotExist(err) {
		return "", fmt.Errorf("path does not exist")
	}

	return cleanPath, nil
}

// CheckPermission 检查文件访问权限
func (s *FileService) CheckPermission(c *gin.Context, path, operation string) error {
	// 获取当前用户信息
	username, exists := c.Get("username")
	if !exists {
		return fmt.Errorf("user not authenticated")
	}

	// 管理员拥有所有权限
	if role, exists := c.Get("role"); exists && role == "admin" {
		return nil
	}

	// 检查用户是否是文件所有者
	// 这里可以添加更复杂的权限检查逻辑
	// 目前简化为：用户可以访问自己的主目录

	// 检查路径是否在用户的主目录下
	userHomeDir := fmt.Sprintf("/home/%s", username)
	if !strings.HasPrefix(path, userHomeDir) {
		return fmt.Errorf("access denied: path outside user home directory")
	}

	return nil
}

// ListDirectory 列出目录内容
func (s *FileService) ListDirectory(path string) ([]FileInfo, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var files []FileInfo
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		fullPath := filepath.Join(path, entry.Name())
		permissions := getFilePermissions(info.Mode())
		owner := getFileOwner(info)

		fileInfo := FileInfo{
			Name:        entry.Name(),
			Path:        fullPath,
			Size:        info.Size(),
			IsDir:       entry.IsDir(),
			ModTime:     info.ModTime(),
			Permissions: permissions,
			Owner:       owner,
		}

		if !entry.IsDir() {
			fileInfo.MimeType = getMimeType(fullPath)
		}

		files = append(files, fileInfo)
	}

	return files, nil
}

// GetFileInfo 获取文件详细信息
func (s *FileService) GetFileInfo(path string) (*FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	return &FileInfo{
		Name:        info.Name(),
		Path:        path,
		Size:        info.Size(),
		IsDir:       info.IsDir(),
		ModTime:     info.ModTime(),
		Permissions: getFilePermissions(info.Mode()),
		MimeType:    getMimeType(path),
	}, nil
}

// MoveToTrash 移动文件到回收站
func (s *FileService) MoveToTrash(path string) error {
	// 创建回收站目录
	trashDir := "/home/.Trash"
	if err := os.MkdirAll(trashDir, 0700); err != nil {
		return err
	}

	// 生成唯一的文件名
	timestamp := time.Now().Format("20060102150405")
	baseName := filepath.Base(path)
	trashName := fmt.Sprintf("%s_%s", timestamp, baseName)
	trashPath := filepath.Join(trashDir, trashName)

	return os.Rename(path, trashPath)
}

// LogAccess 记录文件访问日志
func (s *FileService) LogAccess(userID uint, path, operation string, success bool, errorMsg string) {
	// 如果数据库不可用，跳过日志记录
	if s.db == nil {
		return
	}

	access := &models.FileSystemAccess{
		UserID:    userID,
		Path:      path,
		Operation: operation,
		Success:   success,
		ErrorMsg:  errorMsg,
	}

	if err := s.db.Create(access).Error; err != nil {
		// 记录失败不影响主要操作
		fmt.Printf("Failed to log file access: %v\n", err)
	}
}

// getFilePermissions 获取文件权限字符串
func getFilePermissions(mode os.FileMode) string {
	permissions := ""

	if mode&0400 != 0 {
		permissions += "r"
	} else {
		permissions += "-"
	}

	if mode&0200 != 0 {
		permissions += "w"
	} else {
		permissions += "-"
	}

	if mode&0100 != 0 {
		permissions += "x"
	} else {
		permissions += "-"
	}

	return permissions
}

// getFileOwner 获取文件所有者（简化版）
func getFileOwner(info os.FileInfo) string {
	// 这里应该使用系统调用获取文件所有者
	// 简化实现，返回当前用户
	return "unknown"
}

// getMimeType 获取文件MIME类型
func getMimeType(path string) string {
	ext := filepath.Ext(path)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		return "application/octet-stream"
	}
	return mimeType
}
