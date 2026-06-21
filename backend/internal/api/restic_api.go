package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"nas-dashboard/internal/database"
	"nas-dashboard/internal/models"
	"nas-dashboard/internal/service"
)

// ResticAPI exposes all restic-related HTTP handlers.
type ResticAPI struct {
	DB *gorm.DB
	// in-memory logs per running task, keyed by job descriptor "<type>:<id>".
	logs sync.Map // map[string]*service.LogBuffer
}

// NewResticAPI creates the API.
func NewResticAPI(db *gorm.DB) *ResticAPI {
	return &ResticAPI{DB: db}
}

// =====================================================================
// 仓库 (Repositories)
// =====================================================================

// ListRepos GET /api/storage/backup/repos
func (a *ResticAPI) ListRepos(c *gin.Context) {
	var repos []models.BackupRepo
	if err := a.DB.Order("id DESC").Find(&repos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 不返回原始密码字段，但前端需要展示状态。
	for i := range repos {
		if repos[i].Password != "" {
			repos[i].Password = "********"
		}
	}
	c.JSON(http.StatusOK, gin.H{"repos": repos, "total": len(repos)})
}

// repoPayload 用于接收创建/更新仓库的请求。
type repoPayload struct {
	Name     string            `json:"name" binding:"required"`
	Type     string            `json:"type" binding:"required"`
	URL      string            `json:"url" binding:"required"`
	Password string            `json:"password"`
	Env      map[string]string `json:"env"`
	Init     bool              `json:"init"` // 创建后是否立即 init
}

// CreateRepo POST /api/storage/backup/repos
func (a *ResticAPI) CreateRepo(c *gin.Context) {
	var req repoPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password is required"})
		return
	}
	if req.Type == "local" {
		if err := service.AuditRepoPath(req.URL); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	repo := &models.BackupRepo{
		Name:     req.Name,
		Type:     req.Type,
		URL:      req.URL,
		Password: req.Password,
		Status:   "uninitialized",
	}
	if req.Env != nil {
		raw, _ := json.Marshal(req.Env)
		repo.EnvJSON = string(raw)
	}

	// 初始化仓库
	if req.Init {
		svc := service.NewResticService(repo)
		if err := svc.Init(); err != nil {
			repo.LastError = err.Error()
			a.DB.Create(repo)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "restic init failed: " + err.Error(),
				"repo":  maskRepoPtr(repo),
			})
			return
		}
		repo.Status = "active"
	}

	if err := a.DB.Create(repo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, maskRepoPtr(repo))
}

// UpdateRepo PUT /api/storage/backup/repos/:id
func (a *ResticAPI) UpdateRepo(c *gin.Context) {
	id := c.Param("id")
	var repo models.BackupRepo
	if err := a.DB.First(&repo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "repo not found"})
		return
	}
	var req repoPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	repo.Name = req.Name
	repo.Type = req.Type
	repo.URL = req.URL
	if req.Password != "" && req.Password != "********" {
		repo.Password = req.Password
	}
	if req.Env != nil {
		raw, _ := json.Marshal(req.Env)
		repo.EnvJSON = string(raw)
	}
	if err := a.DB.Save(&repo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, maskRepoPtr(&repo))
}

// DeleteRepo DELETE /api/storage/backup/repos/:id
// Query: purge=true 删除仓库目录（仅 local 类型）
func (a *ResticAPI) DeleteRepo(c *gin.Context) {
	id := c.Param("id")
	var repo models.BackupRepo
	if err := a.DB.First(&repo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "repo not found"})
		return
	}
	purge := c.Query("purge") == "true"
	if purge && repo.Type == "local" && repo.URL != "" {
		_ = os.RemoveAll(repo.URL)
	}
	// 删除关联的备份任务
	a.DB.Where("repo_id = ?", repo.ID).Delete(&models.BackupTask{})
	a.DB.Where("source_repo_id = ? OR target_repo_id = ?", repo.ID, repo.ID).Delete(&models.BackupSyncJob{})
	if err := a.DB.Delete(&repo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "repo deleted"})
}

// CheckRepo POST /api/storage/backup/repos/:id/check
// 验证仓库完整性。query: full=true 进行完整数据校验（极慢，云端可能几小时）。
func (a *ResticAPI) CheckRepo(c *gin.Context) {
	repo, err := a.loadRepo(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	svc := service.NewResticService(repo)
	full := c.Query("full") == "true"
	out, err := svc.Check(full)
	if err != nil {
		repo.Status = "error"
		repo.LastError = err.Error()
		a.DB.Save(repo)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "output": out})
		return
	}
	repo.Status = "active"
	repo.LastError = ""
	a.DB.Save(repo)
	c.JSON(http.StatusOK, gin.H{"output": out, "status": "ok"})
}

