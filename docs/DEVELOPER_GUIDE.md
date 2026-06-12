# NAS Dashboard - Developer Guide

Complete development guide for contributors to the NAS Dashboard project, including setup, workflows, and best practices.

## 📋 Table of Contents

1. [Development Overview](#development-overview)
2. [Development Environment Setup](#development-environment-setup)
3. [Project Structure](#project-structure)
4. [Development Workflow](#development-workflow)
5. [Coding Standards](#coding-standards)
6. [Testing Guidelines](#testing-guidelines)
7. [Debugging Techniques](#debugging-techniques)
8. [Performance Optimization](#performance-optimization)
9. [Contributing Guidelines](#contributing-guidelines)
10. [Release Process](#release-process)

---

## Development Overview

### Technology Stack

#### Frontend
- **Framework**: Vue 3.5+ with Composition API
- **Language**: TypeScript 6.0+
- **Build Tool**: Vite 8.0+
- **State Management**: Pinia 3.0+
- **Styling**: Tailwind CSS 4.3+
- **Charts**: Chart.js 4.5+, ApexCharts 5.15+

#### Backend
- **Language**: Go 1.22+
- **Framework**: Gin (web framework)
- **Database**: PostgreSQL with GORM
- **Authentication**: JWT with bcrypt
- **System Monitoring**: gopsutil
- **Container Management**: Docker SDK

### Development Principles

1. **Clean Code**: Write readable, maintainable code
2. **Testing**: Test everything you commit
3. **Documentation**: Document as you code
4. **Performance**: Optimize for speed and efficiency
5. **Security**: Follow security best practices
6. **Collaboration**: Work well with others

---

## Development Environment Setup

### Prerequisites

#### Required Software

- **Go**: 1.22+ ([Download](https://golang.org/dl/))
- **Node.js**: 20+ ([Download](https://nodejs.org/))
- **Git**: Latest version
- **Docker**: 20.10+ (for container testing)
- **VS Code**: Recommended IDE

#### VS Code Extensions

Recommended extensions for development:

```json
{
  "recommendations": [
    "vue.volar",
    "vue.vscode-typescript-vue-plugin",
    "golang.go",
    "bradlc.vscode-tailwindcss",
    "esbenp.prettier-vscode",
    "dbaeumer.vscode-eslint",
    "ms-vscode.makefile-tools"
  ]
}
```

### Initial Setup

#### 1. Clone Repository

```bash
# Clone the repository
git clone https://github.com/yourusername/nas-dashboard.git
cd nas-dashboard

# Create development branch
git checkout -b develop
```

#### 2. Backend Setup

```bash
cd backend

# Install Go dependencies
go mod download

# Install development tools
go install github.com/cosmtrek/air@latest  # Live reload
go install github.com/git Cosby/cosetty@latest  # API testing

# Verify installation
go version
air --version
```

**Create development configuration:**

```bash
# backend/config/dev.env
PORT=8888
GIN_MODE=debug
JWT_SECRET=dev-secret-key
JWT_ACCESS_DURATION=24h
JWT_REFRESH_DURATION=720h
CORS_ORIGIN=http://localhost:5173
DB_HOST=localhost
DB_PORT=5432
DB_NAME=nas_dashboard_dev
DB_USER=nas_user
DB_PASSWORD=dev_password
LOG_LEVEL=debug
```

#### 3. Frontend Setup

```bash
cd frontend

# Install dependencies
npm install

# Verify installation
npm --version
node --version
```

**Create development configuration:**

```bash
# frontend/.env
VITE_API_URL=http://localhost:8888
VITE_WS_URL=ws://localhost:8888
VITE_DEBUG=true
```

#### 4. Database Setup

**Option A: PostgreSQL (Recommended)**

```bash
# Install PostgreSQL
sudo apt update
sudo apt install postgresql postgresql-contrib

# Create database
sudo -u postgres psql << EOF
CREATE DATABASE nas_dashboard_dev;
CREATE USER nas_user WITH PASSWORD 'dev_password';
GRANT ALL PRIVILEGES ON DATABASE nas_dashboard_dev TO nas_user;
EOF
```

**Option B: SQLite (Quick Development)**

```bash
# No setup required - will use embedded SQLite
# Configure backend to use SQLite
echo "DB_TYPE=sqlite" >> backend/config/dev.env
echo "DB_PATH=./dev.db" >> backend/config/dev.env
```

### Start Development Servers

#### Backend Development Server

```bash
cd backend

# Option 1: Manual restart
go run cmd/server/main.go

# Option 2: Live reload (recommended)
air

# Backend will run on http://localhost:8888
```

#### Frontend Development Server

```bash
cd frontend

# Start development server
npm run dev

# Frontend will run on http://localhost:5173
```

### Verify Setup

**Test Backend:**

```bash
# Check health endpoint
curl http://localhost:8888/api/health

# Should return:
# {"success":true,"data":{"status":"healthy","timestamp":"...","version":"0.1.0"}}
```

**Test Frontend:**

```bash
# Open browser
open http://localhost:5173

# Should see login screen
```

---

## Project Structure

### Overall Structure

```
nas-dashboard/
├── backend/                 # Go backend application
│   ├── cmd/                # Application entry points
│   │   └── server/
│   │       └── main.go    # Application entry point
│   ├── internal/           # Private application code
│   │   ├── api/          # API handlers
│   │   ├── middleware/   # HTTP middleware
│   │   ├── models/       # Data models
│   │   ├── services/     # Business logic
│   │   └── websocket/   # WebSocket handlers
│   ├── pkg/              # Public packages
│   ├── config/          # Configuration files
│   ├── go.mod           # Go dependencies
│   └── go.sum           # Dependency checksums
├── frontend/             # Vue 3 frontend
│   ├── src/
│   │   ├── api/         # API clients
│   │   ├── apps/        # Desktop applications
│   │   ├── components/  # Vue components
│   │   ├── stores/      # Pinia stores
│   │   ├── composables/ # Composition functions
│   │   ├── types/       # TypeScript definitions
│   │   └── utils/       # Utility functions
│   ├── public/          # Static files
│   ├── package.json     # NPM dependencies
│   ├── vite.config.ts   # Vite configuration
│   └── tsconfig.json    # TypeScript configuration
├── docs/                # Documentation
├── docker-compose.yml   # Docker orchestration
└── README.md           # Project overview
```

### Backend Structure Details

```
backend/
├── cmd/server/
│   └── main.go                 # Application entry point
├── internal/
│   ├── api/
│   │   ├── auth.go            # Authentication endpoints
│   │   ├── monitor.go         # System monitoring
│   │   ├── storage.go         # Storage management
│   │   ├── docker.go          # Docker management
│   │   ├── user.go            # User management
│   │   └── service.go         # Service management
│   ├── middleware/
│   │   ├── auth.go            # JWT authentication
│   │   ├── cors.go            # CORS handling
│   │   └── logger.go          # Request logging
│   ├── models/
│   │   ├── user.go            # User model
│   │   ├── database.go        # Database models
│   │   └── system.go          # System models
│   ├── services/
│   │   ├── auth_service.go    # Authentication logic
│   │   ├── monitor_service.go # Monitoring logic
│   │   └── storage_service.go # Storage logic
│   └── websocket/
│       ├── hub.go             # Connection pool
│       └── manager.go         # WebSocket manager
├── pkg/
│   ├── jwt/
│   │   └── jwt.go             # JWT utilities
│   └── response/
│       └── response.go        # API response helpers
└── config/
    ├── dev.env               # Development config
    └── prod.env              # Production config
```

### Frontend Structure Details

```
frontend/
├── src/
│   ├── api/
│   │   ├── client.ts         # Axios configuration
│   │   ├── auth.ts           # Authentication API
│   │   ├── monitor.ts        # Monitoring API
│   │   ├── storage.ts        # Storage API
│   │   └── user.ts           # User API
│   ├── apps/
│   │   ├── StorageManager.vue
│   │   ├── SystemMonitor.vue
│   │   ├── UserManager.vue
│   │   ├── SystemSettings.vue
│   │   ├── AppCenter.vue
│   │   └── PluginStore.vue
│   ├── components/
│   │   ├── Desktop/
│   │   │   ├── DSMDesktop.vue
│   │   │   ├── EnhancedDock.vue
│   │   │   ├── WindowManager.vue
│   │   │   └── widgets/
│   │   │       └── SystemMonitorWidget.vue
│   │   ├── Layout/
│   │   └── Common/
│   ├── stores/
│   │   ├── auth.ts           # Authentication store
│   │   ├── app.ts            # Application state
│   │   └── monitor.ts        # Monitoring state
│   ├── composables/
│   │   ├── useDesktop.ts     # Desktop utilities
│   │   └── useWebSocket.ts  # WebSocket client
│   ├── types/
│   │   ├── desktop.ts        # Desktop types
│   │   ├── api.ts            # API types
│   │   └── index.ts          # Type exports
│   ├── utils/
│   │   ├── format.ts         # Data formatting
│   │   └── validation.ts     # Input validation
│   ├── router/
│   │   └── index.ts          # Route configuration
│   ├── App.vue               # Root component
│   └── main.ts               # Application entry
├── public/                   # Static files
├── package.json              # Dependencies
├── vite.config.ts           # Vite config
└── tsconfig.json            # TypeScript config
```

---

## Development Workflow

### Git Workflow

#### Branch Strategy

```
main (production)
  ↑
develop (development)
  ↑
feature/* (feature branches)
hotfix/* (hotfix branches)
```

#### Feature Development

1. **Create feature branch**

```bash
git checkout develop
git pull origin develop
git checkout -b feature/your-feature-name
```

2. **Make changes**

```bash
# Make your changes
git add .
git commit -m "feat: add your feature"
```

3. **Push and create PR**

```bash
git push origin feature/your-feature-name
# Create pull request on GitHub
```

#### Commit Message Format

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>[optional scope]: <description>

[optional body]

[optional footer]
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes
- `refactor`: Code refactoring
- `test`: Test additions/changes
- `chore`: Build process changes

**Examples:**

```bash
feat(auth): add refresh token support
fix(monitor): correct network data format
docs(api): update authentication documentation
refactor(dashboard): optimize CPU usage calculation
test(user): add user management tests
```

### Code Review Process

#### Submitting Pull Requests

1. **Update documentation**
   ```bash
   # Update relevant documentation
   # Add inline comments for complex code
   ```

2. **Run tests**
   ```bash
   # Backend tests
   cd backend && go test ./...
   
   # Frontend tests
   cd frontend && npm run test
   ```

3. **Create PR**
   - Clear title and description
   - Link related issues
   - Include screenshots for UI changes
   - List breaking changes

#### Review Guidelines

**For Reviewers:**

1. Check code quality and style
2. Verify tests pass
3. Test functionality manually
4. Check for security issues
5. Verify documentation is updated

**For Authors:**

1. Address all review comments
2. Update tests as needed
3. Keep PR focused and small
4. Respond to feedback promptly

---

## Coding Standards

### Go Backend Standards

#### Code Style

**Follow standard Go conventions:**

```go
// Good: Clear naming, proper error handling
func GetCPUUsage() (float64, error) {
    percent, err := cpu.Percent(0, false)
    if err != nil {
        return 0, fmt.Errorf("failed to get CPU usage: %w", err)
    }
    return percent[0], nil
}

// Bad: Unclear naming, poor error handling
func get() ([]float64, error) {
    p, e := cpu.Percent(0, false)
    if e != nil {
        return nil, e
    }
    return p, nil
}
```

#### Package Organization

```go
// Package comments
// Package api provides HTTP handlers for the NAS Dashboard API.
package api

import (
    "github.com/gin-gonic/gin"
    "nas-dashboard/internal/services"
)

// Handler groups related methods
type Handler struct {
    monitorService *services.MonitorService
    dockerService  *services.DockerService
}

// NewHandler creates a new API handler
func NewHandler(ms *services.MonitorService, ds *services.DockerService) *Handler {
    return &Handler{
        monitorService: ms,
        dockerService:  ds,
    }
}
```

#### Error Handling

```go
// Always handle errors
if err != nil {
    log.Printf("Error: %v", err)
    return err
}

// Wrap errors with context
if err := user.Save(); err != nil {
    return fmt.Errorf("failed to save user: %w", err)
}

// Use custom error types
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed for field %s: %s", e.Field, e.Message)
}
```

### Vue Frontend Standards

#### Component Structure

```vue
<template>
  <!-- Use semantic HTML -->
  <div class="cpu-monitor">
    <h2>{{ title }}</h2>
    <div class="cpu-stats">
      <!-- Clear, readable template -->
    </div>
  </div>
</template>

<script setup lang="ts">
// Import dependencies
import { ref, onMounted, computed } from 'vue'
import { useMonitorStore } from '@/stores/monitor'

// Define types
interface CPUData {
  usage: number
  cores: number
  load1: number
}

// Use Composition API
const monitorStore = useMonitorStore()
const cpuData = ref<CPUData | null>(null)

// Computed properties
const cpuPercentage = computed(() => {
  return cpuData.value ? (cpuData.value.usage * 100).toFixed(1) : '0'
})

// Lifecycle hooks
onMounted(() => {
  loadCPUData()
})

// Methods
const loadCPUData = async () => {
  try {
    cpuData.value = await monitorStore.fetchCPUData()
  } catch (error) {
    console.error('Failed to load CPU data:', error)
  }
}
</script>

<style scoped>
/* Use scoped styles */
.cpu-monitor {
  padding: 1rem;
  background: #f5f5f5;
  border-radius: 8px;
}
</style>
```

#### TypeScript Usage

```typescript
// Define interfaces for data structures
interface User {
  id: number
  username: string
  email: string
  role: 'admin' | 'user'
  createdAt: string
}

// Use generics for API responses
interface ApiResponse<T> {
  success: boolean
  data?: T
  error?: string
}

// Type-safe API calls
async function getUser(id: number): Promise<ApiResponse<User>> {
  const response = await axios.get(`/api/users/${id}`)
  return response.data
}

// Type guards for runtime validation
function isAdmin(user: User): boolean {
  return user.role === 'admin'
}
```

---

## Testing Guidelines

### Backend Testing

#### Unit Tests

```go
// internal/api/auth_test.go
package api

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestLoginHandler_ValidCredentials_ReturnsToken(t *testing.T) {
    // Setup
    mockService := new(MockAuthService)
    handler := NewAuthHandler(mockService)
    
    mockService.On("Authenticate", "admin", "admin123").
        Return("valid-token", nil)
    
    // Test
    gin.SetMode(gin.TestMode)
    router := gin.Default()
    router.POST("/login", handler.Login)
    
    req, _ := http.NewRequest("POST", "/login", 
        strings.NewReader(`{"username":"admin","password":"admin123"}`))
    req.Header.Set("Content-Type", "application/json")
    
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    // Assert
    assert.Equal(t, 200, w.Code)
    mockService.AssertExpectations(t)
}
```

#### Integration Tests

```go
// tests/integration/api_test.go
package integration

import (
    "testing"
    "net/http/httptest"
)

func TestAPIFlow_LoginAndGetData(t *testing.T) {
    // Setup test server
    testServer := setupTestServer()
    defer testServer.Close()
    
    // Test login
    loginResp := login(testServer, "admin", "admin123")
    assert.NotEmpty(t, loginResp.Token)
    
    // Test authenticated request
    cpuResp := getCPUData(testServer, loginResp.Token)
    assert.NotNil(t, cpuResp.Usage)
}
```

### Frontend Testing

#### Component Testing

```typescript
// components/Monitor/CPU.spec.ts
import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import CPU from '@/views/Monitor/CPU.vue'

describe('CPU Monitor', () => {
  it('displays CPU usage percentage', () => {
    const wrapper = mount(CPU, {
      props: {
        cpuData: {
          usage: 0.45,
          cores: 24,
          load1: 2.5
        }
      }
    })
    
    expect(wrapper.text()).toContain('45.0%')
  })
  
  it('updates when cpuData prop changes', async () => {
    const wrapper = mount(CPU, {
      props: {
        cpuData: { usage: 0.45, cores: 24, load1: 2.5 }
      }
    })
    
    await wrapper.setProps({ 
      cpuData: { usage: 0.75, cores: 24, load1: 5.0 }
    })
    
    expect(wrapper.text()).toContain('75.0%')
  })
})
```

#### Store Testing

```typescript
// stores/auth.spec.ts
import { setActivePinia, createPinia } from 'pinia'
import { describe, it, expect, beforeEach } from 'vitest'
import { useAuthStore } from '@/stores/auth'

describe('Auth Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
  })
  
  it('stores token after login', () => {
    const authStore = useAuthStore()
    authStore.setToken('test-token')
    
    expect(authStore.token).toBe('test-token')
    expect(localStorage.getItem('access_token')).toBe('test-token')
  })
  
  it('clears token on logout', () => {
    const authStore = useAuthStore()
    authStore.setToken('test-token')
    authStore.logout()
    
    expect(authStore.token).toBeNull()
    expect(localStorage.getItem('access_token')).toBeNull()
  })
})
```

---

## Debugging Techniques

### Backend Debugging

#### Enable Debug Mode

```bash
# backend/config/dev.env
GIN_MODE=debug
LOG_LEVEL=debug
```

#### Debug Logging

```go
// Use structured logging
logger.Debug("Processing request",
    "method", c.Request.Method,
    "path", c.Request.URL.Path,
    "ip", c.ClientIP(),
)

// Log errors with context
logger.Error("Failed to connect to database",
    "error", err.Error(),
    "host", dbHost,
    "port", dbPort,
)
```

#### Debug Endpoints

```go
// Add debug endpoint (only in debug mode)
if gin.Mode() == gin.DebugMode {
    router.GET("/debug/vars", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "environment": config.Env,
            "database": config.DBHost,
            "version": config.Version,
        })
    })
}
```

### Frontend Debugging

#### Browser DevTools

```typescript
// Use browser console for debugging
console.log('API Request:', config)
console.table(data)
console.error('Error:', error)
```

#### Vue DevTools

1. Install Vue DevTools extension
2. Inspect component state
3. Monitor Pinia stores
4. Debug component hierarchy

#### Debug Utilities

```typescript
// utils/debug.ts
export const debugLog = (message: string, data?: any) => {
  if (import.meta.env.VITE_DEBUG === 'true') {
    console.log(`[DEBUG] ${message}`, data)
  }
}

// Usage
debugLog('API Request', { url: '/api/cpu', params: {} })
debugLog('Response Data', response.data)
```

---

## Performance Optimization

### Backend Optimization

#### Database Connection Pool

```go
// Optimize connection pool
db.DB.SetMaxOpenConns(100)
db.DB.SetMaxIdleConns(10)
db.DB.SetConnMaxLifetime(time.Hour)
```

#### Efficient Queries

```go
// Bad: N+1 query problem
users, _ := db.Find(&users).Error
for _, user := range users {
    db.Where("user_id = ?", user.ID).Find(&posts)
}

// Good: Use preloading
db.Preload("Posts").Find(&users)
```

#### Caching

```go
// Implement caching
type Cache struct {
    data map[string]interface{}
    mu   sync.RWMutex
    ttl  time.Duration
}

func (c *Cache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    value, exists := c.data[key]
    return value, exists
}
```

### Frontend Optimization

#### Code Splitting

```typescript
// Lazy load routes
const Dashboard = () => import('@/views/Dashboard.vue')
const StorageManager = () => import('@/apps/StorageManager.vue')

// Lazy load components
const HeavyComponent = defineAsyncComponent(
    () => import('@/components/HeavyComponent.vue')
)
```

#### Virtual Scrolling

```vue
<!-- Use virtual scrolling for large lists -->
<VirtualList
    :items="largeDataset"
    :item-height="50"
    :viewport-height="600"
/>
```

#### Debouncing

```typescript
import { debounce } from '@vueuse/core'

const search = debounce((query: string) => {
    performSearch(query)
}, 300)
```

---

## Contributing Guidelines

### Before Contributing

1. **Read documentation**
   - Project README
   - Architecture documentation
   - API documentation

2. **Set up development environment**
   - Follow setup instructions
   - Verify everything works
   - Run existing tests

3. **Choose an issue**
   - Look for "good first issue" label
   - Comment on issue you want to work on
   - Wait for assignment

### Making Contributions

#### Code Quality

1. **Follow coding standards**
   - Go standards for backend
   - Vue standards for frontend
   - TypeScript best practices

2. **Write tests**
   - Unit tests for new functions
   - Integration tests for APIs
   - Component tests for UI

3. **Update documentation**
   - Update README if needed
   - Add inline comments
   - Update API documentation

#### Pull Request Guidelines

1. **Descriptive title**
   ```
   feat(auth): add refresh token support
   ```

2. **Detailed description**
   - What changes were made
   - Why they were made
   - How they were tested

3. **Link issues**
   - References related issues
   - Closes #123

---

## Release Process

### Versioning

Follow [Semantic Versioning](https://semver.org/):

```
MAJOR.MINOR.PATCH

Example: 0.1.0
- MAJOR: Incompatible changes
- MINOR: New features, backwards compatible
- PATCH: Bug fixes, backwards compatible
```

### Release Checklist

1. **Update version numbers**
   ```bash
   # backend/go.mod
   module github.com/yourusername/nas-dashboard/v0.2.0
   
   # frontend/package.json
   "version": "0.2.0"
   ```

2. **Update changelog**
   ```markdown
   ## [0.2.0] - 2024-01-15
   
   ### Added
   - Refresh token support
   - Widget library
   
   ### Fixed
   - Network monitoring display
   - Token refresh logic
   
   ### Changed
   - Improved error handling
   - Updated dependencies
   ```

3. **Run full test suite**
   ```bash
   go test ./...
   npm run test
   ```

4. **Build production packages**
   ```bash
   # Backend
   cd backend
   go build -o nas-dashboard cmd/server/main.go
   
   # Frontend
   cd frontend
   npm run build
   ```

5. **Create Docker images**
   ```bash
   docker-compose build
   ```

6. **Tag release**
   ```bash
   git tag -a v0.2.0 -m "Release v0.2.0"
   git push origin v0.2.0
   ```

7. **Create GitHub release**
   - Upload binaries
   - Include changelog
   - Add release notes

---

## Conclusion

This developer guide provides comprehensive information for contributing to the NAS Dashboard project. Following these guidelines ensures high-quality contributions that benefit the entire community.

### Next Steps

1. **Set up development environment**
2. **Explore the codebase**
3. **Find an issue to work on**
4. **Make your first contribution**
5. **Join the community**

### Resources

- **Documentation**: `/docs` directory
- **Issues**: [GitHub Issues](https://github.com/yourusername/nas-dashboard/issues)
- **Discussions**: [GitHub Discussions](https://github.com/yourusername/nas-dashboard/discussions)
- **Code of Conduct**: See CODE_OF_CONDUCT.md

---

**Last Updated**: 2026-06-12  
**Version**: 0.1.0  
**Status**: Developer Guide Complete
