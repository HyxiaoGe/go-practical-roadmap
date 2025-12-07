package websocket

import (
	"encoding/json"
	"time"

	"github.com/yourname/02-concurrency-worker/internal/model"
	"github.com/yourname/02-concurrency-worker/internal/pkg/logger"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// Hub WebSocket中心
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan *TaskUpdateMessage
	register   chan *Client
	unregister chan *Client
}

// Client WebSocket客户端
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan *TaskUpdateMessage
}

// TaskUpdateMessage 任务更新消息
type TaskUpdateMessage struct {
	TaskID   string            `json:"task_id"`
	Status   model.TaskStatus  `json:"status"`
	Progress int               `json:"progress,omitempty"`
	Result   interface{}       `json:"result,omitempty"`
	Error    string            `json:"error,omitempty"`
	Timestamp time.Time        `json:"timestamp"`
}

// NewHub 创建新的WebSocket中心
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan *TaskUpdateMessage, 100),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// Run 运行WebSocket中心
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			logger.Info("Client registered", zap.Int("total_clients", len(h.clients)))

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				logger.Info("Client unregistered", zap.Int("total_clients", len(h.clients)))
			}

		case message := <-h.broadcast:
			// 广播消息给所有客户端
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					// 如果发送通道满了，关闭连接
					close(client.send)
					delete(h.clients, client)
					logger.Warn("Client send channel full, disconnecting client")
				}
			}
		}
	}
}

// Register returns the register channel
func (h *Hub) Register() chan<- *Client {
	return h.register
}

// Unregister returns the unregister channel
func (h *Hub) Unregister() chan<- *Client {
	return h.unregister
}

// Broadcast returns the broadcast channel
func (h *Hub) Broadcast() chan<- *TaskUpdateMessage {
	return h.broadcast
}

// NewClient 创建新的WebSocket客户端
func (h *Hub) NewClient(conn *websocket.Conn) *Client {
	return &Client{
		hub:  h,
		conn: conn,
		send: make(chan *TaskUpdateMessage, 256),
	}
}

// Send returns the send channel
func (c *Client) Send() chan<- *TaskUpdateMessage {
	return c.send
}
func (c *Client) ServeWS() {
	// 启动写协程
	go c.writePump()

	// 启动读协程
	c.readPump()
}

// writePump 写协程
func (c *Client) writePump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				// 通道已关闭
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// 发送消息
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				logger.Error("Failed to get next writer", zap.Error(err))
				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				logger.Error("Failed to marshal message", zap.Error(err))
				return
			}

			w.Write(data)
			w.Close()

		case <-time.After(54 * time.Second):
			// 发送ping消息保持连接
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				logger.Error("Failed to send ping message", zap.Error(err))
				return
			}
		}
	}
}

// readPump 读协程
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Error("WebSocket unexpected close error", zap.Error(err))
			}
			break
		}
	}
}

// BroadcastTaskUpdate 广播任务更新
func (h *Hub) BroadcastTaskUpdate(task *model.Task) {
	message := &TaskUpdateMessage{
		TaskID:    task.ID,
		Status:    task.Status,
		Progress:  task.Progress,
		Result:    task.Result,
		Error:     task.Error,
		Timestamp: time.Now(),
	}

	select {
	case h.broadcast <- message:
		logger.Debug("Task update broadcasted", zap.String("task_id", task.ID))
	default:
		logger.Warn("Broadcast channel full, dropping message", zap.String("task_id", task.ID))
	}
}

// GetClientCount 获取客户端数量
func (h *Hub) GetClientCount() int {
	return len(h.clients)
}