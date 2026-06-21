package api

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gorm.io/gorm"

	"nas-dashboard/internal/models"
)

// seedBackupConfig 在 backend 启动时读取 config/backup-seed.json，
// 把其中声明的仓库 / 任务 / 设置同步到数据库（已存在的不覆盖）。
// 这样新机器部署、git clone 后只要 docker compose up 就能直接看到预设的备份配置。
func seedBackupConfig(db *gorm.DB) {
	if db == nil {
		return
	}
	paths := []string{
		"config/backup-seed.json",
		"/data/backup-seed.json",
		"/app/config/backup-seed.json",
	}
	var seedPath string
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			seedPath = p
			break
		}
	}
	if seedPath == "" {
		log.Printf("[seed] no backup-seed.json found, skipping")
		return
	}
	raw, err := os.ReadFile(seedPath)
	if err != nil {
		log.Printf("[seed] read %s failed: %v", seedPath, err)
		return
	}
	// 用 json.RawMessage 解码，跳过 _comment / _doc 等元字段
	var doc struct {
		Repos    []json.RawMessage `json:"repos"`
		Tasks    []json.RawMessage `json:"tasks"`
		Settings json.RawMessage   `json:"settings"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		log.Printf("[seed] parse %s failed: %v", seedPath, err)
		return
	}

	repoNameToID := map[string]uint{}

	// 1. 仓库
	for _, raw := range doc.Repos {
		var spec struct {
			Name     string            `json:"name"`
			Type     string            `json:"type"`
			URL      string            `json:"url"`
			Password string            `json:"password"`
			Env      map[string]string `json:"env"`
			Init     bool              `json:"init"`
		}
		if err := json.Unmarshal(raw, &spec); err != nil {
			continue
		}
		if spec.Name == "" || spec.Type == "" || spec.URL == "" {
			continue
		}
		// 占位符替换
		spec.Password = expandEnv(spec.Password)
		for k, v := range spec.Env {
			spec.Env[k] = expandEnv(v)
		}
		// 查询是否已存在
		var existing models.BackupRepo
		err := db.Where("name = ?", spec.Name).First(&existing).Error
		if err == nil {
			// 已存在，跳过（不覆盖用户修改）
			repoNameToID[spec.Name] = existing.ID
			continue
		}
		if err != gorm.ErrRecordNotFound {
			log.Printf("[seed] query repo %s: %v", spec.Name, err)
			continue
		}
		envJSON, _ := json.Marshal(spec.Env)
		repo := &models.BackupRepo{
			Name:     spec.Name,
			Type:     spec.Type,
			URL:      spec.URL,
			Password: spec.Password,
			EnvJSON:  string(envJSON),
			Status:   "uninitialized",
			Origin:   "seed",
		}
		if err := db.Create(repo).Error; err != nil {
			log.Printf("[seed] create repo %s: %v", spec.Name, err)
			continue
		}
		repoNameToID[spec.Name] = repo.ID
		log.Printf("[seed] created repo %s (id=%d)", repo.Name, repo.ID)
	}

	// 2. 任务
	for _, raw := range doc.Tasks {
		var spec struct {
			Name      string         `json:"name"`
			RepoName  string         `json:"repoName"`
			SourcePath string        `json:"sourcePath"`
			Excludes  string         `json:"excludes"`
			Tags      string         `json:"tags"`
			Retention map[string]int `json:"retention"`
			AutoPrune bool           `json:"autoPrune"`
		}
		if err := json.Unmarshal(raw, &spec); err != nil {
			continue
		}
		if spec.Name == "" || spec.RepoName == "" || spec.SourcePath == "" {
			continue
		}
		repoID, ok := repoNameToID[spec.RepoName]
		if !ok {
			// 没在 seed 里声明、也可能已经在 DB 里
			var r models.BackupRepo
			if err := db.Where("name = ?", spec.RepoName).First(&r).Error; err == nil {
				repoID = r.ID
			} else {
				log.Printf("[seed] task %s: repo %s not found, skipping", spec.Name, spec.RepoName)
				continue
			}
		}
		var existing models.BackupTask
		if err := db.Where("name = ?", spec.Name).First(&existing).Error; err == nil {
			continue
		}
		retentionJSON := ""
		if len(spec.Retention) > 0 {
			b, _ := json.Marshal(spec.Retention)
			retentionJSON = string(b)
		}
		task := &models.BackupTask{
			Name:          spec.Name,
			RepoID:        repoID,
			SourcePath:    spec.SourcePath,
			Excludes:      spec.Excludes,
			Tags:          spec.Tags,
			RetentionJSON: retentionJSON,
			AutoPrune:     spec.AutoPrune,
			Enabled:       false,
			Status:        "idle",
		}
		if err := db.Create(task).Error; err != nil {
			log.Printf("[seed] create task %s: %v", spec.Name, err)
			continue
		}
		log.Printf("[seed] created task %s (repo=%s)", task.Name, spec.RepoName)
	}

	// 3. 设置
	if len(doc.Settings) > 0 {
		var s struct {
			DefaultHostname string `json:"defaultHostname"`
			DefaultExcludes string `json:"defaultExcludes"`
			DefaultTags     string `json:"defaultTags"`
			AutoCheck       bool   `json:"autoCheck"`
			ConfirmPurge    bool   `json:"confirmPurge"`
		}
		if err := json.Unmarshal(doc.Settings, &s); err == nil {
			saveCfg := func(key, val string) {
				var cfg models.SystemConfig
				if err := db.Where("key = ?", key).First(&cfg).Error; err == gorm.ErrRecordNotFound {
					db.Create(&models.SystemConfig{
						Key: key, Value: val, Type: "string",
						Category: "backup", IsPublic: false,
					})
				}
			}
			saveCfg("backup.settings.default_hostname", s.DefaultHostname)
			saveCfg("backup.settings.default_excludes", s.DefaultExcludes)
			saveCfg("backup.settings.default_tags", s.DefaultTags)
			saveCfg("backup.settings.auto_check", boolStr(s.AutoCheck))
			saveCfg("backup.settings.confirm_purge", boolStr(s.ConfirmPurge))
		}
	}
}

// expandEnv 把 ${VAR} 替换成 os.Getenv("VAR")；找不到时留空（不报错）。
func expandEnv(s string) string {
	if !strings.Contains(s, "${") {
		return s
	}
	return os.Expand(s, func(k string) string {
		// os.Expand 已经剥掉了 ${}  ; 直接查环境变量
		return os.Getenv(k)
	})
}

// 触发 seeding：在 InitAPI 里调用。如果 db 还没好就跳过。
func triggerSeedOnce(db *gorm.DB) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[seed] panic: %v", r)
		}
	}()
	// 用绝对路径定位 config/ 目录（main.go 在 cmd/server，cwd=/root 或 /app）
	candidates := []string{
		"config/backup-seed.json",
		"/app/config/backup-seed.json",
	}
	_ = filepath.Glob // 占位，避免某些静态检查误删 import
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			seedBackupConfig(db)
			return
		}
	}
	// 没找到也无所谓，用户可在 UI 里手动建
	_ = fmt.Sprintf
}
