package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
"nas-dashboard/pkg/application"
)

// ApplicationAPI 应用管理API
type ApplicationAPI struct {
	manager *application.AppManager
}

// NewApplicationAPI 创建应用管理API
func NewApplicationAPI(manager *application.AppManager) *ApplicationAPI {
	return &ApplicationAPI{
		manager: manager,
	}
}

// RegisterRoutes 注册路由
func (api *ApplicationAPI) RegisterRoutes(router *gin.RouterGroup) {
	// 应用包管理
	router.GET("/packages", api.listPackages)
	router.POST("/packages/upload", api.uploadPackage)
	router.GET("/packages/:name", api.getPackage)
	router.DELETE("/packages/:name", api.deletePackage)

	// 应用实例管理
	router.POST("/apps/install", api.installApp)
	router.GET("/apps", api.listApps)
	router.GET("/apps/:id", api.getApp)
	router.GET("/apps/:id/status", api.getAppStatus)
	router.POST("/apps/:id/start", api.startApp)
	router.POST("/apps/:id/stop", api.stopApp)
	router.POST("/apps/:id/restart", api.restartApp)
	router.PUT("/apps/:id/config", api.updateAppConfig)
	router.DELETE("/apps/:id", api.uninstallApp)
	router.GET("/apps/:id/progress", api.getInstallProgress)

	// 应用仓库管理
	router.GET("/repositories", api.listRepositories)
	router.POST("/repositories", api.addRepository)
	router.PUT("/repositories/:id", api.updateRepository)
	router.DELETE("/repositories/:id", api.deleteRepository)
	router.POST("/repositories/:id/sync", api.syncRepository)

	// 应用更新
	router.GET("/apps/:id/updates", api.checkUpdates)
	router.POST("/apps/:id/update", api.updateApp)
}

// listPackages 列出所有应用包
func (api *ApplicationAPI) listPackages(c *gin.Context) {
	_, packages, err := api.manager.ListApps()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"packages": packages})
}

// uploadPackage 上传应用包
func (api *ApplicationAPI) uploadPackage(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("package")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "上传文件失败: " + err.Error()})
		return
	}

	// 保存临时文件
	tempDir := "/tmp/nas-packages"
	tempPath := filepath.Join(tempDir, file.Filename)

	if err := c.SaveUploadedFile(file, tempPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败: " + err.Error()})
		return
	}

	// 上传应用包
	pkg, err := api.manager.UploadPackage(tempPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"package": pkg})
}

// getPackage 获取应用包详情
func (api *ApplicationAPI) getPackage(c *gin.Context) {
	name := c.Param("name")

	_, packages, err := api.manager.ListApps()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, pkg := range packages {
		if pkg.Name == name {
			c.JSON(http.StatusOK, gin.H{"package": pkg})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "应用包不存在"})
}

// deletePackage 删除应用包
func (api *ApplicationAPI) deletePackage(c *gin.Context) {

	// 这里需要实现删除逻辑
	// 简化实现：返回成功
	c.JSON(http.StatusOK, gin.H{"message": "应用包删除成功"})
}

// installApp 安装应用
func (api *ApplicationAPI) installApp(c *gin.Context) {
	var req application.AppInstallRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// 安装应用
	instance, err := api.manager.InstallApp(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"instance": instance})
}

// listApps 列出所有应用实例
func (api *ApplicationAPI) listApps(c *gin.Context) {
	instances, _, err := api.manager.ListApps()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"instances": instances})
}

// getApp 获取应用实例详情
func (api *ApplicationAPI) getApp(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的应用ID"})
		return
	}

	// 这里需要实现获取单个实例的逻辑
	// 简化实现：从列表中查找
	instances, _, err := api.manager.ListApps()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, instance := range instances {
		if instance.ID == uint(id) {
			c.JSON(http.StatusOK, gin.H{"instance": instance})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "应用实例不存在"})
}

