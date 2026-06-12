# NAS Dashboard - Security Considerations

Comprehensive security guide for deploying and maintaining the NAS Dashboard.

## 📋 Table of Contents

1. [Security Overview](#security-overview)
2. [Authentication Security](#authentication-security)
3. [API Security](#api-security)
4. [WebSocket Security](#websocket-security)
5. [Data Protection](#data-protection)
6. [Network Security](#network-security)
7. [Production Deployment Security](#production-deployment-security)
8. [Known Vulnerabilities](#known-vulnerabilities)
9. [Security Best Practices](#security-best-practices)

---

## Security Overview

### Current Security Posture

The NAS Dashboard implements several security measures but has known vulnerabilities that must be addressed before production deployment.

### Security Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    Security Layers                      │
├─────────────────────────────────────────────────────────┤
│  1. Network Layer        (Firewall, TLS, CORS)        │
│  2. Authentication       (JWT, Token Management)      │
│  3. Authorization        (Role-Based Access)           │
│  4. Application Logic   (Input Validation)             │
│  5. Data Layer          (Encryption, Sanitization)    │
└─────────────────────────────────────────────────────────┘
```

---

## Authentication Security

### Current Implementation

#### JWT Token Configuration

**Current Settings:**
```bash
JWT_ACCESS_DURATION=24h      # 24-hour access tokens
JWT_REFRESH_DURATION=720h    # 30-day refresh tokens
```

**Security Assessment:**
- ✅ Proper token expiration
- ✅ Refresh token mechanism
- ❌ **CRITICAL**: JWT secret must be changed from default
- ❌ **HIGH**: No token rotation on refresh

#### Recommended Changes

**File: `/home/hserver/nas-dashboard/backend/pkg/jwt/jwt.go`**

```go
// Current (INSECURE)
var jwtSecret = []byte("your-secret-key")

// Recommended (SECURE)
import (
    "crypto/rand"
    "encoding/base64"
    "os"
)

func getJWTSecret() []byte {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        panic("JWT_SECRET must be set")
    }
    return []byte(secret)
}

// Generate secure secret for new installations
func generateSecureSecret() (string, error) {
    bytes := make([]byte, 32)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return base64.URLEncoding.EncodeToString(bytes), nil
}
```

#### Password Security

**Current Implementation:**
```go
// File: /home/hserver/nas-dashboard/backend/internal/api/auth.go
var users = map[string]string{
    "admin": "admin123",  // PLAINTEXT PASSWORD - INSECURE
}
```

**Security Issues:**
- ❌ **CRITICAL**: Passwords stored in plaintext
- ❌ **CRITICAL**: Hardcoded users in memory
- ❌ **HIGH**: No password complexity requirements
- ❌ **MEDIUM**: No account lockout mechanism

**Secure Implementation:**

```go
import (
    "golang.org/x/crypto/bcrypt"
    "github.com/google/uuid"
)

type User struct {
    ID           uuid.UUID `json:"id"`
    Username     string    `json:"username"`
    PasswordHash string    `json:"-"`
    Email        string    `json:"email"`
    Role         string    `json:"role"`
    CreatedAt    time.Time `json:"created_at"`
    LastLogin    time.Time `json:"last_login"`
}

// Hash password
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword(
        []byte(password), 
        bcrypt.DefaultCost, // Cost factor 10
    )
    return string(bytes), err
}

// Verify password
func VerifyPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword(
        []byte(hash), 
        []byte(password),
    )
    return err == nil
}

// Validate password strength
func ValidatePassword(password string) error {
    if len(password) < 8 {
        return errors.New("password must be at least 8 characters")
    }
    
    hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
    hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
    hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
    hasSpecial := regexp.MustCompile(`[!@#$%^&*]`).MatchString(password)
    
    if !hasUpper || !hasLower || !hasDigit || !hasSpecial {
        return errors.New(
            "password must contain uppercase, lowercase, digit, and special character",
        )
    }
    
    return nil
}

// Account lockout
type FailedLoginAttempt struct {
    Username  string
    Attempts  int
    LastAttempt time.Time
    LockedUntil time.Time
}

var failedLogins = make(map[string]*FailedLoginAttempt)
const maxAttempts = 5
const lockoutDuration = 30 * time.Minute

func IsAccountLocked(username string) bool {
    attempt, exists := failedLogins[username]
    if !exists {
        return false
    }
    
    if time.Now().Before(attempt.LockedUntil) {
        return true
    }
    
    // Reset if lockout expired
    delete(failedLogins, username)
    return false
}

func RecordFailedLogin(username string) {
    attempt, exists := failedLogins[username]
    if !exists {
        attempt = &FailedLoginAttempt{Username: username}
        failedLogins[username] = attempt
    }
    
    attempt.Attempts++
    attempt.LastAttempt = time.Now()
    
    if attempt.Attempts >= maxAttempts {
        attempt.LockedUntil = time.Now().Add(lockoutDuration)
    }
}
```

#### Token Management

**Secure Token Refresh:**

```go
type TokenPayload struct {
    UserID    uuid.UUID `json:"user_id"`
    Username  string    `json:"username"`
    Role      string    `json:"role"`
    TokenID   string    `json:"token_id"` // Unique identifier
}

// Generate tokens
func GenerateTokens(user User) (accessToken, refreshToken string, err error) {
    tokenID := uuid.New().String()
    
    accessPayload := TokenPayload{
        UserID:   user.ID,
        Username: user.Username,
        Role:     user.Role,
        TokenID:  tokenID,
    }
    
    accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "payload": accessPayload,
        "exp":     time.Now().Add(24 * time.Hour).Unix(),
        "type":    "access",
    }).SignedString(getJWTSecret())
    
    refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id":  user.ID.String(),
        "token_id": tokenID,
        "exp":      time.Now().Add(30 * 24 * time.Hour).Unix(),
        "type":     "refresh",
    }).SignedString(getJWTSecret())
    
    return
}

// Token blacklist for logout
var tokenBlacklist = make(map[string]bool)

func InvalidateToken(tokenString string) error {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return getJWTSecret(), nil
    })
    
    if err != nil {
        return err
    }
    
    if claims, ok := token.Claims.(jwt.MapClaims); ok {
        tokenID := claims["token_id"].(string)
        tokenBlacklist[tokenID] = true
        return nil
    }
    
    return errors.New("invalid token")
}