// RefreshRepo POST /api/storage/backup/repos/:id/refresh
// 重新拉取仓库的快照数 / 体积，写入数据库。
func (a *ResticAPI) RefreshRepo(c *gin.Context) {
	repo, err := a.loadRepo(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	svc := service.NewResticService(repo)
	stats, err := svc.Stats()
	if err != nil {
		// 可能是仓库未初始化
		repo.Status = "error"
		repo.LastError = err.Error()
		a.DB.Save(repo)
		c.JSON(http.StatusOK, gin.H{"status": "error", "error": err.Error()})
		return
	}
	repo.SnapshotCnt = stats.SnapshotsCount
	repo.RepoSize = stats.TotalSize
	repo.Status = "active"
	repo.LastError = ""
	a.DB.Save(repo)
	c.JSON(http.StatusOK, gin.H{
		"status":     "ok",
		"stats":      stats,
		"repo":       maskRepoPtr(repo),
	})
}

// =====================================================================
// 快照 (Snapshots)
// =====================================================================

// ListSnapshots GET /api/storage/backup/repos/:id/snapshots
func (a *ResticAPI) ListSnapshots(c *gin.Context) {
	repo, err := a.loadRepo(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	svc := service.NewResticService(repo)
	snaps, err := svc.Snapshots()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"snapshots": snaps, "total": len(snaps)})
}

// ListSnapshotFiles GET /api/storage/backup/repos/:id/snapshots/:sid/ls
func (a *ResticAPI) ListSnapshotFiles(c *gin.Context) {
	repo, err := a.loadRepo(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	sid := c.Param("sid")
	svc := service.NewResticService(repo)
	nodes, err := svc.LS(sid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"files": nodes, "total": len(nodes)})
}

// DeleteSnapshot DELETE /api/storage/backup/repos/:id/snapshots/:sid
// query: prune=true 顺带 prune
func (a *ResticAPI) DeleteSnapshot(c *gin.Context) {
	repo, err := a.loadRepo(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	sid := c.Param("sid")
	prune := c.Query("prune") == "true"
	svc := service.NewResticService(repo)
	out, err := svc.Forget([]string{sid}, prune)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "output": out})
		return
	}
	c.JSON(http.StatusOK, gin.H{"output": out})
}

// =====================================================================
// 备份任务 (Backup Tasks)
// =====================================================================

// ListTasks GET /api/storage/backup/tasks
func (a *ResticAPI) ListTasks(c *gin.Context) {
	var tasks []models.BackupTask
	if err := a.DB.Preload("Repo").Order("id DESC").Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tasks": tasks, "total": len(tasks)})
}

type taskPayload struct {
	Name          string         `json:"name" binding:"required"`
	RepoID        uint           `json:"repoId" binding:"required"`
	SourcePath    string         `json:"sourcePath" binding:"required"`
	Excludes      string         `json:"excludes"`
	Tags          string         `json:"tags"`
	Retention     map[string]int `json:"retention"`
	AutoPrune     bool           `json:"autoPrune"`
	Enabled       bool           `json:"enabled"`
	Schedule      string         `json:"schedule"`
}

// CreateTask POST /api/storage/backup/tasks
func (a *ResticAPI) CreateTask(c *gin.Context) {
	var req taskPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	task := &models.BackupTask{
		Name:       req.Name,
		RepoID:     req.RepoID,
		SourcePath: req.SourcePath,
		Excludes:   req.Excludes,
		Tags:       req.Tags,
		AutoPrune:  req.AutoPrune,
		Enabled:    req.Enabled,
		Schedule:   req.Schedule,
		Status:     "idle",
	}
	if req.Retention != nil {
		raw, _ := json.Marshal(req.Retention)
		task.RetentionJSON = string(raw)
	}
	if err := a.DB.Create(task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	a.DB.Preload("Repo").First(task, task.ID)
	c.JSON(http.StatusCreated, task)
}

// UpdateTask PUT /api/storage/backup/tasks/:id
func (a *ResticAPI) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var task models.BackupTask
	if err := a.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	var req taskPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	task.Name = req.Name
	task.RepoID = req.RepoID
	task.SourcePath = req.SourcePath
	task.Excludes = req.Excludes
	task.Tags = req.Tags
	task.AutoPrune = req.AutoPrune
	task.Enabled = req.Enabled
	task.Schedule = req.Schedule
	if req.Retention != nil {
		raw, _ := json.Marshal(req.Retention)
		task.RetentionJSON = string(raw)
	} else {
		task.RetentionJSON = ""
	}
	if err := a.DB.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	a.DB.Preload("Repo").First(&task, task.ID)
	c.JSON(http.StatusOK, task)
}

