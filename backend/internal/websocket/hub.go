package websocket

import (
	"log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Message WebSocket 消息
type Message struct {
	Type      string                 `json:"type"`
	Data      interface{}            `json:"data"`
	Timestamp int64                  `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// Client WebSocket 客户端
type Client struct {
	ID       string
	Hub      *Hub
	Conn     *websocket.Conn
	Send     chan *Message
	Username string
	Role     string
}

// Hub WebSocket 连接池
type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan *Message
	mu         sync.RWMutex
}

// NewHub 创建新的 Hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message, 256),
	}
}

// Run 运行 Hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("Client registered: %s (total: %d)", client.ID, len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			}
			h.mu.Unlock()
			log.Printf("Client unregistered: %s (total: %d)", client.ID, len(h.clients))

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					log.Printf("Client %s send channel is full, skipping", client.ID)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Broadcast 广播消息到所有客户端
func (h *Hub) Broadcast(message *Message) {
	h.broadcast <- message
}

// BroadcastToUser 广播消息到特定用户
func (h *Hub) BroadcastToUser(username string, message *Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		if client.Username == username {
			select {
			case client.Send <- message:
			default:
				log.Printf("Client %s send channel is full, skipping", client.ID)
			}
		}
	}
}

// BroadcastToRole 广播消息到特定角色
func (h *Hub) BroadcastToRole(role string, message *Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		if client.Role == role {
			select {
			case client.Send <- message:
			default:
				log.Printf("Client %s send channel is full, skipping", client.ID)
			}
		}
	}
}

// GetClientCount 获取客户端数量
func (h *Hub) GetClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// GetClients 获取所有客户端
func (h *Hub) GetClients() []*Client {
	h.mu.RLock()
	defer h.mu.RUnlock()

	clients := make([]*Client, 0, len(h.clients))
	for client := range h.clients {
		clients = append(clients, client)
	}
	return clients
}

// RemoveClient 移除客户端
func (h *Hub) RemoveClient(client *Client) {
	h.unregister <- client
}

// writePump 写入泵
func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				// Hub 关闭了通道
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteJSON(message); err != nil {
				log.Printf("Error writing to WebSocket: %v", err)
				return
			}

		case <-ticker.C:
			// 发送心跳
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// readPump 读取泵
func (c *Client) readPump() {
	defer func() {
		c.Hub.RemoveClient(c)
		c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		var message Message
		if err := c.Conn.ReadJSON(&message); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// 处理客户端消息
		c.handleMessage(&message)
	}
}

// handleMessage 处理客户端消息
func (c *Client) handleMessage(message *Message) {
	// 根据消息类型处理
	switch message.Type {
	case "ping":
		// 响应心跳
		pong := &Message{
			Type:      "pong",
			Timestamp: time.Now().Unix(),
		}
		c.Send <- pong

	case "subscribe":
		// 处理订阅请求
		c.handleSubscribe(message)

	case "unsubscribe":
		// 处理取消订阅请求
		c.handleUnsubscribe(message)

	default:
		log.Printf("Unknown message type: %s", message.Type)
	}
}

// handleSubscribe 处理订阅请求
func (c *Client) handleSubscribe(message *Message) {
	// 这里可以实现订阅特定频道的逻辑
	log.Printf("Client %s subscribed to channels", c.ID)
}

// handleUnsubscribe 处理取消订阅请求
func (c *Client) handleUnsubscribe(message *Message) {
	// 这里可以实现取消订阅特定频道的逻辑
	log.Printf("Client %s unsubscribed from channels", c.ID)
}

// GenerateClientID 生成客户端ID
func GenerateClientID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString 生成随机字符串
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