func IsTokenBlacklisted(tokenID string) bool {
    return tokenBlacklist[tokenID]
}
```

---

## API Security

### CORS Configuration

**Current Implementation:**
```go
// File: /home/hserver/nas-dashboard/backend/internal/middleware/cors.go
c.Header("Access-Control-Allow-Origin", "*")
```

**Security Issues:**
- ❌ **CRITICAL**: Allows all origins
- ❌ **HIGH**: No origin validation
- ❌ **MEDIUM**: No credential handling

**Secure Implementation:**

```go
func CORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        origin := c.Request.Header.Get("Origin")
        
        // Validate origin
        allowedOrigins := []string{
            "https://your-domain.com",
            "https://www.your-domain.com",
        }
        
        allowed := false
        for _, allowedOrigin := range allowedOrigins {
            if origin == allowedOrigin {
                allowed = true
                break
            }
        }
        
        if allowed {
            c.Header("Access-Control-Allow-Origin", origin)
            c.Header("Access-Control-Allow-Credentials", "true")
        }
        
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
        c.Header("Access-Control-Max-Age", "86400")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        
        c.Next()
    }
}
```

### Rate Limiting

**Current Status:** ❌ Missing

**Recommended Implementation:**

```go
import (
    "sync"
    "time"
)

type RateLimiter struct {
    clients map[string]*ClientInfo
    mu      sync.RWMutex
}

type ClientInfo struct {
    Requests    []time.Time
    LastRequest time.Time
}

var (
    limiter = RateLimiter{
        clients: make(map[string]*ClientInfo),
    }
    
    requestsPerMinute = 100
    cleanupInterval   = 5 * time.Minute
)