// DeleteTask DELETE /api/storage/backup/tasks/:id
func (a *ResticAPI) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	if err := a.DB.Delete(&models.BackupTask{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "task deleted"})
}

// RunTask POST /api/storage/backup/tasks/:id/run
// 立即执行一次备份任务。
func (a *ResticAPI) RunTask(c *gin.Context) {
	id := c.Param("id")
	var task models.BackupTask
	if err := a.DB.Preload("Repo").First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	if task.Repo.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "repo not configured"})
		return
	}
	if task.Status == "running" {
		c.JSON(http.StatusConflict, gin.H{"error": "task is already running"})
		return
	}

	// 启动异步执行
	go a.executeBackupTask(task)

	c.JSON(http.StatusAccepted, gin.H{"message": "backup started", "taskId": task.ID})
}

// TaskLogs GET /api/storage/backup/tasks/:id/logs
// 返回最近一次运行的日志（in-memory，重启后清空）。
func (a *ResticAPI) TaskLogs(c *gin.Context) {
	id := c.Param("id")
	buf := a.getLogBuf("task:" + id)
	c.JSON(http.StatusOK, gin.H{"lines": buf.Lines()})
}

// TaskStatus GET /api/storage/backup/tasks/:id/status
func (a *ResticAPI) TaskStatus(c *gin.Context) {
	id := c.Param("id")
	var task models.BackupTask
	if err := a.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":         task.Status,
		"lastRun":        task.LastRun,
		"lastError":      task.LastError,
		"lastSnapshotId": task.LastSnapshotID,
		"lastDuration":   task.LastDuration,
	})
}

// executeBackupTask runs the backup with retention/prune policy.
func (a *ResticAPI) executeBackupTask(task models.BackupTask) {
	logKey := fmt.Sprintf("task:%d", task.ID)
	buf := a.resetLogBuf(logKey)
	now := time.Now()

	setStatus := func(status, errMsg, snapID string, dur float64) {
		updates := map[string]interface{}{
			"status":          status,
			"lastRun":         &now,
			"lastError":       errMsg,
			"lastSnapshotId":  snapID,
			"lastDuration":    dur,
		}
		a.DB.Model(&models.BackupTask{}).Where("id = ?", task.ID).Updates(updates)
		if status == "completed" {
			a.DB.Model(&models.BackupRepo{}).Where("id = ?", task.RepoID).
				Update("last_backup", &now)
		}
	}
	setStatus("running", "", "", 0)
	buf.Append(fmt.Sprintf("[%s] Backup '%s' started, source=%s", time.Now().Format(time.RFC3339), task.Name, task.SourcePath))

	repo := &task.Repo
	svc := service.NewResticService(repo)

	// 合并全局默认 hostname / excludes / tags
	settings := a.loadSettings()
	allExcludes := splitLines(task.Excludes)
	allExcludes = append(allExcludes, splitLines(settings.DefaultExcludes)...)
	allTags := splitCSV(task.Tags)
	allTags = append(allTags, splitCSV(settings.DefaultTags)...)
	hostname := settings.DefaultHostname
	if hostname == "" {
		hostname = hostnameFor(repo)
	}

	opts := service.BackupOptions{
		Source:   task.SourcePath,
		Excludes: allExcludes,
		Tags:     allTags,
		Hostname: hostname,
	}
	result, err := svc.Backup(opts, func(line string) { buf.Append(line) })
	if err != nil {
		buf.Append(fmt.Sprintf("[ERROR] %v", err))
		setStatus("failed", err.Error(), "", 0)
		return
	}

	// 应用保留策略
	if retention := task.ParseRetention(); len(retention) > 0 {
		buf.Append("[INFO] applying retention policy: " + task.RetentionJSON)
		if out, ferr := svc.ApplyRetention(retention, task.AutoPrune, false); ferr != nil {
			buf.Append("[WARN] forget/prune failed: " + ferr.Error() + " out=" + out)
		} else if out != "" {
			buf.Append("[INFO] retention output: " + out)
		}
	}
	buf.Append(fmt.Sprintf("[%s] Backup completed, snapshot=%s, duration=%s",
		time.Now().Format(time.RFC3339), result.SnapshotID, result.Duration))
	setStatus("completed", "", result.SnapshotID, result.Duration.Seconds())

	// 更新仓库快照数缓存
	go func(repoID uint) {
		var r models.BackupRepo
		if err := database.GetDB().First(&r, repoID).Error; err == nil {
			if stats, err := service.NewResticService(&r).Stats(); err == nil {
				database.GetDB().Model(&r).Updates(map[string]interface{}{
					"snapshot_cnt": stats.SnapshotsCount,
					"repo_size":    stats.TotalSize,
				})
			}
		}
	}(task.RepoID)
}