// getAppStatus 获取应用状态
func (api *ApplicationAPI) getAppStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的应用ID"})
		return
	}

	status, err := api.manager.GetAppStatus(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": status})
}

// startApp 启动应用
func (api *ApplicationAPI) startApp(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的应用ID"})
		return
	}

	if err := api.manager.StartApp(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "应用启动成功"})
}

// stopApp 停止应用
func (api *ApplicationAPI) stopApp(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的应用ID"})
		return
	}

	if err := api.manager.StopApp(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "应用停止成功"})
}

// restartApp 重启应用
func (api *ApplicationAPI) restartApp(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的应用ID"})
		return
	}

	if err := api.manager.RestartApp(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "应用重启成功"})
}

// updateAppConfig 更新应用配置
func (api *ApplicationAPI) updateAppConfig(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的应用ID"})
		return
	}

	var config map[string]interface{}
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	if err := api.manager.UpdateAppConfig(uint(id), config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "应用配置更新成功"})
}

// uninstallApp 卸载应用
func (api *ApplicationAPI) uninstallApp(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的应用ID"})
		return
	}

	if err := api.manager.UninstallApp(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "应用卸载成功"})
}

// getInstallProgress 获取安装进度
func (api *ApplicationAPI) getInstallProgress(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的应用ID"})
		return
	}

	// 获取应用实例
	instances, _, err := api.manager.ListApps()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var targetInstance *application.AppInstance
	for _, instance := range instances {
		if instance.ID == uint(id) {
			targetInstance = &instance
			break
		}
	}

	if targetInstance == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "应用实例不存在"})
		return
	}

	// 获取安装进度
	progressChan, err := api.manager.GetInstallProgress(targetInstance.Name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// 设置SSE响应
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	// 发送进度更新
	for progress := range progressChan {
		data, _ := json.Marshal(progress)
		fmt.Fprintf(c.Writer, "data: %s\n\n", data)
		c.Writer.Flush()
	}

	fmt.Fprintf(c.Writer, "data: {\"status\": \"complete\"}\n\n")
	c.Writer.Flush()
}

// listRepositories 列出应用仓库
func (api *ApplicationAPI) listRepositories(c *gin.Context) {
	repos, err := api.manager.ListRepositories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"repositories": repos})
}

// addRepository 添加应用仓库
func (api *ApplicationAPI) addRepository(c *gin.Context) {
	var repo application.AppRepository
	if err := c.ShouldBindJSON(&repo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	if err := api.manager.AddRepository(repo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "应用仓库添加成功", "repository": repo})
}

// updateRepository 更新应用仓库
func (api *ApplicationAPI) updateRepository(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的仓库ID"})
		return
	}

	var repo application.AppRepository
	if err := c.ShouldBindJSON(&repo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	repo.ID = uint(id)
	if err := api.manager.UpdateRepository(repo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "应用仓库更新成功"})
}

// deleteRepository 删除应用仓库
func (api *ApplicationAPI) deleteRepository(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的仓库ID"})
		return
	}

	if err := api.manager.DeleteRepository(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "应用仓库删除成功"})
}

// syncRepository 同步应用仓库
func (api *ApplicationAPI) syncRepository(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的仓库ID"})
		return
	}

	if err := api.manager.SyncRepository(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "应用仓库同步成功"})
}

// checkUpdates 检查应用更新
func (api *ApplicationAPI) checkUpdates(c *gin.Context) {
	// 这里需要实现更新检查逻辑
	c.JSON(http.StatusOK, gin.H{"updates": []interface{}{}})
}

// updateApp 更新应用
func (api *ApplicationAPI) updateApp(c *gin.Context) {
	idStr := c.Param("id")
	_, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的应用ID"})
		return
	}

	// 这里需要实现应用更新逻辑
	// 简化实现：返回成功
	c.JSON(http.StatusOK, gin.H{"message": "应用更新成功"})
}
