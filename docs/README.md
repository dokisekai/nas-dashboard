# NAS Dashboard - Complete Documentation

Welcome to the NAS Dashboard comprehensive documentation suite. This modern web-based monitoring and management system is designed for Network Attached Storage (NAS) devices and Linux servers.

## 🌟 Overview

**NAS Dashboard** is a full-featured system management platform that provides:

- **Real-time System Monitoring**: Live CPU, memory, disk, and network metrics
- **Docker Container Management**: Start, stop, and monitor containers
- **Storage Management**: Disk monitoring, SMB shares, and file system operations
- **User Administration**: Complete user management with SSH key support
- **Desktop Environment**: Modern DSM-style interface with widgets and windows
- **Plugin System**: Extensible architecture for custom functionality
- **Application Center**: Easy installation of system applications

## 🚀 Quick Start

### Prerequisites

- **Go**: 1.22+ (backend)
- **Node.js**: 20+ (frontend)
- **Docker**: 20.10+ (optional, for containerized deployment)
- **PostgreSQL**: 14+ (for database features)

### Installation Options

#### 1. Docker Deployment (Recommended)

```bash
# Clone the repository
git clone https://github.com/yourusername/nas-dashboard.git
cd nas-dashboard

# Start with Docker Compose
docker-compose up -d

# Access the application
open http://localhost:3000
```

**Default Credentials:**
- Username: `admin`
- Password: `admin123`

#### 2. Manual Installation

```bash
# Backend setup
cd backend
go mod download
go run cmd/server/main.go

# Frontend setup (new terminal)
cd frontend
npm install
npm run dev
```

Access at http://localhost:5173

#### 3. Production Deployment

See [DEPLOYMENT.md](./DEPLOYMENT.md) for detailed production deployment instructions.

## 📚 Documentation Navigation

### For Users

1. **[USER_GUIDE.md](./USER_GUIDE.md)** - End-user documentation
   - Interface overview
   - Daily operations
   - Desktop features
   - Application management

### For Developers

1. **[DEVELOPER_GUIDE.md](./DEVELOPER_GUIDE.md)** - Development setup
   - Environment setup
   - Code organization
   - Testing procedures
   - Contribution guidelines

2. **[PLUGIN_DEVELOPMENT.md](./PLUGIN_DEVELOPMENT.md)** - Plugin development
   - Plugin architecture
   - SDK reference
   - Development examples
   - Best practices

3. **[API.md](./API.md)** - API documentation
   - REST endpoints
   - WebSocket protocol
   - Authentication
   - Data models

### For System Administrators

1. **[DEPLOYMENT.md](./DEPLOYMENT.md)** - Deployment guide
   - Production setup
   - Security configuration
   - Performance optimization
   - Monitoring

2. **[ARCHITECTURE.md](./ARCHITECTURE.md)** - System architecture
   - Architecture overview
   - Component design
   - Data flow
   - Technology choices

### For Security

1. **[SECURITY_CONSIDERATIONS.md](./docs/SECURITY_CONSIDERATIONS.md)** - Security guide
   - Security best practices
   - Vulnerability assessment
   - Hardening procedures
   - Audit requirements

## 🏗️ Project Structure

```
nas-dashboard/
├── backend/                 # Go backend application
│   ├── cmd/                # Application entry points
│   ├── internal/           # Private application code
│   │   ├── api/           # API handlers
│   │   ├── middleware/    # HTTP middleware
│   │   ├── models/        # Data models
│   │   ├── services/       # Business logic
│   │   └── websocket/     # WebSocket handlers
│   ├── pkg/               # Public packages
│   └── config/            # Configuration files
├── frontend/              # Vue 3 frontend
│   ├── src/
│   │   ├── api/          # API clients
│   │   ├── apps/         # Desktop applications
│   │   ├── components/   # Vue components
│   │   ├── stores/       # Pinia stores
│   │   ├── composables/  # Composition functions
│   │   ├── types/        # TypeScript definitions
│   │   └── utils/        # Utilities
│   └── public/           # Static files
├── docs/                 # Documentation (this directory)
└── docker-compose.yml    # Container orchestration
```

## 🔧 Technology Stack

### Frontend
- **Framework**: Vue 3.5+ with Composition API
- **Language**: TypeScript 6.0+
- **Build Tool**: Vite 8.0+
- **State Management**: Pinia 3.0+
- **Routing**: Vue Router 4.6+
- **Styling**: Tailwind CSS 4.3+
- **Charts**: Chart.js 4.5+, ApexCharts 5.15+
- **Drag & Drop**: vue-draggable-plus 0.6+
- **Utilities**: @vueuse/core 14.3+

