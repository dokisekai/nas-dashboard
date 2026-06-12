package websocket

import (
	"log"
	"nas-dashboard/pkg/system"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // 生产环境应该限制来源
		},
	}
)

// WebSocketManager WebSocket 管理器
type WebSocketManager struct {
	hub           *Hub
	monitorTicker *time.Ticker
	mu            sync.RWMutex
}

// NewWebSocketManager 创建 WebSocket 管理器
func NewWebSocketManager() *WebSocketManager {
	return &WebSocketManager{
		hub: NewHub(),
	}
}

// Start 启动 WebSocket 管理器
func (m *WebSocketManager) Start() {
	// 启动 Hub
	go m.hub.Run()

	// 启动监控数据推送
	m.startMonitorPush()

	log.Println("WebSocket manager started")
}

// Stop 停止 WebSocket 管理器
func (m *WebSocketManager) Stop() {
	if m.monitorTicker != nil {
		m.monitorTicker.Stop()
	}
	log.Println("WebSocket manager stopped")
}

// HandleWebSocket 处理 WebSocket 连接
func (m *WebSocketManager) HandleWebSocket(c *gin.Context) {
	// 获取用户信息
	username, exists := c.Get("username")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	role, _ := c.Get("role")

	// 升级连接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	// 创建客户端
	client := &Client{
		ID:       GenerateClientID(),
		Hub:      m.hub,
		Conn:     conn,
		Send:     make(chan *Message, 256),
		Username: username.(string),
		Role:     role.(string),
	}

	// 注册客户端
	m.hub.register <- client

	// 启动读写泵
	go client.writePump()
	go client.readPump()

	// 发送欢迎消息
	welcome := &Message{
		Type:      "connected",
		Data:      gin.H{"clientId": client.ID, "username": username},
		Timestamp: time.Now().Unix(),
	}
	client.Send <- welcome
}

// startMonitorPush 启动监控数据推送
func (m *WebSocketManager) startMonitorPush() {
	m.monitorTicker = time.NewTicker(2 * time.Second)

	go func() {
		for range m.monitorTicker.C {
			m.pushMonitorData()
		}
	}()
}

// pushMonitorPush 推送监控数据
func (m *WebSocketManager) pushMonitorData() {
	// 获取监控数据
	cpu, err := system.GetCPUInfo()
	if err != nil {
		log.Printf("Error getting CPU info: %v", err)
		return
	}

	mem, err := system.GetMemoryInfo()
	if err != nil {
		log.Printf("Error getting memory info: %v", err)
		return
	}

	disk, err := system.GetDiskInfo()
	if err != nil {
		log.Printf("Error getting disk info: %v", err)
		return
	}

	network, err := system.GetNetworkInfo()
	if err != nil {
		log.Printf("Error getting network info: %v", err)
		return
	}

	// 创建消息
	message := &Message{
		Type: "monitor_data",
		Data: gin.H{
			"cpu":     cpu,
			"memory":  mem,
			"disk":    disk,
			"network": network,
		},
		Timestamp: time.Now().Unix(),
	}

	// 广播消息
	m.hub.Broadcast(message)
}

// BroadcastSystemEvent 广播系统事件
func (m *WebSocketManager) BroadcastSystemEvent(eventType string, data interface{}) {
	message := &Message{
		Type:      eventType,
		Data:      data,
		Timestamp: time.Now().Unix(),
	}

	m.hub.Broadcast(message)
}

// BroadcastUserEvent 广播用户事件
func (m *WebSocketManager) BroadcastUserEvent(username string, eventType string, data interface{}) {
	message := &Message{
		Type:      eventType,
		Data:      data,
		Timestamp: time.Now().Unix(),
	}

	m.hub.BroadcastToUser(username, message)
}

// GetConnectedClients 获取连接的客户端数量
func (m *WebSocketManager) GetConnectedClients() int {
	return m.hub.GetClientCount()
}

// GetClientList 获取客户端列表
func (m *WebSocketManager) GetClientList() []map[string]interface{} {
	clients := m.hub.GetClients()
	result := make([]map[string]interface{}, 0, len(clients))

	for _, client := range clients {
		result = append(result, map[string]interface{}{
			"id":       client.ID,
			"username": client.Username,
			"role":     client.Role,
		})
	}

	return result
}

// DisconnectUser 断开指定用户的所有连接
func (m *WebSocketManager) DisconnectUser(username string) {
	clients := m.hub.GetClients()

	for _, client := range clients {
		if client.Username == username {
			client.Conn.Close()
		}
	}
}

// HandleWebSocketAPI 处理 WebSocket API 请求
func HandleWebSocketAPI(c *gin.Context) {
	wsManager := GetWebSocketManager()

	action := c.Param("action")

	switch action {
	case "clients":
		clients := wsManager.GetClientList()
		c.JSON(200, gin.H{
			"clients": clients,
			"total":   len(clients),
		})

	case "count":
		count := wsManager.GetConnectedClients()
		c.JSON(200, gin.H{"count": count})

	case "disconnect":
		var req struct {
			Username string `json:"username" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		wsManager.DisconnectUser(req.Username)
		c.JSON(200, gin.H{"message": "User disconnected"})

	default:
		c.JSON(400, gin.H{"error": "Unknown action"})
	}
}

// 全局 WebSocket 管理器实例
var globalWSManager *WebSocketManager

// InitWebSocketManager 初始化全局 WebSocket 管理器
func InitWebSocketManager() {
	if globalWSManager == nil {
		globalWSManager = NewWebSocketManager()
		globalWSManager.Start()
	}
}

// GetWebSocketManager 获取全局 WebSocket 管理器
func GetWebSocketManager() *WebSocketManager {
	return globalWSManager
}

// BroadcastMonitorData 广播监控数据（兼容旧接口）
func BroadcastMonitorData(c *gin.Context) {
	wsManager := GetWebSocketManager()

	// 手动触发一次监控数据推送
	wsManager.pushMonitorData()

	c.JSON(200, gin.H{"message": "Monitor data broadcasted"})
}
