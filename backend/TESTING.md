# NAS Dashboard Backend 测试指南

## 测试概述

本指南提供了 NAS Dashboard 后端的完整测试策略，包括单元测试、集成测试和 API 测试。

## 测试环境设置

### 前置条件
```bash
# 安装测试依赖
go get github.com/stretchr/testify/assert
go get github.com/stretchr/testify/mock
go get github.com/golang/mock/gomock

# 安装测试工具
go install github.com/golang/mock/mockgen@latest
```

### 测试数据库设置
```bash
# 创建测试数据库
sudo -u postgres createdb nasdashboard_test

# 设置环境变量
export DB_NAME=nasdashboard_test
export DB_TEST_MODE=true
```

## 单元测试

### 认证服务测试
```go
// internal/service/auth_test.go
package service

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
    authService := NewAuthService()

    // 测试密码哈希
    hashedPassword, err := authService.HashPassword("testPassword123")
    assert.NoError(t, err)
    assert.NotEmpty(t, hashedPassword)
    assert.NotEqual(t, "testPassword123", hashedPassword)

    // 测试密码验证
    isValid := authService.VerifyPassword(hashedPassword, "testPassword123")
    assert.True(t, isValid)

    // 测试错误密码
    isValid = authService.VerifyPassword(hashedPassword, "wrongPassword")
    assert.False(t, isValid)
}

func TestGenerateToken(t *testing.T) {
    authService := NewAuthService()
    user := &models.User{
        Username: "testuser",
        Role:     "user",
    }

    // 测试访问令牌生成
    accessToken, err := authService.GenerateAccessToken(user)
    assert.NoError(t, err)
    assert.NotEmpty(t, accessToken)

    // 测试令牌验证
    claims, err := authService.ValidateAccessToken(accessToken)
    assert.NoError(t, err)
    assert.Equal(t, "testuser", claims.Username)
    assert.Equal(t, "user", claims.Role)
}

func TestLogin(t *testing.T) {
    // 创建测试用户
    authService := NewAuthService()
    user, err := authService.CreateUser("testuser", "password123", "test@example.com", "Test User", "user")
    assert.NoError(t, err)

    // 测试成功登录
    loggedInUser, accessToken, refreshToken, err := authService.Login("testuser", "password123", "127.0.0.1", "test-agent")
    assert.NoError(t, err)
    assert.Equal(t, user.Username, loggedInUser.Username)
    assert.NotEmpty(t, accessToken)
    assert.NotEmpty(t, refreshToken)

    // 测试错误密码
    _, _, _, err = authService.Login("testuser", "wrongpassword", "127.0.0.1", "test-agent")
    assert.Error(t, err)
    assert.Equal(t, ErrInvalidPassword, err)

    // 测试不存在的用户
    _, _, _, err = authService.Login("nonexistent", "password", "127.0.0.1", "test-agent")
    assert.Error(t, err)
    assert.Equal(t, ErrUserNotFound, err)
}
```

### 文件服务测试
```go
// internal/api/file_management_test.go
package api

import (
    "testing"
    "os"
    "path/filepath"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestValidatePath(t *testing.T) {
    fileService := NewFileService()

    // 创建测试目录
    testDir := "/tmp/nas-test"
    os.MkdirAll(testDir, 0755)
    defer os.RemoveAll(testDir)

    // 测试有效路径
    cleanPath, err := fileService.ValidatePath(testDir)
    assert.NoError(t, err)
    assert.Equal(t, testDir, cleanPath)

    // 测试路径遍历攻击
    _, err = fileService.ValidatePath(testDir + "/../etc")
    assert.Error(t, err)

    // 测试不存在的路径
    _, err = fileService.ValidatePath("/nonexistent/path")
    assert.Error(t, err)
}

func TestCheckPermission(t *testing.T) {
    fileService := NewFileService()

    // 创建模拟上下文
    c, _ := gin.CreateTestContext(nil)
    c.Set("username", "testuser")
    c.Set("role", "user")

    // 测试用户主目录访问
    userHomeDir := "/home/testuser"
    err := fileService.CheckPermission(c, userHomeDir+"/documents", "read")
    assert.NoError(t, err)

    // 测试越权访问
    err = fileService.CheckPermission(c, "/root", "read")
    assert.Error(t, err)
}
```