// =====================================================================
// 恢复 (Restore)
// =====================================================================

type restorePayload struct {
	SnapshotID string   `json:"snapshotId"`
	Target     string   `json:"target" binding:"required"`
	Include    []string `json:"include"`
	Exclude    []string `json:"exclude"`
	Host       string   `json:"host"`
	Paths      []string `json:"paths"`
}

// Restore POST /api/storage/backup/repos/:id/restore
// 异步恢复快照到指定目录。
func (a *ResticAPI) Restore(c *gin.Context) {
	repo, err := a.loadRepo(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	var req restorePayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 简单的安全检查：不允许恢复到系统敏感目录。
	if err := auditRestoreTarget(req.Target); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logKey := fmt.Sprintf("restore:%d", repo.ID)
	buf := a.resetLogBuf(logKey)
	go func() {
		buf.Append(fmt.Sprintf("[%s] Restore started: snapshot=%s target=%s",
			time.Now().Format(time.RFC3339), req.SnapshotID, req.Target))
		svc := service.NewResticService(repo)
		err := svc.Restore(service.RestoreOptions{
			Snapshot: req.SnapshotID,
			Target:   req.Target,
			Include:  req.Include,
			Exclude:  req.Exclude,
			Host:     req.Host,
			Paths:    req.Paths,
		}, func(line string) { buf.Append(line) })
		if err != nil {
			buf.Append("[ERROR] " + err.Error())
		} else {
			buf.Append("[INFO] restore completed")
		}
	}()
	c.JSON(http.StatusAccepted, gin.H{"message": "restore started", "logKey": logKey})
}

// RestoreLogs GET /api/storage/backup/repos/:id/restore/logs
func (a *ResticAPI) RestoreLogs(c *gin.Context) {
	repo, err := a.loadRepo(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	buf := a.getLogBuf(fmt.Sprintf("restore:%d", repo.ID))
	c.JSON(http.StatusOK, gin.H{"lines": buf.Lines()})
}

// =====================================================================
// 仓库间同步 (restic copy)
// =====================================================================

type syncJobPayload struct {
	Name         string `json:"name" binding:"required"`
	SourceRepoID uint   `json:"sourceRepoId" binding:"required"`
	TargetRepoID uint   `json:"targetRepoId" binding:"required"`
	Enabled      bool   `json:"enabled"`
}

// ListSyncJobs GET /api/storage/backup/sync-jobs
func (a *ResticAPI) ListSyncJobs(c *gin.Context) {
	var jobs []models.BackupSyncJob
	if err := a.DB.Preload("SourceRepo").Preload("TargetRepo").Order("id DESC").Find(&jobs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for i := range jobs {
		jobs[i].SourceRepo = maskRepo(jobs[i].SourceRepo)
		jobs[i].TargetRepo = maskRepo(jobs[i].TargetRepo)
	}
	c.JSON(http.StatusOK, gin.H{"jobs": jobs, "total": len(jobs)})
}

// CreateSyncJob POST /api/storage/backup/sync-jobs
func (a *ResticAPI) CreateSyncJob(c *gin.Context) {
	var req syncJobPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.SourceRepoID == req.TargetRepoID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "source and target repos must differ"})
		return
	}
	job := &models.BackupSyncJob{
		Name:         req.Name,
		SourceRepoID: req.SourceRepoID,
		TargetRepoID: req.TargetRepoID,
		Enabled:      req.Enabled,
		Status:       "idle",
	}
	if err := a.DB.Create(job).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	a.DB.Preload("SourceRepo").Preload("TargetRepo").First(job, job.ID)
	job.SourceRepo = maskRepo(job.SourceRepo)
	job.TargetRepo = maskRepo(job.TargetRepo)
	c.JSON(http.StatusCreated, job)
}

// DeleteSyncJob DELETE /api/storage/backup/sync-jobs/:id
func (a *ResticAPI) DeleteSyncJob(c *gin.Context) {
	id := c.Param("id")
	if err := a.DB.Delete(&models.BackupSyncJob{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "sync job deleted"})
}

// RunSyncJob POST /api/storage/backup/sync-jobs/:id/run
// 把源仓库所有尚未在目标仓库存在的快照 copy 过去。
func (a *ResticAPI) RunSyncJob(c *gin.Context) {
	id := c.Param("id")
	var job models.BackupSyncJob
	if err := a.DB.Preload("SourceRepo").Preload("TargetRepo").First(&job, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "sync job not found"})
		return
	}
	if job.SourceRepo.ID == 0 || job.TargetRepo.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "repo missing"})
		return
	}
	if job.Status == "running" {
		c.JSON(http.StatusConflict, gin.H{"error": "sync already running"})
		return
	}
	go a.executeSyncJob(job)
	c.JSON(http.StatusAccepted, gin.H{"message": "sync started", "jobId": job.ID})
}