### Backend
- **Language**: Go 1.22+
- **Framework**: Gin (web framework)
- **Database**: PostgreSQL with GORM
- **Authentication**: JWT (bcrypt password hashing)
- **System Monitoring**: gopsutil
- **Container Management**: Docker SDK
- **WebSocket**: gorilla/websocket

## 🎯 Key Features

### 1. Desktop Environment
- **Window Management**: Drag, resize, minimize, maximize windows
- **Widget System**: Real-time monitoring widgets on desktop
- **Dock Bar**: Quick application access and management
- **Theme Support**: Light/dark modes with custom backgrounds
- **Context Menus**: Right-click menus for quick actions

### 2. System Monitoring
- **Real-time Metrics**: Live CPU, memory, disk, network data
- **Historical Data**: Time-series data with charts
- **Process Management**: View and manage system processes
- **Alert System**: Configurable thresholds and notifications
- **WebSocket Streaming**: Efficient real-time data delivery

### 3. Application Center
- **Storage Manager**: Disk management, SMB shares, file browser
- **System Monitor**: Comprehensive system monitoring
- **User Manager**: User administration with SSH keys
- **System Settings**: Network, services, backup management
- **Plugin Store**: Browse and install plugins
- **App Center**: Application marketplace

### 4. Plugin System
- **Dynamic Loading**: Load plugins at runtime
- **SDK**: Comprehensive plugin development kit
- **Lifecycle Management**: Install, enable, disable, uninstall
- **Permission System**: Granular permission control
- **Marketplace**: Plugin distribution and discovery

### 5. Security Features
- **JWT Authentication**: Secure token-based authentication
- **Password Hashing**: bcrypt for secure password storage
- **Session Management**: Multi-device session control
- **Audit Logging**: Complete operation audit trail
- **Access Control**: Role-based permissions
- **CORS Protection**: Configurable cross-origin policies

## 📊 API Endpoints

### Authentication
- `POST /api/auth/login` - User login
- `POST /api/auth/refresh` - Refresh access token
- `POST /api/auth/logout` - User logout
- `GET /api/auth/me` - Get current user

### Monitoring
- `GET /api/monitor/cpu` - CPU information
- `GET /api/monitor/memory` - Memory information
- `GET /api/monitor/disk` - Disk usage
- `GET /api/monitor/network` - Network statistics
- `GET /api/monitor/system` - System information
- `WS /api/monitor/ws` - Real-time WebSocket

### Storage
- `GET /api/storage` - Storage information
- `POST /api/storage/scan` - Scan storage
- `GET /api/storage/disks` - List disks
- `POST /api/storage/mount` - Mount disk
- `POST /api/storage/umount` - Unmount disk

### Docker
- `GET /api/docker/containers` - List containers
- `POST /api/docker/containers/{id}/start` - Start container
- `POST /api/docker/containers/{id}/stop` - Stop container
- `POST /api/docker/containers/{id}/restart` - Restart container
- `GET /api/docker/containers/{id}/logs` - Container logs

### Users
- `GET /api/users` - List users
- `POST /api/users` - Create user
- `PUT /api/users/{id}` - Update user
- `DELETE /api/users/{id}` - Delete user
- `PUT /api/users/{id}/password` - Change password

### Services
- `GET /api/services` - List services
- `POST /api/services/{name}/start` - Start service
- `POST /api/services/{name}/stop` - Stop service
- `POST /api/services/{name}/restart` - Restart service

## 🔐 Security Best Practices

### Production Deployment Checklist

- [ ] Change default admin password
- [ ] Set strong JWT secret in environment variables
- [ ] Enable HTTPS with valid SSL certificate
- [ ] Configure firewall rules
- [ ] Set up database backups
- [ ] Enable audit logging
- [ ] Configure rate limiting
- [ ] Review CORS settings
- [ ] Implement monitoring and alerting
- [ ] Regular security updates

See [SECURITY_CONSIDERATIONS.md](./docs/SECURITY_CONSIDERATIONS.md) for detailed security information.

## 🛠️ Configuration

### Environment Variables

