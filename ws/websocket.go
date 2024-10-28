package ws

import (
	"log"
	"net/http"

	"github.com/CalebRose/SimNBA/managers"
	"github.com/gorilla/websocket"
)

// Upgrader configures the WebSocket connection
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// In production, you should validate origins for security
		return true
	},
}

// Global map to keep track of connected WebSocket clients
var clients = make(map[*websocket.Conn]bool)

// WebSocketHandler handles WebSocket connection requests
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP request to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	// Add new WebSocket connection to the clients map
	clients[conn] = true
	log.Println("New WebSocket client connected")

	ts := managers.GetTimestamp()
	// Send the latest timestamp to the user immediately upon connection
	err = conn.WriteJSON(ts)
	if err != nil {
		log.Println("Error sending timestamp:", err)
		return
	}

	// Listen for messages (e.g., handle disconnects)
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket client disconnected:", err)
			conn.Close()
			delete(clients, conn)
			break
		}

		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Println("Error writing WebSocket message:", err)
			break
		}
	}
}