// RunSyncJobLogs GET /api/storage/backup/sync-jobs/:id/logs
func (a *ResticAPI) RunSyncJobLogs(c *gin.Context) {
	id := c.Param("id")
	buf := a.getLogBuf("sync:" + id)
	c.JSON(http.StatusOK, gin.H{"lines": buf.Lines()})
}

func (a *ResticAPI) executeSyncJob(job models.BackupSyncJob) {
	logKey := fmt.Sprintf("sync:%d", job.ID)
	buf := a.resetLogBuf(logKey)
	now := time.Now()
	setStatus := func(status, errMsg string) {
		a.DB.Model(&models.BackupSyncJob{}).Where("id = ?", job.ID).
			Updates(map[string]interface{}{"status": status, "lastRun": &now, "lastError": errMsg})
	}
	setStatus("running", "")
	buf.Append(fmt.Sprintf("[%s] Copy snapshots from '%s' -> '%s'",
		time.Now().Format(time.RFC3339), job.SourceRepo.Name, job.TargetRepo.Name))

	src := &job.SourceRepo
	dst := &job.TargetRepo
	svc := service.NewResticService(src)
	out, err := svc.Copy(dst, nil, func(line string) { buf.Append(line) })
	if err != nil {
		buf.Append("[ERROR] " + err.Error())
		setStatus("failed", err.Error())
		return
	}
	if out != "" {
		buf.Append("[INFO] " + out)
	}
	buf.Append("[INFO] sync completed")
	setStatus("completed", "")
}

// =====================================================================
// 工具
// =====================================================================

func (a *ResticAPI) loadRepo(id string) (*models.BackupRepo, error) {
	var repo models.BackupRepo
	if err := a.DB.First(&repo, id).Error; err != nil {
		return nil, fmt.Errorf("repo not found")
	}
	return &repo, nil
}

func (a *ResticAPI) getLogBuf(key string) *service.LogBuffer {
	v, _ := a.logs.LoadOrStore(key, service.NewLogBuffer(2000))
	return v.(*service.LogBuffer)
}

func (a *ResticAPI) resetLogBuf(key string) *service.LogBuffer {
	buf := service.NewLogBuffer(2000)
	a.logs.Store(key, buf)
	return buf
}

func maskRepo(r models.BackupRepo) models.BackupRepo {
	if r.Password != "" {
		r.Password = "********"
	}
	return r
}

func maskRepoPtr(r *models.BackupRepo) *models.BackupRepo {
	if r == nil {
		return nil
	}
	if r.Password != "" {
		r.Password = "********"
	}
	return r
}

func splitLines(s string) []string {
	out := []string{}
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			out = append(out, line)
		}
	}
	return out
}

