package api

import (
	"log"
	"nas-dashboard/internal/models"
	"time"
	"gorm.io/gorm"
)

type Scheduler struct {
	DB      *gorm.DB
	SyncAPI *SyncAPI
}

func NewScheduler(db *gorm.DB) *Scheduler {
	return &Scheduler{
		DB:      db,
		SyncAPI: NewSyncAPI(db),
	}
}

func (s *Scheduler) Start() {
	log.Println("Starting system task scheduler...")
	
	// 每分钟检查一次任务
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for range ticker.C {
			s.checkSyncJobs()
			s.checkBackupTasks()
		}
	}()
}

func (s *Scheduler) checkSyncJobs() {
	var jobs []models.SyncJob
	// 查找启用的且未在运行的任务
	if err := s.DB.Where("enabled = ? AND status != ?", true, "running").Find(&jobs).Error; err != nil {
		return
	}

	for _, job := range jobs {
		if s.shouldRun(job.LastRun, job.Schedule) {
			log.Printf("Triggering scheduled sync job: %s", job.Name)
			go s.SyncAPI.executeSync(job)
		}
	}
}

func (s *Scheduler) checkBackupTasks() {
	// 类似逻辑处理 BackupTasks
}

// shouldRun 简单的计划解析
func (s *Scheduler) shouldRun(lastRun *time.Time, schedule string) bool {
	if schedule == "" {
		return false
	}
	
	now := time.Now()
	
	// 如果从未运行过，立即运行
	if lastRun == nil {
		return true
	}

	// 极简实现：如果是 "daily" 且距离上次运行超过 24 小时
	if schedule == "daily" && now.Sub(*lastRun) >= 24*time.Hour {
		return true
	}
	
	// 如果是 "hourly"
	if schedule == "hourly" && now.Sub(*lastRun) >= 1*time.Hour {
		return true
	}

	// 实际生产环境应使用 cron 库解析 schedule 字符串
	return false
}
