package api

import (
	"net/http"
	"time"

	"github.com/yourname/02-concurrency-worker/internal/api/dto"
	"github.com/yourname/02-concurrency-worker/internal/pkg/websocket"
	"github.com/gin-gonic/gin"
	gorillaWs "github.com/gorilla/websocket"
	"go.uber.org/zap"

	"github.com/yourname/02-concurrency-worker/internal/pkg/logger"
)

var upgrader = gorillaWs.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// 允许所有来源（在生产环境中应该更严格）
		return true
	},
}

// WebSocketHandler WebSocket处理器
type WebSocketHandler struct {
	hub *websocket.Hub
}

// NewWebSocketHandler 创建新的WebSocket处理器
func NewWebSocketHandler(hub *websocket.Hub) *WebSocketHandler {
	return &WebSocketHandler{
		hub: hub,
	}
}

// HandleWebSocket 处理WebSocket连接
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	// 升级HTTP连接到WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error("Failed to upgrade connection to WebSocket", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to establish WebSocket connection"})
		return
	}
	defer conn.Close()

	// 创建新的客户端
	client := h.hub.NewClient(conn)

	// 注册客户端
	h.hub.Register() <- client

	// 发送连接成功的消息
	successMsg := &dto.WebSocketResponse{
		Type:    "connection",
		Message: "Connected to task worker WebSocket",
		Data:    gin.H{"client_count": h.hub.GetClientCount()},
	}

	client.Send() <- &websocket.TaskUpdateMessage{
		TaskID:    "system",
		Status:    "connected",
		Result:    successMsg,
		Timestamp: time.Now(),
	}

	// 服务WebSocket连接
	client.ServeWS()
}