func RateLimitMiddleware() gin.HandlerFunc {
    go cleanupStaleClients()
    
    return func(c *gin.Context) {
        clientIP := c.ClientIP()
        
        limiter.mu.Lock()
        defer limiter.mu.Unlock()
        
        client, exists := limiter.clients[clientIP]
        if !exists {
            client = &ClientInfo{
                Requests: make([]time.Time, 0),
            }
            limiter.clients[clientIP] = client
        }
        
        now := time.Now()
        
        // Remove old requests (older than 1 minute)
        recentRequests := make([]time.Time, 0)
        for _, req := range client.Requests {
            if now.Sub(req) < time.Minute {
                recentRequests = append(recentRequests, req)
            }
        }
        client.Requests = recentRequests
        
        // Check rate limit
        if len(client.Requests) >= requestsPerMinute {
            c.JSON(429, gin.H{
                "error": "Rate limit exceeded",
                "retry_after": "60s",
            })
            c.Abort()
            return
        }
        
        client.Requests = append(client.Requests, now)
        client.LastRequest = now
        
        c.Next()
    }
}

func cleanupStaleClients() {
    ticker := time.NewTicker(cleanupInterval)
    defer ticker.Stop()
    
    for range ticker.C {
        limiter.mu.Lock()
        now := time.Now()
        
        for ip, client := range limiter.clients {
            if now.Sub(client.LastRequest) > cleanupInterval {
                delete(limiter.clients, ip)
            }
        }
        
        limiter.mu.Unlock()
    }
}
```

### Input Validation

**Current Status:** ⚠️ Partial

**Secure Implementation:**

```go
import (
    "regexp"
    "strings"
    "html"
)

func ValidateUsername(username string) error {
    if len(username) < 3 || len(username) > 32 {
        return errors.New("username must be 3-32 characters")
    }
    
    matched, err := regexp.MatchString("^[a-zA-Z0-9_-]+$", username)
    if err != nil || !matched {
        return errors.New("username can only contain letters, numbers, underscores, and hyphens")
    }
    
    return nil
}

func SanitizeInput(input string) string {
    // Remove HTML tags
    input = html.EscapeString(input)
    
    // Trim whitespace
    input = strings.TrimSpace(input)
    
    // Remove null bytes
    input = strings.ReplaceAll(input, "\x00", "")
    
    return input
}

// Validate Docker container ID
func ValidateContainerID(id string) error {
    matched, err := regexp.MatchString("^[a-f0-9]{64}$", id)
    if err != nil || !matched {
        return errors.New("invalid container ID format")
    }
    return nil
}
```

---

## WebSocket Security

### Current Implementation

**File: `/home/hserver/nas-dashboard/backend/internal/api/monitor.go`**

```go
func WSMonitor(c *gin.Context) {
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    // NO AUTHENTICATION CHECK!
}
```

**Security Issues:**
- ❌ **CRITICAL**: No authentication before WebSocket upgrade
- ❌ **HIGH**: No origin validation
- ❌ **MEDIUM**: No rate limiting on connections

**Secure Implementation:**

```go
func WSMonitor(c *gin.Context) {
    // Extract token from query string
    token := c.Query("token")
    if token == "" {
        c.JSON(401, gin.H{"error": "Authentication required"})
        return
    }
    
    // Validate token
    claims, err := validateJWTToken(token)
    if err != nil {
        c.JSON(401, gin.H{"error": "Invalid token"})
        return
    }
    
    // Check origin
    origin := c.Request.Header.Get("Origin")
    if !isValidOrigin(origin) {
        c.JSON(403, gin.H{"error": "Invalid origin"})
        return
    }
    
    // Check connection limit
    if getConnectionCount(claims.UserID) >= maxConnectionsPerUser {
        c.JSON(429, gin.H{"error": "Too many connections"})
        return
    }
    
    // Upgrade to WebSocket
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to upgrade connection"})
        return
    }
    
    // Register connection
    registerConnection(claims.UserID, conn)
    
    // Handle connection
    defer func() {
        unregisterConnection(claims.UserID, conn)
        conn.Close()
    }()
    
    for {
        var msg WebSocketMessage
        if err := conn.ReadJSON(&msg); err != nil {
            break
        }
        
        // Validate message
        if !isValidMessage(msg) {
            conn.WriteJSON(gin.H{"error": "Invalid message"})
            continue
        }
        
        // Process message...
    }
}
```

---

## Data Protection

### Sensitive Data Handling

**Current Issues:**
- ❌ User credentials in memory
- ❌ No encryption for sensitive data
- ❌ Logs may contain sensitive information

**Secure Practices:**

```go
import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
)

