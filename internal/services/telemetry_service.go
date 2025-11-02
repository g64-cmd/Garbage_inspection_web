package services

import (
	"log"
	"sync"
	"github.com/gorilla/websocket"
)

// TelemetryHub 遵循 4.2.7 的设计，作为 WebSocket 广播中枢
type TelemetryHub struct {
	// (使用 Client 结构体代替 conn，以便更好地管理状态)
	Clients map[*TelemetryClient]bool
	// BroadcastChannel (从 MQTTListener 接收消息)
	BroadcastChannel chan []byte
	// Register (注册来自 WebSocketHandler 的新连接)
	Register chan *TelemetryClient
	// Unregister (注销断开的连接)
	Unregister chan *TelemetryClient

	// 确保并发安全
	clientsMutex sync.RWMutex
}

// TelemetryClient 是 Hub 管理的 WebSocket 客户端的包装器
type TelemetryClient struct {
	hub  *TelemetryHub
	conn *websocket.Conn
	// (Send 缓冲通道，防止广播时写入阻塞)
	send chan []byte
}

func NewTelemetryHub() *TelemetryHub {
	return &TelemetryHub{
		Clients:          make(map[*TelemetryClient]bool),
		BroadcastChannel: make(chan []byte, 256), // (带缓冲的通道)
		Register:         make(chan *TelemetryClient),
		Unregister:       make(chan *TelemetryClient),
	}
}

func NewTelemetryClient(hub *TelemetryHub, conn *websocket.Conn) *TelemetryClient {
	return &TelemetryClient{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}
}

// Run 是 Hub 的核心 Goroutine，使用 select 循环处理事件
func (h *TelemetryHub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.clientsMutex.Lock()
			h.Clients[client] = true
			h.clientsMutex.Unlock()
			log.Println("INFO: WebSocket client registered")

		case client := <-h.Unregister:
			h.clientsMutex.Lock()
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.send)
				log.Println("INFO: WebSocket client unregistered")
			}
			h.clientsMutex.Unlock()

		case message := <-h.BroadcastChannel:
			// (从 MQTTListener 收到消息，广播给所有客户端)
			h.clientsMutex.RLock()
			for client := range h.Clients {
				select {
				case client.send <- message:
					// 消息已发送到客户端的缓冲
				default:
					// 缓冲已满，丢弃消息并注销客户端
					log.Println("WARN: Client send buffer full, unregistering")
					go func() { h.Unregister <- client }()
				}
			}
			h.clientsMutex.RUnlock()
		}
	}
}

// (RegisterClient 是一个辅助方法，供 Handler 调用)
func (h *TelemetryHub) RegisterClient(client *TelemetryClient) {
	h.Register <- client
}

// --- Client Goroutines (由 WebSocketHandler 启动) ---

// ReadLoop (4.2.2 提及) 侦听客户端断开连接
func (c *TelemetryClient) ReadLoop() {
	defer func() {
		c.hub.Unregister <- c
		c.conn.Close()
	}()
	// (设置 Ping/Pong 或 ReadDeadline 来检测死连接)
	for {
		// (我们只关心读错误，不关心消息内容)
		if _, _, err := c.conn.ReadMessage(); err != nil {
			log.Printf("INFO: WebSocket read error (client disconnected): %v", err)
			break
		}
	}
}

// WriteLoop 将缓冲区的消息写入 WebSocket
func (c *TelemetryClient) WriteLoop() {
	defer c.conn.Close()
	for message := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("ERROR: WebSocket write error: %v", err)
			return // (ReadLoop 会处理注销)
		}
	}
}