### 配置服务测试
```go
// internal/api/system_config_test.go
package api

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestValidateConfig(t *testing.T) {
    configService := NewConfigService()

    // 测试整数验证
    err := configService.ValidateConfig("int", "12345")
    assert.NoError(t, err)

    err = configService.ValidateConfig("int", "not_a_number")
    assert.Error(t, err)

    // 测试布尔值验证
    err = configService.ValidateConfig("bool", "true")
    assert.NoError(t, err)

    err = configService.ValidateConfig("bool", "invalid")
    assert.Error(t, err)

    // 测试 JSON 验证
    err = configService.ValidateConfig("json", `{"key": "value"}`)
    assert.NoError(t, err)

    err = configService.ValidateConfig("json", "not_json")
    assert.Error(t, err)
}
```

## 集成测试

### API 集成测试
```go
// internal/api/integration_test.go
package api

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func setupTestRouter() *gin.Engine {
    gin.SetMode(gin.TestMode)
    router := gin.New()

    // 设置测试路由
    apiGroup := router.Group("/api")
    {
        apiGroup.POST("/auth/login", Login)
        apiGroup.POST("/auth/refresh", RefreshToken)
    }

    return router
}

func TestLoginIntegration(t *testing.T) {
    // 设置测试数据库
    setupTestDatabase()
    defer cleanupTestDatabase()

    // 创建测试用户
    authService := NewAuthService()
    _, err := authService.CreateUser("testuser", "password123", "test@example.com", "Test User", "user")
    require.NoError(t, err)

    // 设置路由
    router := setupTestRouter()

    // 测试登录
    loginReq := LoginRequest{
        Username: "testuser",
        Password: "password123",
    }
    body, _ := json.Marshal(loginReq)

    req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)

    var response LoginResponse
    err = json.Unmarshal(w.Body.Bytes(), &response)
    require.NoError(t, err)

    assert.NotEmpty(t, response.Token)
    assert.NotEmpty(t, response.RefreshToken)
    assert.Equal(t, "testuser", response.User.Username)
}

func TestProtectedRouteIntegration(t *testing.T) {
    // 设置测试数据库和用户
    setupTestDatabase()
    defer cleanupTestDatabase()

    authService := NewAuthService()
    user, err := authService.CreateUser("admin", "admin123", "admin@example.com", "Admin", "admin")
    require.NoError(t, err)

    // 获取访问令牌
    _, accessToken, _, err := authService.Login("admin", "admin123", "127.0.0.1", "test-agent")
    require.NoError(t, err)

    // 设置受保护路由
    router := gin.New()
    protected := router.Group("/api")
    protected.Use(Auth())
    {
        protected.GET("/config", GetConfigs)
    }

    // 测试未授权访问
    req, _ := http.NewRequest("GET", "/api/config", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    assert.Equal(t, http.StatusUnauthorized, w.Code)

    // 测试授权访问
    req, _ = http.NewRequest("GET", "/api/config", nil)
    req.Header.Set("Authorization", "Bearer "+accessToken)
    w = httptest.NewRecorder()
    router.ServeHTTP(w, req)
    assert.Equal(t, http.StatusOK, w.Code)
}
```

### 数据库集成测试
```go
// internal/database/integration_test.go
package database

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestDatabaseMigration(t *testing.T) {
    // 设置测试数据库
    cfg := &Config{
        Host:     "localhost",
        Port:     "5432",
        User:     "postgres",
        Password: "",
        DBName:   "nasdashboard_test",
        SSLMode:  "disable",
    }

    err := Connect(cfg)
    require.NoError(t, err)
    defer Close()

    // 运行迁移
    err = Migrate()
    require.NoError(t, err)

    // 检查表是否创建
    db := GetDB()
    var count int64

    // 检查用户表
    db.Table("users").Count(&count)
    assert.Greater(t, count, int64(0)) // 默认管理员用户

    // 检查配置表
    db.Table("system_configs").Count(&count)
    assert.Greater(t, count, int64(0)) // 默认配置
}

func TestDatabaseCRUD(t *testing.T) {
    // 设置测试数据库
    setupTestDatabase()
    defer cleanupTestDatabase()

    db := GetDB()

    // 创建测试记录
    config := &models.SystemConfig{
        Key:         "test.config",
        Value:       "test_value",
        Type:        "string",
        Category:    "test",
        Description: "Test configuration",
        IsPublic:    true,
    }

    err := db.Create(config).Error
    require.NoError(t, err)
    assert.NotZero(t, config.ID)

    // 读取记录
    var readConfig models.SystemConfig
    err = db.First(&readConfig, config.ID).Error
    require.NoError(t, err)
    assert.Equal(t, config.Key, readConfig.Key)

    // 更新记录
    readConfig.Value = "updated_value"
    err = db.Save(&readConfig).Error
    require.NoError(t, err)

    // 验证更新
    var updatedConfig models.SystemConfig
    db.First(&updatedConfig, config.ID)
    assert.Equal(t, "updated_value", updatedConfig.Value)

    // 删除记录
    err = db.Delete(&updatedConfig).Error
    require.NoError(t, err)

    // 验证删除
    var deletedConfig models.SystemConfig
    err = db.First(&deletedConfig, config.ID).Error
    assert.Error(t, err)
}
```