// Encrypt sensitive data
func EncryptData(plaintext []byte, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    
    nonce := make([]byte, gcm.NonceSize())
    if _, err := rand.Read(nonce); err != nil {
        return nil, err
    }
    
    ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
    return ciphertext, nil
}

// Decrypt data
func DecryptData(ciphertext []byte, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    
    nonceSize := gcm.NonceSize()
    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
    
    return gcm.Open(nil, nonce, ciphertext, nil)
}

// Secure logging (remove sensitive data)
type SecureLogger struct {
    logger *log.Logger
}

func (sl *SecureLogger) Printf(format string, v ...interface{}) {
    sanitized := sanitizeLogData(v...)
    sl.logger.Printf(format, sanitized...)
}

func sanitizeLogData(v ...interface{}) []interface{} {
    result := make([]interface{}, len(v))
    for i, item := range v {
        if str, ok := item.(string); ok {
            result[i] = sanitizeString(str)
        } else {
            result[i] = item
        }
    }
    return result
}

func sanitizeString(s string) string {
    // Remove potential tokens
    if len(s) > 20 && strings.HasPrefix(s, "ey") {
        return "[REDACTED]"
    }
    return s
}
```

---

## Network Security

### Firewall Configuration

**Recommended Rules:**

```bash
# Configure UFW (Ubuntu/Debian)
sudo ufw default deny incoming
sudo ufw default allow outgoing

# Allow SSH
sudo ufw allow 22/tcp

# Allow HTTP/HTTPS
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Allow application port (if not behind proxy)
sudo ufw allow 8888/tcp

# Rate limiting for SSH
sudo ufw limit 22/tcp

# Enable firewall
sudo ufw enable

# Check status
sudo ufw status
```

### SSL/TLS Configuration

**Nginx Configuration:**

```nginx
server {
    listen 443 ssl http2;
    server_name your-domain.com;

    # SSL certificates
    ssl_certificate /etc/ssl/certs/your-domain.crt;
    ssl_certificate_key /etc/ssl/private/your-domain.key;

    # SSL configuration
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers 'ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256';
    ssl_prefer_server_ciphers off;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;

    # HSTS
    add_header Strict-Transport-Security "max-age=31536000" always;

    # Other security headers
    add_header X-Frame-Options "DENY" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;

    # Location blocks...
}
```

---

## Production Deployment Security

### Security Checklist

#### Pre-Deployment

- [ ] Change JWT_SECRET from default
- [ ] Implement password hashing (bcrypt)
- [ ] Enable HTTPS with valid SSL certificate
- [ ] Configure firewall rules
- [ ] Set up fail2ban for brute-force protection
- [ ] Disable debug mode
- [ ] Implement rate limiting
- [ ] Configure CORS to specific origins
- [ ] Set up security headers
- [ ] Implement log rotation

#### Post-Deployment

- [ ] Monitor security logs
- [ ] Regular security updates
- [ ] Backup encryption
- [ ] Intrusion detection
- [ ] Security audit logs
- [ ] Performance monitoring

### Environment Configuration

**Production `.env`:**

```bash
# Server Mode
GIN_MODE=release
LOG_LEVEL=warn

# Security
JWT_SECRET=<generate-secure-secret>
JWT_ACCESS_DURATION=1h
JWT_REFRESH_DURATION=168h

# CORS
CORS_ORIGIN=https://your-domain.com

# Rate Limiting
RATE_LIMIT_PER_MINUTE=60
RATE_LIMIT_BURST=10