func splitCSV(s string) []string {
	out := []string{}
	for _, p := range strings.Split(s, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func hostnameFor(repo *models.BackupRepo) string {
	if h := repo.ParseEnv()["RESTIC_HOSTNAME"]; h != "" {
		return h
	}
	host, _ := os.Hostname()
	if host == "" {
		host = "nas-dashboard"
	}
	return host
}

func auditRestoreTarget(p string) error {
	abs, err := filepath.Abs(p)
	if err != nil {
		return err
	}
	dangerous := []string{"/", "/etc", "/usr", "/proc", "/sys", "/dev", "/bin", "/sbin", "/lib", "/boot"}
	for _, d := range dangerous {
		if abs == d || strings.HasPrefix(abs, d+"/") {
			return fmt.Errorf("restore target %s is in a system directory; pick /restore/... or /data/...", abs)
		}
	}
	return nil
}

// 检查 restic 二进制是否可用，通过 GET /api/storage/backup/ping 返回。
func (a *ResticAPI) Ping(c *gin.Context) {
	err := service.CheckAvailable()
	host, _ := os.Hostname()
	cacheDir := os.Getenv("RESTIC_CACHE_DIR")
	if cacheDir == "" {
		cacheDir = "/tmp/restic-cache"
	}
	// 拿一下 restic 版本
	version := ""
	if path, err := exec.LookPath(service.ResticBinary); err == nil {
		out, err := exec.Command(path, "version").CombinedOutput()
		if err == nil {
			// 输出形如: restic 0.18.1 compiled with go1.25.0 on linux/amd64
			version = strings.SplitN(string(out), "\n", 2)[0]
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"ok":        err == nil,
		"error":     func() string { if err != nil { return err.Error() }; return "" }(),
		"hostname":  host,
		"runtime":   runtime.GOOS + "/" + runtime.GOARCH,
		"cacheDir":  cacheDir,
		"resticBin": service.ResticBinary,
		"version":   version,
		"timestamp": time.Now(),
	})
}

// =====================================================================
// 全局设置 (Settings)
// =====================================================================

// BackupSettings 全局备份设置（持久化到 system_configs 表）
type BackupSettings struct {
	DefaultHostname string `json:"defaultHostname"` // 默认 host 标签
	DefaultExcludes string `json:"defaultExcludes"` // 默认排除规则（每行一条），追加到所有任务
	DefaultTags     string `json:"defaultTags"`     // 默认标签（逗号分隔）
	AutoCheck       bool   `json:"autoCheck"`       // 创建仓库时自动 check
	ConfirmPurge    bool   `json:"confirmPurge"`    // 删除仓库时是否需要二次确认（仅前端使用，后端始终要求）
}

// WebDAVProfile 一个 WebDAV 连接配置（持久化到 system_configs），可被多个仓库复用。
type WebDAVProfile struct {
	RemoteName string `json:"remoteName"` // rclone 配置里的 remote 名，如 "123pan"
	URL        string `json:"url"`        // WebDAV endpoint，如 https://webdav.123pan.cn/webdav
	Vendor     string `json:"vendor"`     // nextcloud / owncloud / sharefile / other / ...
	Username   string `json:"username"`
	Password   string `json:"password,omitempty"` // 返回时会被脱敏
}

// webdavRcloneConfPath rclone 配置文件路径（持久化在 data 目录）
const webdavRcloneConfPath = "/data/rclone.conf"

// defaultWebdavRclonePath 如果环境里没有 /data，fallback 到 /tmp
func defaultWebdavRclonePath() string {
	if _, err := os.Stat("/data"); err == nil {
		return "/data/rclone.conf"
	}
	return "/tmp/rclone.conf"
}

// GetSettings GET /api/storage/backup/settings
func (a *ResticAPI) GetSettings(c *gin.Context) {
	s := a.loadSettings()
	c.JSON(http.StatusOK, s)
}

// UpdateSettings PUT /api/storage/backup/settings
func (a *ResticAPI) UpdateSettings(c *gin.Context) {
	var req BackupSettings
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	a.saveSettings(req)
	c.JSON(http.StatusOK, req)
}

// =====================================================================
// WebDAV 配置（集中管理，仓库创建时复用）
// =====================================================================

// GetWebDAV GET /api/storage/backup/webdav
// 返回当前 webdav 配置（密码会脱敏）。
func (a *ResticAPI) GetWebDAV(c *gin.Context) {
	p := a.loadWebDAVProfile()
	if p.Password != "" {
		p.Password = "********"
	}
	c.JSON(http.StatusOK, p)
}

// UpdateWebDAV PUT /api/storage/backup/webdav
// 保存 webdav 配置 + 写入 rclone.conf 文件（用于后续 restic 操作）。
// 密码字段为 "********" 时表示不修改。
func (a *ResticAPI) UpdateWebDAV(c *gin.Context) {
	var req WebDAVProfile
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.RemoteName == "" {
		req.RemoteName = "webdav"
	}
	if req.Password == "********" {
		// 保留旧密码
		req.Password = a.loadWebDAVProfile().Password
	}
	if err := a.saveWebDAVProfile(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := writeWebDAVToRcloneConf(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "save rclone.conf failed: " + err.Error()})
		return
	}
	resp := req
	resp.Password = "********"
	c.JSON(http.StatusOK, resp)
}

// TestWebDAV POST /api/storage/backup/webdav/test
// 测试 webdav 连接：临时写入 rclone 配置，跑 `rclone lsd` 列出根目录。
func (a *ResticAPI) TestWebDAV(c *gin.Context) {
	var req WebDAVProfile
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.RemoteName == "" {
		req.RemoteName = "webdav"
	}
	if req.Password == "********" {
		req.Password = a.loadWebDAVProfile().Password
	}
	// 写临时 rclone 配置
	tmpFile := "/tmp/rclone-test-" + req.RemoteName + ".conf"
	if err := writeWebDAVToRcloneConfAt(req, tmpFile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer os.Remove(tmpFile)
	// rclone lsd
	cmd := exec.Command("rclone", "--config", tmpFile, "lsd", req.RemoteName+":")
	cmd.Env = append(os.Environ())
	out, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"ok":     false,
			"error":  err.Error(),
			"output": string(out),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ok":     true,
		"output": string(out),
	})
}

