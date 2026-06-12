# NAS Dashboard - Development Guide

Complete guide for developers contributing to the NAS Dashboard project.

## 📋 Table of Contents

1. [Development Setup](#development-setup)
2. [Project Structure](#project-structure)
3. [Coding Standards](#coding-standards)
4. [Testing Guidelines](#testing-guidelines)
5. [Debug Mode Usage](#debug-mode-usage)
6. [Contributing Guidelines](#contributing-guidelines)
7. [Release Process](#release-process)

---

## Development Setup

### Prerequisites

- **Go**: 1.22+
- **Node.js**: 18.0+
- **Git**: For version control
- **Docker**: 20.10+ (for container management testing)
- **VS Code** (recommended) or your preferred IDE

### Initial Setup

#### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/nas-dashboard.git
cd nas-dashboard
```

#### 2. Backend Development Setup

```bash
cd backend

# Install Go dependencies
go mod download

# Install development tools
go install github.com/cosmtrek/air@latest  # Live reload tool
go install github.com/git Cosby/cosetty@latest  # API testing

# Verify installation
go version
air --version
```

#### 3. Frontend Development Setup

```bash
cd frontend

# Install dependencies
npm install

# Verify installation
npm --version
node --version
```

#### 4. Configure Development Environment

**Backend Config (`backend/config/dev.env`):**
```bash
GIN_MODE=debug
PORT=8888
JWT_SECRET=dev-secret-key-change-in-production
JWT_ACCESS_DURATION=24h
JWT_REFRESH_DURATION=720h
CORS_ORIGIN=http://localhost:5173
LOG_LEVEL=debug
```

**Frontend Config (`frontend/.env`):**
```bash
VITE_API_URL=http://localhost:8888
VITE_WS_URL=ws://localhost:8888
VITE_DEBUG=true
```

### Development Workflow

#### Start Backend with Live Reload

```bash
cd backend
air
```

**Air Configuration (`backend/.air.toml`):**
```toml
root = "."
tmp_dir = "tmp"

[build]
cmd = "go build -o ./tmp/main ."
bin = "tmp/main"
include_ext = ["go", "env"]
exclude_ext = ["tmp"]
exclude_dir = ["tmp", "vendor"]
delay = 1000
stop_on_error = true
```

#### Start Frontend Development Server

```bash
cd frontend
npm run dev
```

**Access the application at:** `http://localhost:5173`

---

## Project Structure

### Overall Structure

```
nas-dashboard/
├── backend/                 # Go backend application
│   ├── cmd/                # Application entry points
│   ├── internal/           # Private application code
│   ├── pkg/                # Public packages
│   ├── config/             # Configuration files
│   └── docs/               # Backend documentation
├── frontend/               # Vue 3 frontend application
│   ├── src/                # Source code
│   │   ├── api/           # API client
│   │   ├── assets/        # Static assets
│   │   ├── components/    # Vue components
│   │   ├── stores/       # Pinia stores
│   │   ├── views/        # Page components
│   │   └── utils/        # Utility functions
│   ├── public/            # Public static files
│   └── dist/              # Build output
├── docs/                   # Project documentation
├── docker-compose.yml      # Docker orchestration
└── README.md              # Project overview
```

### Backend Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go        # Application entry point
├── internal/
│   ├── api/               # API handlers
│   │   ├── auth.go        # Authentication endpoints
│   │   ├── docker.go      # Docker management
│   │   ├── monitor.go     # System monitoring
│   │   ├── storage.go     # Storage management
│   │   ├── user.go        # User management
│   │   └── service.go     # Service management
│   ├── middleware/        # HTTP middleware
│   │   ├── auth.go        # JWT authentication
│   │   ├── cors.go        # CORS handling
│   │   └── logger.go      # Request logging
│   ├── models/            # Data models
│   │   ├── user.go        # User model
│   │   ├── container.go   # Container model
│   │   └── system.go      # System model
│   ├── services/          # Business logic
│   │   ├── docker.go      # Docker service
│   │   ├── monitor.go     # Monitoring service
│   │   └── storage.go     # Storage service
│   └── websocket/         # WebSocket handlers
│       └── monitor.go     # Monitoring WebSocket
├── pkg/                   # Public packages
│   ├── jwt/               # JWT utilities
│   └── response/          # API response helpers
├── config/                # Configuration
│   ├── dev.env           # Development config
│   └── prod.env          # Production config
├── go.mod                # Go dependencies
└── go.sum                # Dependency checksums
```

### Frontend Structure

```
frontend/
├── src/
│   ├── api/              # API client
│   │   ├── client.ts     # Axios configuration
│   │   ├── auth.ts       # Authentication API
│   │   ├── monitor.ts    # Monitoring API
│   │   ├── docker.ts     # Docker API
│   │   └── user.ts       # User API
│   ├── assets/           # Static assets
│   │   ├── css/          # Global styles
│   │   └── images/       # Images and icons
│   ├── components/       # Vue components
│   │   ├── Layout/       # Layout components
│   │   │   ├── Sidebar.vue
│   │   │   ├── Header.vue
│   │   │   └── Main.vue
│   │   ├── Monitor/      # Monitoring components
│   │   ├── Storage/      # Storage components
│   │   └── Docker/       # Docker components
│   ├── router/           # Vue Router
│   │   └── index.ts      # Route configuration
│   ├── stores/           # Pinia stores
│   │   ├── auth.ts       # Authentication store
│   │   ├── monitor.ts    # Monitoring store
│   │   └── user.ts       # User store
│   ├── utils/            # Utility functions
│   │   ├── format.ts     # Data formatting
│   │   ├── validation.ts # Input validation
│   │   └── websocket.ts  # WebSocket client
│   ├── views/            # Page components
│   │   ├── Login.vue     # Login page
│   │   ├── Dashboard.vue # Dashboard
│   │   ├── Monitor/      # Monitoring pages
│   │   │   ├── CPU.vue
│   │   │   ├── Memory.vue
│   │   │   ├── Disk.vue
│   │   │   └── Network.vue
│   │   ├── Storage.vue   # Storage page
│   │   ├── Docker.vue    # Docker page
│   │   └── Users.vue     # User management
│   ├── App.vue           # Root component
│   └── main.ts           # Application entry
├── public/               # Public static files
├── package.json          # NPM dependencies
├── vite.config.ts        # Vite configuration
└── tsconfig.json         # TypeScript configuration
```

---

## Coding Standards

### Go Backend Standards

#### Code Style

Follow standard Go conventions:

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

#### Configuration Management

```go
// Use environment variables for configuration
type Config struct {
    Port            string
    GinMode         string
    JWTSecret       string
    AccessDuration  time.Duration
    RefreshDuration time.Duration
    CORSOrigin      string
}

func LoadConfig() *Config {
    return &Config{
        Port:           getEnv("PORT", "8888"),
        GinMode:        getEnv("GIN_MODE", "debug"),
        JWTSecret:      getEnv("JWT_SECRET", "secret"),
        AccessDuration: parseDuration(getEnv("JWT_ACCESS_DURATION", "24h")),
        // ...
    }
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

#### State Management with Pinia

```typescript
// stores/monitor.ts
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useMonitorStore = defineStore('monitor', () => {
  // State
  const cpuData = ref(null)
  const memoryData = ref(null)
  const isConnected = ref(false)

  // Actions
  async function fetchCPUData() {
    const response = await api.getCPUData()
    cpuData.value = response.data
  }

  function setConnected(status: boolean) {
    isConnected.value = status
  }

  // Getters
  const cpuUsage = computed(() => {
    return cpuData.value?.usage || 0
  })

  return {
    cpuData,
    memoryData,
    isConnected,
    fetchCPUData,
    setConnected,
    cpuUsage
  }
})
```

#### API Client Organization

```typescript
// api/client.ts - Base configuration
import axios from 'axios'

const client = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8888',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Request interceptor
client.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('access_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

// Response interceptor
client.interceptors.response.use(
  (response) => response,
  async (error) => {
    if (error.response?.status === 401) {
      // Token refresh logic
      await refreshAccessToken()
      return client.request(error.config)
    }
    return Promise.reject(error)
  }
)

export default client

// api/monitor.ts - Specific API endpoints
import client from './client'

export const monitorApi = {
  getCPUData: () => client.get('/api/monitor/cpu'),
  getMemoryData: () => client.get('/api/monitor/memory'),
  // ...
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
    "testing"
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

#### Component Testing with Vitest

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

## Debug Mode Usage

### Backend Debug Mode

**Enable in `config/dev.env`:**
```bash
GIN_MODE=debug
LOG_LEVEL=debug
```

**Debug features:**
- Detailed request logging
- Stack traces on errors
- Pretty-printed JSON responses
- Development-friendly error messages

**View debug logs:**
```bash
# Real-time logs
tail -f /var/log/nas-dashboard/debug.log

# Filter by level
grep "ERROR" /var/log/nas-dashboard/debug.log
grep "WARN" /var/log/nas-dashboard/debug.log
```

### Frontend Debug Mode

**Enable in `.env`:**
```bash
VITE_DEBUG=true
```

**Debug features:**
- Detailed API call logging
- Request/response body inspection
- Error stack traces
- Performance timing

**Use browser console:**
```javascript
// API calls are logged
// Request details, timing, and errors shown
```

**Debug utilities:**
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

## Contributing Guidelines

### Workflow

1. **Fork the repository**
   ```bash
   # On GitHub, fork the repository to your account
   ```

2. **Clone your fork**
   ```bash
   git clone https://github.com/yourusername/nas-dashboard.git
   cd nas-dashboard
   ```

3. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

4. **Make changes and test**
   ```bash
   # Make your changes
   # Test thoroughly
   npm run test
   go test ./...
   ```

5. **Commit changes**
   ```bash
   git add .
   git commit -m "feat: add user profile page"
   ```

6. **Push to your fork**
   ```bash
   git push origin feature/your-feature-name
   ```

7. **Create pull request**
   - Describe your changes
   - Reference related issues
   - Ensure tests pass

### Commit Message Format

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
- `style`: Code style changes (formatting)
- `refactor`: Code refactoring
- `test`: Test additions or changes
- `chore`: Build process or auxiliary tool changes

**Examples:**
```bash
feat(auth): add refresh token support
fix(dashboard): correct network data format
docs(api): update authentication documentation
refactor(monitor): optimize CPU usage calculation
test(docker): add container management tests
```

### Code Review Process

1. **Submit pull request**
   - Clear description of changes
   - Related issues linked
   - Tests included

2. **Automated checks**
   - CI/CD pipeline runs
   - Code quality checks
   - Test coverage verification

3. **Manual review**
   - Code style compliance
   - Functionality verification
   - Security considerations

4. **Approval and merge**
   - At least one approval required
   - No unresolved comments
   - All checks pass

### Issue Reporting

**Bug Report Template:**
```markdown
## Description
Clear description of the bug

## Steps to Reproduce
1. Go to...
2. Click on...
3. See error...

## Expected Behavior
What should happen

## Actual Behavior
What actually happens

## Environment
- OS: Ubuntu 22.04
- Browser: Chrome 90
- Version: 0.1.0

## Screenshots
If applicable, add screenshots

## Logs
Error logs or console output
```

**Feature Request Template:**
```markdown
## Description
Feature description

## Use Case
Why is this feature needed?

## Proposed Solution
How should it work?

## Alternatives
What other approaches were considered?

## Additional Context
Any additional information
```

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
   - Alert system
   - File manager
   
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

8. **Deploy to production**
   ```bash
   docker-compose pull
   docker-compose up -d
   ```

---

## Additional Resources

### Documentation

- [API Documentation](./API_DOCUMENTATION.md)
- [Installation Guide](./INSTALLATION.md)
- [Security Considerations](./SECURITY_CONSIDERATIONS.md)

### External Resources

- [Go Documentation](https://golang.org/doc/)
- [Vue 3 Documentation](https://vuejs.org/guide/)
- [Gin Framework](https://gin-gonic.com/docs/)
- [Pinia Documentation](https://pinia.vuejs.org/)

### Tools

- **Air**: Live reload for Go
- **Vitest**: Unit testing for Vue
- **ESLint**: JavaScript/TypeScript linting
- **GoLint**: Go code linting

---

**Last Updated**: 2026-06-12
**Version**: 0.1.0-alpha
