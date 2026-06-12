# NAS Dashboard - Project Overview

## 📋 Project Description

The NAS Dashboard is a comprehensive web-based monitoring and management system designed for Network Attached Storage (NAS) devices and Linux servers. It provides real-time system monitoring, Docker container management, storage management, and user administration through an intuitive web interface.

### 🎯 Project Goals

- **Real-time Monitoring**: Provide live system metrics (CPU, memory, disk, network) through WebSocket connections
- **Container Management**: Easy Docker container management with start, stop, and restart capabilities
- **Storage Management**: Monitor disk usage and manage storage pools
- **User Administration**: Simple user management with role-based access control
- **Responsive Design**: Modern, mobile-friendly interface accessible from any device
- **Secure Access**: JWT-based authentication with refresh token support

## 🏗️ Architecture Overview

### System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                         Browser Client                        │
│                    (Vue 3 + TypeScript)                      │
└───────────────────────────┬─────────────────────────────────┘
                            │
                            │ HTTPS/WSS
                            │
┌───────────────────────────┴─────────────────────────────────┐
│                      Backend Server                           │
│                   (Go + Gin Framework)                        │
│                                                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐       │
│  │    API       │  │   WebSocket  │  │  Middleware  │       │
│  │  Handlers    │  │   Server     │  │   Layer      │       │
│  └──────────────┘  └──────────────┘  └──────────────┘       │
│                                                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐       │
│  │    System    │  │    Docker    │  │   Storage    │       │
│  │  Monitoring  │  │   Manager    │  │   Manager    │       │
│  └──────────────┘  └──────────────┘  └──────────────┘       │
└───────────────────────────┬─────────────────────────────────┘
                            │
                            │ System Calls
                            │