## API 测试

### 使用 curl 测试
```bash
# 测试健康检查
curl http://localhost:8888/health

# 测试登录
curl -X POST http://localhost:8888/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

# 测试受保护路由
TOKEN="eyJhbGciOiJIUzI1NiIs..."
curl -X GET http://localhost:8888/api/config \
  -H "Authorization: Bearer $TOKEN"

# 测试文件上传
curl -X POST http://localhost:8888/api/files/upload \
  -H "Authorization: Bearer $TOKEN" \
  -F "path=/home/user/documents" \
  -F "file=@/path/to/file.pdf"

# 测试备份创建
curl -X POST http://localhost:8888/api/backups \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "test-backup",
    "type": "full",
    "includeDB": true,
    "includeFiles": true,
    "filePaths": ["/etc", "/home"]
  }'
```

### 使用 Postman 测试
导入以下 Postman 集合：

```json
{
  "info": {
    "name": "NAS Dashboard API",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "item": [
    {
      "name": "Authentication",
      "item": [
        {
          "name": "Login",
          "request": {
            "method": "POST",
            "header": [],
            "body": {
              "mode": "raw",
              "raw": "{\n  \"username\": \"admin\",\n  \"password\": \"admin123\"\n}"
            },
            "url": {
              "raw": "http://localhost:8888/api/auth/login",
              "protocol": "http",
              "host": ["localhost"],
              "port": "8888",
              "path": ["api", "auth", "login"]
            }
          }
        }
      ]
    },
    {
      "name": "File Management",
      "item": [
        {
          "name": "List Files",
          "request": {
            "method": "POST",
            "header": [],
            "body": {
              "mode": "raw",
              "raw": "{\n  \"path\": \"/home/user\"\n}"
            },
            "url": {
              "raw": "http://localhost:8888/api/files/list",
              "protocol": "http",
              "host": ["localhost"],
              "port": "8888",
              "path": ["api", "files", "list"]
            }
          }
        }
      ]
    }
  ]
}
```

## WebSocket 测试

### 使用 wscat 测试
```bash
# 安装 wscat
npm install -g wscat

# 连接 WebSocket
wscat -c "ws://localhost:8888/ws/monitor"

# 发送消息
{"type":"ping"}

# 订阅监控数据
{"type":"subscribe","data":{"channels":["monitor"]}}
```

### 使用 JavaScript 测试
```javascript
// test-websocket.js
const WebSocket = require('ws');

const ws = new WebSocket('ws://localhost:8888/ws/monitor');

ws.on('open', function open() {
  console.log('WebSocket connected');

  // 发送心跳
  ws.send(JSON.stringify({
    type: 'ping'
  }));
});

ws.on('message', function message(data) {
  const msg = JSON.parse(data);
  console.log('Received:', msg);

  if (msg.type === 'monitor_data') {
    console.log('CPU:', msg.data.cpu);
    console.log('Memory:', msg.data.memory);
  }
});

ws.on('error', function error(error) {
  console.error('WebSocket error:', error);
});

ws.on('close', function close() {
  console.log('WebSocket disconnected');
});
```

## 性能测试

### 使用 Apache Bench
```bash
# 安装 Apache Bench
sudo apt-get install apache2-utils

# 测试登录端点
ab -n 1000 -c 10 -p login.json -T application/json \
  http://localhost:8888/api/auth/login

# 测试监控数据端点
TOKEN="eyJhbGciOiJIUzI1NiIs..."
ab -n 1000 -c 10 \
  -H "Authorization: Bearer $TOKEN" \
  http://localhost:8888/api/monitor/cpu
```

### 使用 wrk
```bash
# 安装 wrk
git clone https://github.com/wg/wrk.git
cd wrk && make

# 运行性能测试
./wrk -t4 -c100 -d30s --latency \
  -H "Authorization: Bearer $TOKEN" \
  http://localhost:8888/api/config
```

## 负载测试

