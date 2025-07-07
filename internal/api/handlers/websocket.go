package handlers

import (
	"net/http"

	"block-flow/internal/engine"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for now
	},
}

// WebSocketHandler handles WebSocket connections
type WebSocketHandler struct {
	engine *engine.Engine
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler(engine *engine.Engine) *WebSocketHandler {
	return &WebSocketHandler{
		engine: engine,
	}
}

// HandleWebSocket handles WebSocket connections for real-time updates
func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade connection", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	// Handle WebSocket communication
	for {
		// Read message from client
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		// Echo message back (placeholder implementation)
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			break
		}
	}
}
