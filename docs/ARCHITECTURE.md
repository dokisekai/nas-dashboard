# NAS Dashboard - System Architecture

Complete architectural documentation for the NAS Dashboard system, including design decisions, component interactions, and technology choices.

## 📋 Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [System Architecture](#system-architecture)
3. [Component Architecture](#component-architecture)
4. [Data Flow Architecture](#data-flow-architecture)
5. [Security Architecture](#security-architecture)
6. [Deployment Architecture](#deployment-architecture)
7. [Technology Choices](#technology-choices)
8. [Design Patterns](#design-patterns)
9. [Performance Considerations](#performance-considerations)
10. [Scalability Architecture](#scalability-architecture)

---

## Architecture Overview

### High-Level Architecture

The NAS Dashboard follows a **modern web application architecture** with clear separation of concerns:

```
┌─────────────────────────────────────────────────────────────┐
│                        Presentation Layer                     │
│                    (Vue 3 + TypeScript)                       │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   Desktop    │  │  Widget      │  │   Window     │      │
│  │   System     │  │  System      │  │  Manager     │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└───────────────────────────┬─────────────────────────────────┘
                            │
                            │ REST API + WebSocket
                            │
┌───────────────────────────┴─────────────────────────────────┐
│                      Application Layer                        │
│                      (Go + Gin Framework)                     │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   API        │  │  WebSocket  │  │  Middleware  │      │
│  │  Handlers    │  │   Server    │  │    Layer     │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   Plugin     │  │   Market    │  │   Service    │      │
│  │   Loader     │  │   Place     │  │   Manager    │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└───────────────────────────┬─────────────────────────────────┘
                            │
                            │ System Calls + Docker API
                            │
┌───────────────────────────┴─────────────────────────────────┐
│                      Data Access Layer                       │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   System     │  │   Docker     │  │  Database    │      │
│  │  Monitoring  │  │   Client     │  │   (GORM)     │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└───────────────────────────┬─────────────────────────────────┘
                            │
                            │
┌───────────────────────────┴─────────────────────────────────┐
│                    Infrastructure Layer                      │
│           (Linux System, Docker Daemon, PostgreSQL)          │
└──────────────────────────────────────────────────────────────┘
```

### Architectural Principles

1. **Separation of Concerns**: Clear boundaries between layers
2. **Modularity**: Independent, reusable components
3. **Scalability**: Horizontal and vertical scaling support
4. **Security**: Defense-in-depth approach
5. **Performance**: Optimized for speed and efficiency
6. **Maintainability**: Clean code and comprehensive documentation
7. **Extensibility**: Plugin system for custom functionality

---

## System Architecture

### Three-Tier Architecture

#### 1. Presentation Tier (Frontend)

**Technologies**: Vue 3, TypeScript, Vite, Tailwind CSS

**Responsibilities**:
- User interface rendering
- Client-side state management
- User input handling
- API communication
- Real-time data visualization

**Key Components**:
- **Desktop System**: Window management, widgets, dock bar
- **Application Center**: Storage manager, system monitor, user manager
- **Plugin System**: Dynamic plugin loading and execution
- **State Management**: Pinia stores for global state
- **API Client**: Axios-based HTTP client
- **WebSocket Client**: Real-time data streaming

#### 2. Application Tier (Backend)

**Technologies**: Go, Gin Framework, GORM, JWT

**Responsibilities**:
- Business logic execution
- Request processing and validation
- Authentication and authorization
- Data transformation
- External service integration
- WebSocket connection management

**Key Components**:
- **API Handlers**: REST endpoints
- **WebSocket Server**: Real-time communication
- **Middleware**: Auth, CORS, logging
- **Service Layer**: Business logic
- **Plugin System**: Dynamic plugin loading
- **Database Layer**: Data persistence

#### 3. Data Tier

**Technologies**: PostgreSQL, Linux System APIs, Docker API

**Responsibilities**:
- Data persistence
- System information collection
- Container management
- File system operations

**Key Components**:
- **Database**: User data, configuration, audit logs
- **System APIs**: CPU, memory, disk, network monitoring
- **Docker API**: Container management
- **File System**: Storage operations

---

## Component Architecture

### Frontend Components

#### 1. Desktop System

```typescript
DSMDesktop.vue (Root Desktop Component)
├── EnhancedDock.vue (Dock Bar)
│   ├── Application Menu
│   ├── Running Applications
│   └── System Tray
├── WindowManager.vue (Window Manager)
│   ├── WindowSnap.vue (Window Snapping)
│   ├── DesktopWindow.vue (Individual Windows)
│   └── Tab Management
├── WidgetSystem (Widget Management)
│   ├── WidgetLibrary.vue (Widget Browser)
│   ├── WidgetConfig.vue (Widget Configuration)
│   └── Individual Widgets
│       ├── SystemMonitorWidget.vue
│       ├── WeatherWidget.vue
│       ├── CalendarWidget.vue
│       ├── QuickNoteWidget.vue
│       └── StorageWidget.vue
└── ThemeManager.vue (Theme Management)
    ├── Background Management
    ├── Color Schemes
    └── Display Settings
```

**Architecture Patterns**:
- **Component Composition**: Reusable component hierarchy
- **Event Bus**: Cross-component communication
- **State Management**: Pinia stores for global state
- **Props Drilling**: Controlled data flow

#### 2. Application Center

```typescript
ApplicationCenter.vue
├── StorageManager.vue
│   ├── Disk Management
│   ├── SMB Shares
│   └── File Browser
├── SystemMonitor.vue
│   ├── CPU Monitoring
│   ├── Memory Monitoring
│   ├── Disk Monitoring
│   └── Network Monitoring
├── UserManager.vue
│   ├── User CRUD
│   ├── SSH Keys
│   └── Group Management
├── SystemSettings.vue
│   ├── Network Configuration
│   ├── Service Management
│   └── Backup Management
├── AppCenter.vue
│   ├── Application Browser
│   ├── Installation Management
│   └── Update System
└── PluginStore.vue
    ├── Plugin Browser
    ├── Installation
    └── Configuration
```

#### 3. Plugin System

```typescript
PluginSystem/
├── core/PluginLoader.ts (Dynamic Loading)
│   ├── Plugin Loading
│   ├── Lifecycle Management
│   └── Dependency Resolution
├── sdk/ (Plugin Development Kit)
│   ├── Context Factory
│   ├── API Wrapper
│   ├── Storage System
│   ├── Logger
│   └── Utilities
├── manager/PluginManager.ts
│   ├── Installation
│   ├── Configuration
│   └── Permissions
└── marketplace/PluginMarketplace.ts
    ├── Discovery
    ├── Search
    └── Installation
```

### Backend Components

#### 1. API Layer

```go
internal/api/
├── auth.go (Authentication)
│   ├── Login
│   ├── Logout
│   ├── Token Refresh
│   └── User Info
├── monitor.go (System Monitoring)
│   ├── CPU Stats
│   ├── Memory Stats
│   ├── Disk Stats
│   └── Network Stats
├── storage.go (Storage Management)
│   ├── Disk Operations
│   ├── SMB Management
│   └── File Operations
├── docker.go (Container Management)
│   ├── Container List
│   ├── Container Control
│   └── Container Logs
├── user.go (User Management)
│   ├── User CRUD
│   ├── Password Management
│   └── SSH Keys
└── service.go (Service Management)
    ├── Service Control
    └── Service Status
```

#### 2. Service Layer

```go
internal/services/
├── auth_service.go (Authentication Logic)
├── monitor_service.go (Monitoring Logic)
├── storage_service.go (Storage Logic)
├── docker_service.go (Docker Logic)
├── user_service.go (User Management Logic)
└── plugin_service.go (Plugin Logic)
```

#### 3. WebSocket Layer

```go
internal/websocket/
├── hub.go (Connection Management)
│   ├── Connection Pool
│   ├── Message Broadcasting
│   └── Client Registration
└── manager.go (WebSocket Manager)
    ├── Connection Handler
    ├── Message Handler
    └── Reconnection Logic
```

---

## Data Flow Architecture

### Request/Response Flow

#### REST API Request Flow

```
┌─────────────┐
│   Browser   │
└──────┬──────┘
       │ HTTP Request
       │ (JWT Token)
       ▼
┌─────────────┐
│   Router    │
│   (Gin)     │
└──────┬──────┘
       │ Route
       ▼
┌─────────────┐
│ Middleware  │
│  - CORS     │
│  - Auth     │
│  - Logger   │
└──────┬──────┘
       │ Validated
       ▼
┌─────────────┐
│ API Handler │
└──────┬──────┘
       │ Call Service
       ▼
┌─────────────┐
│  Service    │
│   Layer     │
└──────┬──────┘
       │ Query Data
       ▼
┌─────────────┐
│ Data Access │
│   Layer     │
└──────┬──────┘
       │ System Call
       ▼
┌─────────────┐
│   System    │
│  / Docker   │
└─────────────┘
```

#### WebSocket Data Flow

```
┌─────────────┐
│   Browser   │
└──────┬──────┘
       │ WebSocket Connect
       │ (JWT Token)
       ▼
┌─────────────┐
│   WebSocket │
│   Manager   │
└──────┬──────┘
       │ Authenticate
       ▼
┌─────────────┐
│     Hub     │
│  (Pool)     │
└──────┬──────┘
       │ Register
       ▼
┌─────────────┐
│  Monitor    │
│  Service    │
└──────┬──────┘
       │ Push Data
       ▼
┌─────────────┐
│ Broadcast   │
│   to Hub    │
└──────┬──────┘
       │ Send
       ▼
┌─────────────┐
│   Clients   │
└─────────────┘
```

### State Management Flow

#### Frontend State Flow (Pinia)

```
┌─────────────┐
│   Component │
└──────┬──────┘
       │ Action
       ▼
┌─────────────┐
│   Store     │
│  (Pinia)    │
└──────┬──────┘
       │ API Call
       ▼
┌─────────────┐
│ API Client  │
│  (Axios)    │
└──────┬──────┘
       │ HTTP Request
       ▼
┌─────────────┐
│   Backend   │
└──────┬──────┘
       │ Response
       ▼
┌─────────────┐
│ Update State│
└──────┬──────┘
       │ Reactivity
       ▼
┌─────────────┐
│   Component │
│   Update    │
└─────────────┘
```

---

## Security Architecture

### Authentication Flow

```
┌─────────────┐
│   User      │
└──────┬──────┘
       │ Login Request
       ▼
┌─────────────┐
│   Frontend  │
└──────┬──────┘
       │ POST /api/auth/login
       ▼
┌─────────────┐
│   Backend   │
└──────┬──────┘
       │ Validate Credentials
       ▼
┌─────────────┐
│  Database   │
└──────┬──────┘
       │ User Data
       ▼
┌─────────────┐
│  bcrypt     │
│  Verify     │
└──────┬──────┘
       │ Generate JWT
       ▼
┌─────────────┐
│ Return Token│
└─────────────┘
```

### Authorization Architecture

```
┌─────────────┐
│   Request   │
└──────┬──────┘
       │ JWT Token
       ▼
┌─────────────┐
│  Middleware │
│  - Validate │
│  - Parse    │
└──────┬──────┘
       │ User Info
       ▼
┌─────────────┐
│ Permission  │
│   Check     │
└──────┬──────┘
       │ Allowed?
       ▼
┌─────────────┐
│ API Handler │
└─────────────┘
```

### Security Layers

1. **Network Security**
   - TLS/SSL encryption
   - CORS protection
   - Rate limiting

2. **Authentication Security**
   - JWT token validation
   - bcrypt password hashing
   - Session management

3. **Authorization Security**
   - Role-based access control
   - Permission checks
   - Audit logging

4. **Data Security**
   - Input validation
   - SQL injection prevention
   - XSS protection

---

## Deployment Architecture

### Development Deployment

```
┌────────────────────────────────────────┐
│         Developer Machine               │
│  ┌──────────────┐  ┌──────────────┐   │
│  │  Frontend    │  │   Backend    │   │
│  │  (Vite Dev)  │  │  (Go Run)    │   │
│  │  Port: 5173  │  │  Port: 8888  │   │
│  └──────────────┘  └──────────────┘   │
│         │                 │           │
│         └─────────────────┘           │
│                   │                    │
│              Local DB                 │
└────────────────────────────────────────┘
```

### Production Deployment

```
┌─────────────────────────────────────────────────────────┐
│                    Production Server                     │
│  ┌─────────────────────────────────────────────────┐  │
│  │              Reverse Proxy (Nginx)                │  │
│  │                    SSL/TLS                       │  │
│  └────────────┬──────────────────────┬─────────────┘  │
│               │                      │                │
│  ┌────────────┴─────┐    ┌───────────┴──────────────┐ │
│  │  Frontend Server  │    │   Backend Server        │ │
│  │  (Static Files)   │    │   (Go Application)      │ │
│  └───────────────────┘    └──────────┬──────────────┘ │
│                                      │                  │
│  ┌───────────────────────────────────┼──────────────┐ │
│  │         PostgreSQL Database        │              │ │
│  └───────────────────────────────────┴──────────────┘ │
│  ┌───────────────────────────────────────────────────┐│
│  │            Docker Daemon                            ││
│  └───────────────────────────────────────────────────┘│
└─────────────────────────────────────────────────────────┘
```

### Docker Deployment

```
┌─────────────────────────────────────────────────────────┐
│              Docker Compose Environment                 │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐ │
│  │   nginx      │  │  frontend    │  │   backend    │ │
│  │   container  │  │  container   │  │  container   │ │
│  │   Port: 80   │  │  Port: 3000  │  │  Port: 8888  │ │
│  └──────────────┘  └──────────────┘  └──────────────┘ │
│         │                    │                    │       │
│         └────────────────────┼────────────────────┘       │
│                              │                            │
│  ┌───────────────────────────┼──────────────────────────┐│
│  │      postgres container    │                          ││
│  └───────────────────────────┴──────────────────────────┘│
└──────────────────────────────────────────────────────────┘
```

---

## Technology Choices

### Frontend Technology Choices

#### Vue 3 + Composition API
**Why?**
- Reactive and composable
- Excellent TypeScript support
- Smaller bundle size
- Better performance than Vue 2
- Large ecosystem and community

**Benefits**:
- Easy state management with reactive refs
- Reusable composition functions
- Type safety with TypeScript
- Fast rendering with virtual DOM

#### TypeScript
**Why?**
- Type safety prevents bugs
- Better IDE support
- Self-documenting code
- Easier refactoring
- Large projects maintainability

#### Vite
**Why?**
- Instant server start
- Lightning-fast HMR
- Optimized build
- Native ES modules
- Great developer experience

#### Tailwind CSS
**Why?**
- Utility-first approach
- No custom CSS to maintain
- Consistent design system
- Easy customization
- Small production bundle

### Backend Technology Choices

#### Go
**Why?**
- High performance
- Strong concurrency
- Simple deployment (single binary)
- Strong type system
- Excellent standard library

#### Gin Framework
**Why?**
- Fast HTTP router
- Middleware support
- JSON validation
- Group routing
- Good documentation

#### GORM + PostgreSQL
**Why?**
- Type-safe ORM
- Migration support
- Relationship handling
- ACID compliance
- Complex queries

#### JWT Authentication
**Why?**
- Stateless authentication
- Cross-platform support
- Built-in expiration
- Easy to implement
- Industry standard

---

## Design Patterns

### 1. Repository Pattern

**Data access abstraction**:

```go
type UserRepository interface {
    Create(user *User) error
    FindByID(id int) (*User, error)
    Update(user *User) error
    Delete(id int) error
}

type userRepository struct {
    db *gorm.DB
}

func (r *userRepository) Create(user *User) error {
    return r.db.Create(user).Error
}
```

### 2. Factory Pattern

**Plugin creation**:

```typescript
export function createPlugin(context: PluginContext): Plugin {
    return {
        manifest,
        lifecycle: {
            onInstall: async () => { /* ... */ },
            onActivate: async () => { /* ... */ }
        }
    };
}
```

### 3. Observer Pattern

**WebSocket data streaming**:

```go
type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.clients[client] = true
        case message := <-h.broadcast:
            for client := range h.clients {
                client.send <- message
            }
        }
    }
}
```

### 4. Strategy Pattern

**Monitoring strategies**:

```typescript
interface MonitoringStrategy {
    collect(): Promise<MonitoringData>;
    transform(data: MonitoringData): DisplayData;
    display(data: DisplayData): void;
}

class CPUMonitoringStrategy implements MonitoringStrategy {
    async collect(): Promise<MonitoringData> {
        return await api.getCPUData();
    }
    // ...
}
```

### 5. Singleton Pattern

**Global state management**:

```typescript
export const useAuthStore = defineStore('auth', () => {
    // Single instance across app
    const token = ref<string | null>(null);
    const user = ref<User | null>(null);
    
    return { token, user };
});
```

---

## Performance Considerations

### Frontend Performance

#### 1. Code Splitting
```typescript
// Lazy load routes
const Dashboard = () => import('@/views/Dashboard.vue');
const StorageManager = () => import('@/apps/StorageManager.vue');

// Lazy load components
const HeavyComponent = defineAsyncComponent(
    () => import('@/components/HeavyComponent.vue')
);
```

#### 2. Virtual Scrolling
```vue
<VirtualList
    :items="largeDataset"
    :item-height="50"
    :viewport-height="600"
/>
```

#### 3. Debouncing
```typescript
import { debounce } from '@vueuse/core';

const search = debounce((query: string) => {
    performSearch(query);
}, 300);
```

#### 4. Memoization
```typescript
const expensiveValue = computed(() => {
    return heavyCalculation(source.value);
});
```

### Backend Performance

#### 1. Connection Pooling
```go
db.DB.SetMaxOpenConns(100)
db.DB.SetMaxIdleConns(10)
db.DB.SetConnMaxLifetime(time.Hour)
```

#### 2. Caching
```go
// In-memory cache for frequently accessed data
type Cache struct {
    data map[string]interface{}
    mu   sync.RWMutex
}

func (c *Cache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    value, exists := c.data[key]
    return value, exists
}
```

#### 3. Goroutines
```go
// Parallel system monitoring
func GetSystemInfo() *SystemInfo {
    var wg sync.WaitGroup
    var cpuInfo CPUInfo
    var memInfo MemoryInfo
    
    wg.Add(2)
    
    go func() {
        cpuInfo = getCPUInfo()
        wg.Done()
    }()
    
    go func() {
        memInfo = getMemoryInfo()
        wg.Done()
    }()
    
    wg.Wait()
    return &SystemInfo{CPU: cpuInfo, Memory: memInfo}
}
```

---

## Scalability Architecture

### Horizontal Scaling

```
                    ┌──────────────┐
                    │   Load       │
                    │  Balancer    │
                    └──────┬───────┘
                           │
        ┌──────────────────┼──────────────────┐
        │                  │                  │
┌───────┴───────┐  ┌───────┴───────┐  ┌───────┴───────┐
│   Backend     │  │   Backend     │  │   Backend     │
│   Instance 1  │  │   Instance 2  │  │   Instance 3  │
└───────────────┘  └───────────────┘  └───────────────┘
        │                  │                  │
        └──────────────────┼──────────────────┘
                           │
                  ┌────────┴────────┐
                  │  Shared         │
                  │  Database       │
                  │  (PostgreSQL)   │
                  └─────────────────┘
```

### Vertical Scaling

- **Database Optimization**: Indexing, query optimization
- **Caching Layer**: Redis for session and data caching
- **CDN**: Static asset delivery
- **Database Replication**: Read replicas for scaling reads

### Microservices Architecture (Future)

```
┌─────────────────────────────────────────────────────────┐
│                  API Gateway                             │
└─────────────────────────────────────────────────────────┘
         │              │              │              │
┌────────┴────┐   ┌──────┴──────┐   ┌──┴──────────┐   ┌──┴─────┐
│  Monitor    │   │   Storage   │   │    Docker   │   │  User   │
│  Service    │   │   Service   │   │   Service   │   │ Service │
└─────────────┘   └─────────────┘   └─────────────┘   └────────┘
```

---

## Monitoring and Observability

### Application Monitoring

```
┌─────────────────────────────────────────────────────────┐
│              Monitoring Stack                            │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐ │
│  │   Metrics    │  │   Logging    │  │  Tracing     │ │
│  │  (Prometheus)│  │  (ELK Stack) │  │  (Jaeger)    │ │
│  └──────────────┘  └──────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────┘
```

### Health Checks

```go
func HealthCheck() gin.HandlerFunc {
    return func(c *gin.Context) {
        status := map[string]string{
            "status": "healthy",
            "timestamp": time.Now().Format(time.RFC3339),
            "version": "0.1.0",
        }
        
        // Check database
        if err := db.DB.Ping(); err != nil {
            status["status"] = "unhealthy"
            status["database"] = "down"
        }
        
        c.JSON(200, status)
    }
}
```

---

## Future Architecture Improvements

### Planned Enhancements

1. **GraphQL API**: Alternative to REST for complex queries
2. **Message Queue**: RabbitMQ/Kafka for async processing
3. **Microservices**: Split monolith into services
4. **Caching Layer**: Redis for performance
5. **CDN Integration**: Global static asset delivery
6. **Real-time Collaboration**: WebRTC for collaborative features
7. **Advanced Monitoring**: Distributed tracing and metrics

### Architecture Evolution Path

```
Current → Enhanced → Distributed
Monolith   Modular   Microservices
```

---

## Conclusion

The NAS Dashboard architecture is designed with:

- **Modularity**: Easy to maintain and extend
- **Performance**: Optimized for speed and efficiency
- **Security**: Multiple layers of protection
- **Scalability**: Ready to grow with demand
- **Maintainability**: Clean code and documentation

The architecture supports current requirements while being flexible enough for future enhancements and scaling needs.

---

**Last Updated**: 2026-06-12  
**Version**: 0.1.0  
**Status**: Active Development
