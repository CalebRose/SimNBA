package ws

import (
	"log"
	"sync"

	"github.com/CalebRose/SimNBA/managers"
	"github.com/CalebRose/SimNBA/structs"
)

// latestTS caches the most recently broadcast timestamp so new connections
// don't each hit the database.
var (
	latestTS    structs.Timestamp
	hasLatestTS bool
	latestMu    sync.RWMutex
)

// currentTimestamp returns the cached timestamp, falling back to the
// database until the first broadcast has populated the cache.
func currentTimestamp() structs.Timestamp {
	latestMu.RLock()
	ts, ok := latestTS, hasLatestTS
	latestMu.RUnlock()
	if ok {
		return ts
	}
	return managers.GetTimestamp()
}

// BroadcastTSUpdate sends the updated timestamp to all connected WebSocket clients
func BroadcastTSUpdate(ts structs.Timestamp) error {
	latestMu.Lock()
	latestTS = ts
	hasLatestTS = true
	latestMu.Unlock()

	mu.Lock()
	defer mu.Unlock()
	for c := range clients {
		select {
		case c.send <- ts:
		default:
			// Client is too slow; drop it. Closing send stops its writePump,
			// which closes the connection and unblocks readPump.
			log.Println("Dropping slow WebSocket client")
			delete(clients, c)
			close(c.send)
		}
	}
	return nil
}