┌───────────────────────────┴─────────────────────────────────┐
│                    Linux System                               │
│              (Docker, Filesystem, Processes)                 │
└──────────────────────────────────────────────────────────────┘
```

### Component Architecture

#### Frontend Components

- **Authentication Layer**: JWT-based authentication with automatic token refresh
- **State Management**: Pinia stores for auth, monitor data, and user state
- **API Client**: Axios-based client with interceptors for token management
- **WebSocket Client**: Real-time data streaming with automatic reconnection
- **UI Components**: Reusable Vue components with TailwindCSS styling

#### Backend Components

- **API Layer**: RESTful endpoints for all operations
- **WebSocket Server**: Real-time monitoring data streaming
- **Middleware Layer**: Authentication, CORS, and logging
- **Service Layer**: Business logic for system operations
- **Data Access Layer**: System monitoring and Docker integration

## 🔧 Technology Stack

### Frontend

| Technology | Version | Purpose |
|------------|---------|---------|
| **Vue.js** | 3.4+ | Progressive JavaScript framework |
| **TypeScript** | 5.3+ | Type-safe JavaScript |
| **Vite** | 5.0+ | Build tool and development server |
| **Pinia** | 2.1+ | State management |
| **Vue Router** | 4.2+ | Client-side routing |
| **Axios** | 1.6+ | HTTP client |
| **TailwindCSS** | 3.4+ | Utility-first CSS framework |
| **Chart.js** | 4.4+ | Data visualization |
| **SockJS/Stomp** | - | WebSocket client |

### Backend

| Technology | Version | Purpose |
|------------|---------|---------|
| **Go** | 1.22+ | Programming language |
| **Gin** | 1.9+ | Web framework |
| **gopsutil** | 4.23+ | System monitoring |
| **Docker SDK** | 24.0+ | Container management |
| **JWT-Go** | 5.2+ | Token authentication |
| **CORS** | 1.4+ | Cross-origin resource sharing |
| **WebSocket** | - | Real-time communication |

### System Requirements

- **Operating System**: Linux (Ubuntu 20.04+, Debian 11+, CentOS 8+)
- **Go Runtime**: 1.22 or higher
- **Node.js**: 18.0 or higher
- **Docker**: 20.10 or higher (for container management)
- **Memory**: Minimum 2GB RAM (4GB recommended)
- **Storage**: 100MB free disk space for application

## 📊 Feature Matrix

### ✅ Fully Implemented Features

| Feature | Description | Status |
|---------|-------------|--------|
| **Authentication** | JWT-based login/logout with refresh tokens | ✅ Complete |
| **Dashboard** | System overview with statistics cards | ✅ Complete |
| **CPU Monitoring** | Real-time CPU usage and load metrics | ✅ Complete |
| **Memory Monitoring** | RAM and swap usage tracking | ✅ Complete |
| **Disk Monitoring** | Multiple disk usage and I/O metrics | ✅ Complete |
| **Network Monitoring** | Interface statistics and traffic data | ✅ Complete |
| **Storage Management** | Disk usage overview and management | ✅ Complete |
| **Docker Management** | Container listing, start, stop, restart | ✅ Complete |
| **User Management** | User CRUD operations with UI | ✅ Complete |
| **WebSocket Monitoring** | Real-time system data streaming | ✅ Complete |
| **Responsive Design** | Mobile-friendly interface | ✅ Complete |
| **Dark Mode** | Dark theme support | ✅ Complete |

### 🚧 Partially Implemented Features

| Feature | Description | Status | Notes |
|---------|-------------|--------|-------|
| **Service Management** | Systemd service management | ⚠️ Partial | Backend complete, frontend needs UI |
| **Group Management** | User group management | ⚠️ Partial | API exists, no UI |
| **Alert System** | System alerts and notifications | ❌ Missing | Backend and frontend needed |
| **File Manager** | File browsing and management | ❌ Missing | Complete implementation needed |
| **Log Viewer** | System log viewing | ❌ Missing | Complete implementation needed |
| **Backup System** | Backup and restore | ❌ Missing | Complete implementation needed |
| **Historical Data** | Time-series monitoring data | ❌ Missing | Database integration needed |

### 🎯 Planned Features

| Feature | Priority | Est. Effort |
|---------|----------|-------------|
| **Alert System** | P0 | 16 hours |
| **File Manager** | P1 | 24 hours |
| **Log Viewer** | P1 | 16 hours |
| **Historical Monitoring** | P1 | 20 hours |
| **Backup Management** | P1 | 16 hours |
| **Docker Compose Support** | P2 | 12 hours |
| **System Updates** | P2 | 8 hours |
| **Terminal Access** | P3 | 20 hours |
| **Mobile App** | P3 | 40 hours |

## 🔌 API Architecture

### REST API Endpoints

The backend exposes RESTful APIs for all operations:

- **Authentication**: `/api/auth/*` - Login, logout, token refresh
- **System Monitoring**: `/api/monitor/*` - System metrics and information
- **Docker Management**: `/api/docker/*` - Container operations
- **Storage**: `/api/storage/*` - Disk and storage management
- **Users**: `/api/users/*` - User management
- **Services**: `/api/services/*` - Systemd service management

### WebSocket Protocol

Real-time monitoring data is streamed through WebSocket:

- **Endpoint**: `ws://host:8888/api/monitor/ws`
- **Protocol**: JSON message format
- **Data Types**: CPU, memory, disk, network metrics
- **Update Rate**: 1 second intervals

See [API_DOCUMENTATION.md](./API_DOCUMENTATION.md) for detailed API specifications.

## 🎨 User Interface

### Page Structure

```
├── Login Page
├── Dashboard (System Overview)
├── Monitor Pages
│   ├── CPU Monitor
│   ├── Memory Monitor
│   ├── Disk Monitor
│   └── Network Monitor
├── Storage Management
├── Docker Container Management
├── Service Management
└── User Management
```

### Navigation

- **Sidebar**: Main navigation menu
- **Header**: User info, notifications, settings
- **Breadcrumbs**: Page navigation hierarchy
- **Quick Actions**: Common operations on each page

### Responsive Breakpoints

- **Desktop**: 1280px+ - Full functionality
- **Tablet**: 768px-1279px - Optimized layout
- **Mobile**: <768px - Simplified interface

## 🔒 Security Architecture

### Authentication Flow

1. User submits credentials to `/api/auth/login`
2. Backend validates credentials and returns JWT tokens
3. Frontend stores access token (24h) and refresh token (30d)
4. All API requests include access token in Authorization header
5. Access token refresh automatically on expiry
6. Logout clears tokens and redirects to login

### Authorization

- **Public Routes**: Login page only
- **Protected Routes**: All application pages
- **Token Validation**: JWT middleware on all protected endpoints
- **WebSocket Auth**: Token-based authentication for WebSocket connections

### Security Measures

- **JWT Token Expiration**: 24-hour access tokens
- **Refresh Tokens**: 30-day refresh tokens with rotation
- **CORS Protection**: Configurable origin whitelist
- **Rate Limiting**: (To be implemented) Request throttling
- **Password Hashing**: bcrypt for password storage
- **HTTPS Required**: Production deployment only

## 📈 Performance Considerations

### Frontend Optimization

- **Code Splitting**: Route-based chunking with Vite
- **Lazy Loading**: Components loaded on demand
- **Tree Shaking**: Unused code elimination
- **Asset Optimization**: Image and CSS compression
- **Caching**: Aggressive caching for static assets

### Backend Optimization

- **Goroutine Pool**: Concurrent request handling
- **Connection Pooling**: Database and Docker client pooling
- **Response Caching**: System metrics caching (1-second TTL)
- **WebSocket Throttling**: Rate-limited data streaming
- **Memory Management**: Efficient data structure usage

### Monitoring Performance

- **WebSocket Connections**: Single connection per client
- **Polling Fallback**: Automatic fallback for WebSocket failures
- **Data Sampling**: Configurable monitoring intervals
- **Resource Usage**: Minimal CPU and memory footprint

## 🚀 Deployment Architecture

### Development Environment

```
Frontend: http://localhost:5173 (Vite dev server)
Backend:  http://localhost:8888 (Go server)
```

### Production Environment

```
Frontend: nginx (Serve static files, reverse proxy)
Backend:  Go server (systemd service)
WebSocket: nginx proxy with WebSocket support
```

### Docker Deployment (Recommended)

```
Container 1: Backend (Go server)
Container 2: Frontend (nginx)
Network:    Bridge network for inter-container communication
Volumes:    Host system mounts for Docker socket and data persistence
```

## 📝 Development Status

### Current Version: 0.1.0 (Alpha)

**Completed:**
- ✅ Core authentication system
- ✅ Real-time monitoring infrastructure
- ✅ Basic container management
- ✅ Responsive UI framework
- ✅ API and WebSocket architecture

**In Progress:**
- ⚠️ Security hardening
- ⚠️ Error handling improvements
- ⚠️ Testing coverage

**Next Milestone: 0.2.0 (Beta)**
- Alert system implementation
- File manager integration
- Enhanced security measures
- Comprehensive testing

## 🤝 Contributing

This project is in active development. Contributions are welcome in the form of:
- Bug reports and issue submissions
- Feature requests and proposals
- Pull requests for improvements
- Documentation enhancements
- Test case development

See [DEVELOPMENT_GUIDE.md](./DEVELOPMENT_GUIDE.md) for detailed contribution guidelines.

## 📄 License

[Specify your license here]

## 📞 Support

For issues, questions, or contributions:
- **Issues**: [GitHub Issues](https://github.com/yourusername/nas-dashboard/issues)
- **Documentation**: [Project Wiki](https://github.com/yourusername/nas-dashboard/wiki)
- **Discussions**: [GitHub Discussions](https://github.com/yourusername/nas-dashboard/discussions)

---

**Last Updated**: 2026-06-12
**Version**: 0.1.0-alpha