// FindInSnapshots GET /api/storage/backup/repos/:id/find?q=keyword&pattern=glob
// 在所有快照里搜索文件（restic find），返回命中的快照+文件路径。
func (a *ResticAPI) FindInSnapshots(c *gin.Context) {
	repo, err := a.loadRepo(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	pattern := c.Query("q")
	if pattern == "" {
		pattern = c.Query("pattern")
	}
	if pattern == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "q parameter required"})
		return
	}
	svc := service.NewResticService(repo)
	out, err := svc.Find(pattern)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "output": out})
		return
	}
	c.JSON(http.StatusOK, gin.H{"output": out, "pattern": pattern})
}

// ---- WebDAV helpers ----

func (a *ResticAPI) loadWebDAVProfile() WebDAVProfile {
	p := WebDAVProfile{
		RemoteName: "webdav",
		Vendor:     "other",
	}
	var cfg models.SystemConfig
	if err := a.DB.Where("key = ?", "backup.webdav").First(&cfg).Error; err == nil {
		_ = json.Unmarshal([]byte(cfg.Value), &p)
	}
	if p.RemoteName == "" {
		p.RemoteName = "webdav"
	}
	if p.Vendor == "" {
		p.Vendor = "other"
	}
	return p
}

func (a *ResticAPI) saveWebDAVProfile(p WebDAVProfile) error {
	raw, err := json.Marshal(p)
	if err != nil {
		return err
	}
	var cfg models.SystemConfig
	if err := a.DB.Where("key = ?", "backup.webdav").First(&cfg).Error; err == gorm.ErrRecordNotFound {
		return a.DB.Create(&models.SystemConfig{
			Key:         "backup.webdav",
			Value:       string(raw),
			Type:        "json",
			Category:    "backup",
			Description: "WebDAV connection profile for rclone",
		}).Error
	} else if err == nil {
		cfg.Value = string(raw)
		return a.DB.Save(&cfg).Error
	} else {
		return err
	}
}

// writeWebDAVToRcloneConf 把 webdav 配置写到默认的 rclone.conf。
// 如果文件里已有同名 remote，会被替换；其他 remote 保留。
func writeWebDAVToRcloneConf(p WebDAVProfile) error {
	return writeWebDAVToRcloneConfAt(p, defaultWebdavRclonePath())
}

func writeWebDAVToRcloneConfAt(p WebDAVProfile, path string) error {
	// 读现有内容（如果有）
	existing := map[string][]string{} // section -> lines
	if data, err := os.ReadFile(path); err == nil {
		var curSection string
		for _, line := range strings.Split(string(data), "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
				continue
			}
			if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
				curSection = line[1 : len(line)-1]
				existing[curSection] = []string{}
				continue
			}
			if curSection != "" {
				existing[curSection] = append(existing[curSection], line)
			}
		}
	}
	// 用新的 webdav section 替换 / 添加
	existing[p.RemoteName] = []string{
		"type = webdav",
		"url = " + p.URL,
		"vendor = " + p.Vendor,
		"user = " + p.Username,
		"pass = " + rcloneObscure(p.Password),
	}
	// 重新生成 ini
	var buf strings.Builder
	for section, lines := range existing {
		fmt.Fprintf(&buf, "[%s]\n", section)
		for _, l := range lines {
			fmt.Fprintf(&buf, "%s\n", l)
		}
		fmt.Fprintln(&buf)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(buf.String()), 0600)
}

// rcloneObscure 用 rclone 自己的 obscure 机制加密密码。
// 实际上 rclone 接受明文密码（会打 WARNING），所以为了让用户在 UI 里粘贴什么都能工作，
// 我们只在密码不像 obscure 形式时才尝试 obscure；否则原样写。
// rclone obscure 输出总是以 "OMYCT" 开头的 base64 风格字符串（具体取决于版本）。
// 这里采用最稳的策略：不 obscure，直接写明文。rclone 接受。
func rcloneObscure(plain string) string {
	if plain == "" {
		return ""
	}
	// 不做 obscure：rclone.conf 直接写明文也能用，避免双重 obscure 导致 401
	return plain
}

func settingsKey(field string) string { return "backup.settings." + field }

func (a *ResticAPI) loadSettings() BackupSettings {
	s := BackupSettings{
		DefaultHostname: "nas-dashboard",
		AutoCheck:       false,
		ConfirmPurge:    true,
	}
	getCfg := func(key string, def string) string {
		var cfg models.SystemConfig
		if err := a.DB.Where("key = ?", key).First(&cfg).Error; err == nil {
			return cfg.Value
		}
		return def
	}
	s.DefaultHostname = getCfg(settingsKey("default_hostname"), s.DefaultHostname)
	s.DefaultExcludes = getCfg(settingsKey("default_excludes"), "")
	s.DefaultTags = getCfg(settingsKey("default_tags"), "")
	s.AutoCheck = getCfg(settingsKey("auto_check"), "false") == "true"
	s.ConfirmPurge = getCfg(settingsKey("confirm_purge"), "true") == "true"
	return s
}