### 数据库连接池测试
```go
func TestDatabaseConnectionPool(t *testing.T) {
    setupTestDatabase()
    defer cleanupTestDatabase()

    var wg sync.WaitGroup
    connections := 100

    for i := 0; i < connections; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()

            db := GetDB()
            var config models.SystemConfig
            err := db.First(&config).Error

            assert.NoError(t, err)
        }()
    }

    wg.Wait()
}
```

### 并发文件操作测试
```go
func TestConcurrentFileOperations(t *testing.T) {
    fileService := NewFileService()
    testDir := "/tmp/nas-concurrent-test"
    os.MkdirAll(testDir, 0755)
    defer os.RemoveAll(testDir)

    var wg sync.WaitGroup
    operations := 50

    for i := 0; i < operations; i++ {
        wg.Add(1)
        go func(index int) {
            defer wg.Done()

            // 创建文件
            fileName := fmt.Sprintf("test_file_%d.txt", index)
            filePath := filepath.Join(testDir, fileName)
            content := []byte("test content")

            err := os.WriteFile(filePath, content, 0644)
            assert.NoError(t, err)

            // 读取文件
            readContent, err := os.ReadFile(filePath)
            assert.NoError(t, err)
            assert.Equal(t, content, readContent)
        }(i)
    }

    wg.Wait()
}
```

## 测试覆盖率

### 生成覆盖率报告
```bash
# 运行测试并生成覆盖率
go test -coverprofile=coverage.out ./...

# 查看覆盖率
go tool cover -func=coverage.out

# 生成 HTML 报告
go tool cover -html=coverage.out -o coverage.html

# 在浏览器中查看
firefox coverage.html
```

### 设置覆盖率目标
```bash
# 检查覆盖率是否达到目标
go test -coverprofile=coverage.out ./...
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
if (( $(echo "$COVERAGE < 80" | bc -l) )); then
    echo "Coverage $COVERAGE% is below 80%"
    exit 1
fi
```

## 持续集成

### GitHub Actions 配置
```yaml
# .github/workflows/test.yml
name: Test

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]

jobs:
  test:
    runs-on: ubuntu-latest

    services:
      postgres:
        image: postgres:12
        env:
          POSTGRES_USER: postgres
          POSTGRES_PASSWORD: postgres
          POSTGRES_DB: nasdashboard_test
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5

    steps:
    - uses: actions/checkout@v2

    - name: Set up Go
      uses: actions/setup-go@v2
      with:
        go-version: 1.25.0

    - name: Install dependencies
      run: go mod download

    - name: Run tests
      run: go test -v -race -coverprofile=coverage.out ./...
      env:
        DB_HOST: localhost
        DB_PORT: 5432
        DB_USER: postgres
        DB_PASSWORD: postgres
        DB_NAME: nasdashboard_test
        DB_SSLMODE: disable

    - name: Upload coverage
      uses: codecov/codecov-action@v2
      with:
        files: coverage.out
```

## 故障排除

### 常见测试问题

1. **数据库连接失败**
```bash
# 检查 PostgreSQL 是否运行
sudo systemctl status postgresql

# 重启 PostgreSQL
sudo systemctl restart postgresql
```

2. **端口占用**
```bash
# 查找占用端口的进程
lsof -i :8888

# 终止进程
kill -9 <PID>
```

3. **权限错误**
```bash
# 设置正确的权限
chmod +x ./scripts/test-setup.sh
sudo ./scripts/test-setup.sh
```

## 最佳实践

### 测试命名规范
```go
// 好的测试名称
func TestUserLoginWithValidCredentials(t *testing.T) {}
func TestUserLoginWithInvalidPassword(t *testing.T) {}

// 避免模糊的名称
func TestLogin1(t *testing.T) {}
func TestLogin2(t *testing.T) {}
```

### 测试结构
```go
func TestFeature(t *testing.T) {
    // Arrange
    setupTestEnvironment()
    testData := createTestData()

    // Act
    result, err := functionUnderTest(testData)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, expected, result)

    // Cleanup
    cleanupTestEnvironment()
}
```

### 表驱动测试
```go
func TestValidateConfig(t *testing.T) {
    tests := []struct {
        name    string
        config  string
        type    string
        wantErr bool
    }{
        {"valid int", "123", "int", false},
        {"invalid int", "abc", "int", true},
        {"valid bool", "true", "bool", false},
        {"invalid bool", "yes", "bool", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := configService.ValidateConfig(tt.type, tt.config)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

这个测试指南提供了全面的测试策略，确保 NAS Dashboard 后端的质量和可靠性。
