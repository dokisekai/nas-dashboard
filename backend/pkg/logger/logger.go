package logger

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

// 日志级别
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

var (
	logger       *log.Logger
	currentLevel LogLevel = INFO
	logFile      *os.File
)

// 初始化日志系统
func Init(logPath string, level LogLevel) error {
	currentLevel = level

	// 确保日志目录存在
	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
		return err
	}

	// 打开日志文件
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	logFile = file
	logger = log.New(file, "", log.LstdFlags)

	// 同时输出到控制台
	log.SetOutput(file)

	Info("Logger initialized", map[string]interface{}{
		"level": levelToString(currentLevel),
		"path":  logPath,
	})

	return nil
}

// 关闭日志系统
func Close() {
	if logFile != nil {
		logFile.Close()
	}
}

// 日志格式化
type LogFields map[string]interface{}

func formatLog(level string, message string, fields LogFields) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logStr := timestamp + " [" + level + "] " + message

	if len(fields) > 0 {
		logStr += " |"
		for key, value := range fields {
			logStr += " " + key + "=" + formatValue(value)
		}
	}

	return logStr
}

func formatValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case error:
		return v.Error()
	default:
		return ""
	}
}

func levelToString(level LogLevel) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// 日志方法
func Debug(message string, fields LogFields) {
	if currentLevel <= DEBUG {
		logStr := formatLog(levelToString(DEBUG), message, fields)
		if logger != nil {
			logger.Println(logStr)
		}
	}
}

func Info(message string, fields LogFields) {
	if currentLevel <= INFO {
		logStr := formatLog(levelToString(INFO), message, fields)
		if logger != nil {
			logger.Println(logStr)
		}
	}
}

func Warn(message string, fields LogFields) {
	if currentLevel <= WARN {
		logStr := formatLog(levelToString(WARN), message, fields)
		if logger != nil {
			logger.Println(logStr)
		}
	}
}

func Error(message string, fields LogFields) {
	if currentLevel <= ERROR {
		logStr := formatLog(levelToString(ERROR), message, fields)
		if logger != nil {
			logger.Println(logStr)
		}
	}
}

func Fatal(message string, fields LogFields) {
	logStr := formatLog(levelToString(FATAL), message, fields)
	if logger != nil {
		logger.Println(logStr)
	}
	os.Exit(1)
}

// API日志专用方法
func APIRequest(method, path string, fields LogFields) {
	Info("API Request", mergeFields(fields, LogFields{
		"method": method,
		"path":   path,
	}))
}

func APIResponse(code int, message string, duration int64) {
	Info("API Response", LogFields{
		"code":     code,
		"message":  message,
		"duration": duration,
	})
}

func APIError(method, path string, err error, code int) {
	Error("API Error", LogFields{
		"method": method,
		"path":   path,
		"error":  err.Error(),
		"code":   code,
	})
}

// 辅助方法
func mergeFields(fields ...LogFields) LogFields {
	result := make(LogFields)
	for _, f := range fields {
		for k, v := range f {
			result[k] = v
		}
	}
	return result
}