**Backend (`.env`):**
```bash
# Server
PORT=8888
GIN_MODE=release

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=nas_dashboard
DB_USER=nas_user
DB_PASSWORD=your_secure_password

# JWT
JWT_SECRET=your_jwt_secret_key_change_this
JWT_ACCESS_DURATION=24h
JWT_REFRESH_DURATION=720h

# CORS
CORS_ORIGIN=https://your-domain.com

# Logging
LOG_LEVEL=info
```

**Frontend (`.env`):**
```bash
VITE_API_URL=https://api.your-domain.com
VITE_WS_URL=wss://api.your-domain.com
VITE_APP_NAME=NAS Dashboard
VITE_DEBUG=false
```

## 📈 Performance Optimization

### Backend
- Database connection pooling
- WebSocket connection management
- Efficient system monitoring queries
- Lazy loading for large datasets
- Response caching for static data

### Frontend
- Code splitting and lazy loading
- Virtual scrolling for large lists
- Image optimization and lazy loading
- Debounced API calls
- Efficient state management

## 🐛 Troubleshooting

### Common Issues

**1. Cannot connect to backend**
- Check backend is running on correct port
- Verify CORS settings
- Check firewall rules

**2. Authentication fails**
- Verify JWT credentials
- Check token expiration
- Review browser console for errors

**3. WebSocket connection issues**
- Verify WebSocket URL
- Check proxy configuration
- Review network connectivity

**4. Docker management not working**
- Verify Docker service running
- Check API socket permissions
- Review Docker daemon configuration

For detailed troubleshooting, see [TROUBLESHOOTING.md](./docs/TROUBLESHOOTING.md).

## 🤝 Contributing

We welcome contributions! Please see [DEVELOPER_GUIDE.md](./DEVELOPER_GUIDE.md) for:

- Development environment setup
- Coding standards
- Testing procedures
- Pull request process
- Code review guidelines

## 📝 Version History

### Version 0.1.0 (Current)
- Initial release
- Basic system monitoring
- Docker container management
- User management
- Desktop environment
- Plugin system
- Application center

## 📞 Support

### Getting Help

- **Documentation**: This comprehensive guide
- **Issues**: [GitHub Issues](https://github.com/yourusername/nas-dashboard/issues)
- **Discussions**: [GitHub Discussions](https://github.com/yourusername/nas-dashboard/discussions)
- **Email**: support@example.com

### Reporting Issues

When reporting issues, please include:
- NAS Dashboard version
- Operating system and version
- Browser and version (if applicable)
- Steps to reproduce
- Expected vs actual behavior
- Relevant logs and screenshots

## 📄 License

MIT License - see [LICENSE](../LICENSE) for details.

## 🙏 Acknowledgments

- **Vue.js Team** - Excellent frontend framework
- **Gin Team** - Powerful Go web framework
- **Tailwind CSS** - Amazing utility-first CSS framework
- **Chart.js** - Beautiful charting library
- **Open Source Community** - All the amazing open-source tools used

---

## 🎯 Quick Reference

### Start Development
```bash
# Backend
cd backend && go run cmd/server/main.go

# Frontend
cd frontend && npm run dev
```

### Build for Production
```bash
# Backend
cd backend && go build -o nas-dashboard cmd/server/main.go

# Frontend
cd frontend && npm run build
```

### Docker Deployment
```bash
docker-compose up -d
```

### Run Tests
```bash
# Backend
cd backend && go test ./...

# Frontend
cd frontend && npm run test
```

---

**Last Updated**: 2026-06-12  
**Version**: 0.1.0  
**Status**: Alpha (Active Development)

---

## 📖 Recommended Reading Order

### For New Users
1. [Quick Start](#quick-start)
2. [USER_GUIDE.md](./USER_GUIDE.md)
3. [API.md](./API.md) (if integrating with other systems)

### For Developers
1. [DEVELOPER_GUIDE.md](./DEVELOPER_GUIDE.md)
2. [ARCHITECTURE.md](./ARCHITECTURE.md)
3. [API.md](./API.md)
4. [PLUGIN_DEVELOPMENT.md](./PLUGIN_DEVELOPMENT.md)

### For System Administrators
1. [DEPLOYMENT.md](./DEPLOYMENT.md)
2. [ARCHITECTURE.md](./ARCHITECTURE.md)
3. [SECURITY_CONSIDERATIONS.md](./docs/SECURITY_CONSIDERATIONS.md)
4. [docs/TROUBLESHOOTING.md](./docs/TROUBLESHOOTING.md)

---

**Next Steps**: Choose your role above and follow the recommended reading order, or jump to any specific document using the navigation provided.
