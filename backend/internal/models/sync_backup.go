package models

import (
	"encoding/json"
	"time"
	"gorm.io/gorm"
)

// SyncJob 数据同步任务模型
type SyncJob struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`

	Name        string `gorm:"not null" json:"name"`
	Description string `json:"description"`
	
	// 同步源和目标
	SourcePath  string `gorm:"not null" json:"sourcePath"`
	DestPath    string `gorm:"not null" json:"destPath"`
	
	// 同步类型: local, remote (rclone), rsync
	Type        string `json:"type"` 
	
	// 计划配置
	Enabled     bool   `json:"enabled"`
	Schedule    string `json:"schedule"` // cron 表达式
	
	// 策略
	DeleteExtra bool   `json:"deleteExtra"` // --delete
	Checksum    bool   `json:"checksum"`    // --checksum
	
	// 状态
	Status      string     `json:"status"` // idle, running, failed, completed
	LastRun     *time.Time `json:"lastRun"`
	LastError   string     `json:"lastError"`
}

// BackupRepo Restic 备份仓库模型
type BackupRepo struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`

	Name string `gorm:"not null" json:"name"`
	// 仓库类型: local, s3, sftp, rest, b2, azure, gs, rclone
	Type string `gorm:"not null" json:"type"`
	// Restic 仓库规格（不含 scheme 部分），如 "/data/restic/repo1" 或 "s3.amazonaws.com/mybucket/path"
	URL string `json:"url"`
	// 完整的 restic -r 参数（由 Type+URL+其他参数拼成，运行时构造）
	// 留空则由后端自动拼接

	// 加密仓库密码（必填，restic init / 每次操作都需要）
	Password string `json:"password,omitempty"`
	// 凭据 / 连接参数（JSON 字符串）：根据后端不同填 access key/secret、SFTP 私钥、B2 keyID 等
	// 这些键会作为环境变量传给 restic（例如 AWS_ACCESS_KEY_ID、RESTIC_REPOSITORY 等）
	EnvJSON string `gorm:"type:text" json:"-"`

	// 起源：
	//   "" / "manual" — 用户在 UI 中手动创建
	//   "external:<container>" — 从外部 Docker 容器自动导入（绑定到该容器）
	// 当 Origin 是 external 时，立即备份 = docker start 对应容器，而不是本地 restic。
	Origin string `gorm:"default:'manual'" json:"origin"`

	// 状态: active / uninitialized / error
	Status      string     `gorm:"default:'uninitialized'" json:"status"`
	LastError   string     `gorm:"type:text" json:"lastError"`
	LastBackup  *time.Time `json:"lastBackup"`
	SnapshotCnt int        `json:"snapshotCount"`
	RepoSize    int64      `json:"repoSize"`
}

// MaskedEnv 返回脱敏后的环境变量（供前端展示）
func (r *BackupRepo) MaskedEnv() map[string]string {
	env := r.ParseEnv()
	out := make(map[string]string, len(env))
	for k, v := range env {
		if len(v) == 0 {
			continue
		}
		out[k] = maskValue(v)
	}
	return out
}

func maskValue(v string) string {
	if len(v) <= 4 {
		return "****"
	}
	return v[:2] + "****" + v[len(v)-2:]
}

// ParseEnv 解析 EnvJSON 为 map
func (r *BackupRepo) ParseEnv() map[string]string {
	out := map[string]string{}
	if r.EnvJSON == "" {
		return out
	}
	_ = json.Unmarshal([]byte(r.EnvJSON), &out)
	return out
}

// BackupTask Restic 备份任务
type BackupTask struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`

	Name string `gorm:"not null" json:"name"`
	// 关联仓库
	RepoID uint       `gorm:"not null" json:"repoId"`
	Repo   BackupRepo `gorm:"foreignKey:RepoID" json:"repo"`

	// 备份源路径（在容器内可访问的路径，例如 /host/data 或 /data/restic-data）
	SourcePath string `gorm:"not null" json:"sourcePath"`
	// 排除规则，每行一条（restic --exclude 语法）
	Excludes string `gorm:"type:text" json:"excludes"`
	// 标签（逗号分隔），用作 restic --tag
	Tags string `json:"tags"`
	// 保留策略 JSON，例如 {"keep_daily":7,"keep_weekly":4,"keep_monthly":12}
	// 留空表示不自动 forget
	RetentionJSON string `gorm:"type:text" json:"retention"`
	// 是否在备份完成后自动 prune
	AutoPrune bool `gorm:"default:false" json:"autoPrune"`

	// 触发模式：
	//   "" / "local"  — 在本容器内直接执行 restic backup
	//   "external"    — 通过 docker start 触发关联的外部容器（要求 Repo.Origin 以 external: 开头）
	TriggerMode string `gorm:"default:'local'" json:"triggerMode"`
	// 当 TriggerMode=external 时绑定的外部容器名；留空则用 Repo.Origin 的容器名
	ExternalContainer string `json:"externalContainer"`

	// 调度
	Enabled  bool   `gorm:"default:false" json:"enabled"`
	Schedule string `json:"schedule"` // cron 表达式（仅展示，目前由用户手动触发）

	// 状态字段
	Status         string     `gorm:"default:'idle'" json:"status"` // idle / running / completed / failed
	LastRun        *time.Time `json:"lastRun"`
	LastError      string     `gorm:"type:text" json:"lastError"`
	LastSnapshotID string     `json:"lastSnapshotId"`
	LastDuration   float64    `json:"lastDuration"` // 秒
}

// ParseRetention 解析保留策略
func (t *BackupTask) ParseRetention() map[string]int {
	out := map[string]int{}
	if t.RetentionJSON == "" {
		return out
	}
	// 用 encoding/json 解析 map[string]json.Number，再转 int
	var raw map[string]float64
	if err := json.Unmarshal([]byte(t.RetentionJSON), &raw); err == nil {
		for k, v := range raw {
			out[k] = int(v)
		}
	}
	return out
}

// BackupSyncJob 仓库到仓库的快照同步任务（restic copy）
type BackupSyncJob struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`

	Name string `gorm:"not null" json:"name"`
	// 源仓库
	SourceRepoID uint       `gorm:"not null" json:"sourceRepoId"`
	SourceRepo   BackupRepo `gorm:"foreignKey:SourceRepoID" json:"sourceRepo"`
	// 目标仓库
	TargetRepoID uint       `gorm:"not null" json:"targetRepoId"`
	TargetRepo   BackupRepo `gorm:"foreignKey:TargetRepoID" json:"targetRepo"`

	Enabled bool `gorm:"default:false" json:"enabled"`

	Status    string     `gorm:"default:'idle'" json:"status"`
	LastRun   *time.Time `json:"lastRun"`
	LastError string     `gorm:"type:text" json:"lastError"`
}
