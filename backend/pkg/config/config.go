package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// 配置管理器
type Config struct {
	mu       sync.RWMutex
	data     map[string]interface{}
	filePath string
}

var (
	instance *Config
	once     sync.Once
)

// 获取配置实例
func GetInstance() *Config {
	once.Do(func() {
		instance = &Config{
			data:     make(map[string]interface{}),
			filePath: getConfigPath(),
		}
		instance.load()
	})
	return instance
}

// 获取配置文件路径
func getConfigPath() string {
	// 可以从环境变量获取
	if path := os.Getenv("CONFIG_FILE"); path != "" {
		return path
	}

	// 默认路径
	return "/etc/nas-dashboard/config.json"
}

// 加载配置
func (c *Config) load() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 检查文件是否存在
	if _, err := os.Stat(c.filePath); os.IsNotExist(err) {
		// 文件不存在，创建默认配置
		c.data = c.getDefaultConfig()
		return c.save()
	}

	// 读取文件
	data, err := os.ReadFile(c.filePath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// 解析JSON
	if err := json.Unmarshal(data, &c.data); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	return nil
}

// 保存配置
func (c *Config) save() error {
	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(c.filePath), 0755); err != nil {
		return err
	}

	// 序列化配置
	data, err := json.MarshalIndent(c.data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(c.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// 获取默认配置
func (c *Config) getDefaultConfig() map[string]interface{} {
	return map[string]interface{}{
		"server": map[string]interface{}{
			"host": "0.0.0.0",
			"port": 8888,
			"tls":  false,
		},
		"database": map[string]interface{}{
			"type":     "sqlite",
			"path":     "/var/lib/nas-dashboard/data.db",
			"timeout":  10,
		},
		"logging": map[string]interface{}{
			"level":    "info",
			"file":     "/var/log/nas-dashboard/app.log",
			"max_size": 100, // MB
			"max_age":  30,  // days
		},
		"auth": map[string]interface{}{
			"jwt_secret":     "change-me-in-production",
			"token_duration": 86400, // 24 hours
			"refresh_duration": 604800, // 7 days
		},
		"websocket": map[string]interface{}{
			"enabled":       true,
			"ping_interval": 30,
			"ping_timeout":  60,
		},
		"monitoring": map[string]interface{}{
			"enabled":         true,
			"update_interval": 2, // seconds
			"history_days":    7,
		},
	}
}

// 获取字符串配置
func (c *Config) GetString(key string, defaultValue string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if value, ok := c.data[key].(string); ok {
		return value
	}
	return defaultValue
}

// 获取整数配置
func (c *Config) GetInt(key string, defaultValue int) int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if value, ok := c.data[key].(int); ok {
		return value
	}
	if value, ok := c.data[key].(float64); ok {
		return int(value)
	}
	return defaultValue
}

// 获取布尔配置
func (c *Config) GetBool(key string, defaultValue bool) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if value, ok := c.data[key].(bool); ok {
		return value
	}
	return defaultValue
}

// 获取嵌套配置
func (c *Config) GetNested(keys ...string) interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	current := c.data
	for _, key := range keys {
		if value, ok := current[key].(map[string]interface{}); ok {
			current = value
		} else {
			return nil
		}
	}
	return current
}

// 设置配置
func (c *Config) Set(key string, value interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = value
	return c.save()
}

// 设置嵌套配置
func (c *Config) SetNested(value interface{}, keys ...string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	current := c.data
	for i, key := range keys {
		if i == len(keys)-1 {
			// 最后一个key，设置值
			current[key] = value
		} else {
			// 中间key，确保是map
			if _, ok := current[key].(map[string]interface{}); !ok {
				current[key] = make(map[string]interface{})
			}
			current = current[key].(map[string]interface{})
		}
	}

	return c.save()
}

// 获取所有配置
func (c *Config) GetAll() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// 返回副本
	result := make(map[string]interface{})
	for k, v := range c.data {
		result[k] = v
	}
	return result
}

// 重新加载配置
func (c *Config) Reload() error {
	return c.load()
}

// 便捷方法 - 获取服务器配置
func (c *Config) GetServerConfig() (host string, port int, tls bool) {
	serverConfig := c.GetNested("server").(map[string]interface{})
	host = c.GetString("server.host", "0.0.0.0")
	port = c.GetInt("server.port", 8888)
	tls = c.GetBool("server.tls", false)
	return
}

// 获取数据库配置
func (c *Config) GetDatabaseConfig() map[string]interface{} {
	return c.GetNested("database").(map[string]interface{})
}

// 获取日志配置
func (c *Config) GetLoggingConfig() map[string]interface{} {
	return c.GetNested("logging").(map[string]interface{})
}

// 获取JWT配置
func (c *Config) GetJWTConfig() (secret string, tokenDuration, refreshDuration int) {
	secret = c.GetString("auth.jwt_secret", "change-me-in-production")
	tokenDuration = c.GetInt("auth.token_duration", 86400)
	refreshDuration = c.GetInt("auth.refresh_duration", 604800)
	return
}