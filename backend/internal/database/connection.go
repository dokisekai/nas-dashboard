package database

import (
	"fmt"
	"log"
	"nas-dashboard/internal/models"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Config 数据库配置
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// LoadConfig 从环境变量加载数据库配置
func LoadConfig() *Config {
	// 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "nasdashboard"),
		Password: getEnv("DB_PASSWORD", ""),
		DBName:   getEnv("DB_NAME", "nasdashboard"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}

// Connect 连接到数据库
func Connect(cfg *Config) error {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// 配置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	// 设置空闲连接池中的最大连接数
	sqlDB.SetMaxIdleConns(10)

	// 设置数据库的最大打开连接数
	sqlDB.SetMaxOpenConns(100)

	// 设置连接可重用的最长时间
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Successfully connected to database")
	return nil
}

// Migrate 运行数据库迁移
func Migrate() error {
	if DB == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	log.Println("Running database migrations...")

	// 自动迁移所有模型
	err := DB.AutoMigrate(
		&models.User{},
		&models.SSHKey{},
		&models.Session{},
		&models.SystemConfig{},
		&models.BackupRecord{},
		&models.OperationLog{},
		&models.Plugin{},
		&models.FileSystemAccess{},
	)

	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	// 创建默认管理员用户（如果不存在）
	if err := createDefaultAdmin(); err != nil {
		log.Printf("Warning: failed to create default admin: %v", err)
	}

	// 创建默认系统配置
	if err := createDefaultConfigs(); err != nil {
		log.Printf("Warning: failed to create default configs: %v", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// createDefaultAdmin 创建默认管理员用户
func createDefaultAdmin() error {
	var count int64
	DB.Model(&models.User{}).Where("role = ?", "admin").Count(&count)

	if count > 0 {
		return nil // 已经存在管理员用户
	}

	// 创建默认管理员
	admin := &models.User{
		Username:    "admin",
		DisplayName: "Administrator",
		Email:       "admin@localhost",
		Role:        "admin",
		IsActive:    true,
	}

	// 注意：这里应该在实际使用时调用 AuthService 来设置密码哈希
	// 这里只是创建基础结构
	if err := DB.Create(admin).Error; err != nil {
		return err
	}

	log.Println("Default admin user created")
	return nil
}

// createDefaultConfigs 创建默认系统配置
func createDefaultConfigs() error {
	configs := []models.SystemConfig{
		{
			Key:         "system.name",
			Value:       "NAS Dashboard",
			Type:        "string",
			Category:    "general",
			Description: "系统名称",
			IsPublic:    true,
		},
		{
			Key:         "system.timezone",
			Value:       "Asia/Shanghai",
			Type:        "string",
			Category:    "general",
			Description: "系统时区",
			IsPublic:    true,
		},
		{
			Key:         "security.session_timeout",
			Value:       "86400",
			Type:        "int",
			Category:    "security",
			Description: "会话超时时间（秒）",
			IsPublic:    false,
		},
		{
			Key:         "security.max_login_attempts",
			Value:       "5",
			Type:        "int",
			Category:    "security",
			Description: "最大登录尝试次数",
			IsPublic:    false,
		},
		{
			Key:         "backup.auto_backup_enabled",
			Value:       "true",
			Type:        "bool",
			Category:    "backup",
			Description: "启用自动备份",
			IsPublic:    false,
		},
		{
			Key:         "backup.retention_days",
			Value:       "30",
			Type:        "int",
			Category:    "backup",
			Description: "备份保留天数",
			IsPublic:    false,
		},
	}

	for _, config := range configs {
		var existingConfig models.SystemConfig
		if err := DB.Where("key = ?", config.Key).First(&existingConfig).Error; err == gorm.ErrRecordNotFound {
			if err := DB.Create(&config).Error; err != nil {
				log.Printf("Failed to create config %s: %v", config.Key, err)
			}
		}
	}

	return nil
}

// Close 关闭数据库连接
func Close() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}