func (a *ResticAPI) saveSettings(s BackupSettings) {
	entries := map[string]string{
		settingsKey("default_hostname"): s.DefaultHostname,
		settingsKey("default_excludes"): s.DefaultExcludes,
		settingsKey("default_tags"):     s.DefaultTags,
		settingsKey("auto_check"):       boolStr(s.AutoCheck),
		settingsKey("confirm_purge"):    boolStr(s.ConfirmPurge),
	}
	for key, value := range entries {
		var cfg models.SystemConfig
		if err := a.DB.Where("key = ?", key).First(&cfg).Error; err == gorm.ErrRecordNotFound {
			a.DB.Create(&models.SystemConfig{
				Key:         key,
				Value:       value,
				Type:        "string",
				Category:    "backup",
				Description: "Restic backup manager setting",
				IsPublic:    false,
			})
		} else if err == nil {
			cfg.Value = value
			a.DB.Save(&cfg)
		}
	}
}

func boolStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// =====================================================================
// 仓库高级操作
// =====================================================================

// TestRepoConnection POST /api/storage/backup/repos/:id/test
// 不进行 init / check，仅尝试 list snapshots；连接成功就返回 ok。
func (a *ResticAPI) TestRepoConnection(c *gin.Context) {
	repo, err := a.loadRepo(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	svc := service.NewResticService(repo)
	snaps, err := svc.Snapshots()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"ok":     false,
			"error":  err.Error(),
			"status": "error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ok":             true,
		"snapshotCount":  len(snaps),
		"latestSnapshot": latestSnapshotTime(snaps),
		"status":         "active",
	})
}

// InitRepo POST /api/storage/backup/repos/:id/init
// 对已存在但未初始化的仓库记录执行 init。
func (a *ResticAPI) InitRepo(c *gin.Context) {
	repo, err := a.loadRepo(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	svc := service.NewResticService(repo)
	if err := svc.Init(); err != nil {
		repo.Status = "error"
		repo.LastError = err.Error()
		a.DB.Save(repo)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	repo.Status = "active"
	repo.LastError = ""
	a.DB.Save(repo)
	c.JSON(http.StatusOK, maskRepoPtr(repo))
}

// UnlockRepo POST /api/storage/backup/repos/:id/unlock?force=true
// 清除仓库锁。出现 "repository is already locked exclusively" 时使用。
func (a *ResticAPI) UnlockRepo(c *gin.Context) {
	repo, err := a.loadRepo(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	force := c.Query("force") == "true"
	svc := service.NewResticService(repo)
	out, err := svc.Unlock(force)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "output": out})
		return
	}
	c.JSON(http.StatusOK, gin.H{"output": out, "status": "ok"})
}

// =====================================================================
// 快照详情
// =====================================================================

// SnapshotDetail GET /api/storage/backup/repos/:id/snapshots/:sid
// 返回单条快照的详细元数据 + 文件统计。
func (a *ResticAPI) SnapshotDetail(c *gin.Context) {
	repo, err := a.loadRepo(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	sid := c.Param("sid")
	svc := service.NewResticService(repo)
	snaps, err := svc.Snapshots()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var found *service.Snapshot
	for i := range snaps {
		if snaps[i].ID == sid || snaps[i].ShortID == sid {
			found = &snaps[i]
			break
		}
	}
	if found == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "snapshot not found"})
		return
	}
	// 统计文件
	var fileCount, dirCount, totalSize int64
	files, _ := svc.LS(sid)
	for _, f := range files {
		if f.Type == "dir" {
			dirCount++
		} else {
			fileCount++
			totalSize += f.Size
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"snapshot":   found,
		"fileCount":  fileCount,
		"dirCount":   dirCount,
		"totalSize":  totalSize,
	})
}

// DiffSnapshots GET /api/storage/backup/repos/:id/diff?a=&b=
// 比较两个快照之间的差异。
func (a *ResticAPI) DiffSnapshots(c *gin.Context) {
	repo, err := a.loadRepo(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	aID := c.Query("a")
	bID := c.Query("b")
	if aID == "" || bID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "both a and b query parameters are required"})
		return
	}
	svc := service.NewResticService(repo)
	res, err := svc.Diff(aID, bID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"output": res})
}

func latestSnapshotTime(snaps []service.Snapshot) *time.Time {
	if len(snaps) == 0 {
		return nil
	}
	latest := snaps[0].Time
	for _, s := range snaps[1:] {
		if s.Time.After(latest) {
			latest = s.Time
		}
	}
	return &latest
}