# Security Headers
ENABLE_HSTS=true
ENABLE_CSP=true
```

### Secure Secret Generation

```bash
# Generate secure JWT secret
openssl rand -base64 32

# Generate encryption key
openssl rand -hex 32

# Generate SSL certificate
sudo certbot certonly --standalone -d your-domain.com
```

---

## Known Vulnerabilities

### Critical Issues

| Issue | Severity | Status | Fix Required |
|-------|----------|--------|---------------|
| Default JWT secret | Critical | Open | Generate secure secret |
| Plaintext passwords | Critical | Open | Implement bcrypt |
| No WebSocket auth | Critical | Open | Add token validation |
| Hardcoded users | High | Open | Implement database |
| Missing rate limiting | High | Open | Add rate limiter |
| CORS allows all origins | High | Open | Configure specific origins |

### Vulnerability Timeline

| Version | Status | Vulnerabilities | Timeline |
|---------|--------|-----------------|----------|
| 0.1.0 | Current | 6 critical/high | Immediate fix |
| 0.2.0 | Planned | 0 | Q1 2024 |

---

## Security Best Practices

### Regular Security Tasks

**Daily:**
- Monitor application logs
- Check for failed login attempts
- Verify system integrity

**Weekly:**
- Review security logs
- Check for unusual activity
- Backup verification

**Monthly:**
- Update dependencies
- Security audit
- Penetration testing
- Review access controls

**Quarterly:**
- Full security review
- Incident response drill
- Security training

### Dependency Management

```bash
# Check for vulnerabilities
npm audit
go list -json -m all | nancy sleuth

# Update dependencies
npm update
go get -u ./...
go mod tidy

# Automated security scanning
trivy fs .
```

### Incident Response

**Security Incident Response Plan:**

1. **Detection**
   - Monitor logs for anomalies
   - Set up alerts for suspicious activity
   - Regular security scans

2. **Containment**
   - Isolate affected systems
   - Block malicious IPs
   - Disable compromised accounts

3. **Eradication**
   - Remove malware or backdoors
   - Patch vulnerabilities
   - Update security measures

4. **Recovery**
   - Restore from clean backups
   - Verify system integrity
   - Monitor for recurrence

5. **Lessons Learned**
   - Document incident
   - Update procedures
   - Train staff

---

## Security Monitoring

### Log Monitoring

```bash
# Monitor authentication attempts
tail -f /var/log/nas-dashboard/auth.log | grep -i "failed\|success"

# Monitor API access
tail -f /var/log/nas-dashboard/api.log

# Monitor WebSocket connections
tail -f /var/log/nas-dashboard/websocket.log
```

### Intrusion Detection

```bash
# Install fail2ban
sudo apt install fail2ban -y

# Configure for NAS Dashboard
sudo nano /etc/fail2ban/jail.local
```

**fail2ban configuration:**
```ini
[nas-dashboard-auth]
enabled = true
port = http,https
filter = nas-dashboard-auth
logpath = /var/log/nas-dashboard/auth.log
maxretry = 5
bantime = 3600
findtime = 600
```

---

## Compliance

### Data Protection

- **GDPR**: User data protection and privacy
- **SOC 2**: Security controls and monitoring
- **ISO 27001**: Information security management

### Security Standards

- **OWASP Top 10**: Web application security
- **NIST Framework**: Cybersecurity guidelines
- **CIS Controls**: Security best practices

---

## Resources

### Security Tools

- **OWASP ZAP**: Web application security scanner
- **Nmap**: Network security scanner
- **Wireshark**: Network protocol analyzer
- **Metasploit**: Penetration testing framework

### Documentation

- [OWASP Security Guidelines](https://owasp.org/www-project-top-ten/)
- [Go Security Best Practices](https://golang.org/doc/security)
- [Vue Security Guide](https://vuejs.org/guide/best-practices/security.html)

---

**Last Updated**: 2026-06-12
**Security Version**: 0.1.0-alpha
**Next Security Review**: 2026-07-12
