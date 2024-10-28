package ws

import (
	"log"

	"github.com/CalebRose/SimNBA/structs"
)

// BroadcastTSUpdate sends the updated timestamp to all connected WebSocket clients
func BroadcastTSUpdate(ts structs.Timestamp) {
	for conn := range clients {
		err := conn.WriteJSON(ts)
		if err != nil {
			log.Println("Error broadcasting to WebSocket client:", err)
			conn.Close()
			delete(clients, conn) // Remove client on error
		}
	}
}
