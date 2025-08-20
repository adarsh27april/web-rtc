package srv

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"signaling-server-webrtc/pkg"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for dev
	},
}

func ServeWS(hub *pkg.Hub, w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("roomId")
	clientId := r.URL.Query().Get("clientId")

	if roomID == "" || clientId == "" {
		http.Error(w, "Missing roomId or clientId", http.StatusBadRequest)
		return
	}

	hub.Mu.RLock()
	_, exist := hub.Rooms[roomID][clientId]
	hub.Mu.RUnlock()

	if !exist {
		http.Error(w, "Unauthorized Client", http.StatusBadRequest)
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	client := &pkg.Client{
		Connection: conn,
		ClientId:   clientId,
		RoomID:     roomID,
		Send:       make(chan []byte, 256),
	}

	hub.Register <- client // register the client

	// these are Per-client goroutines
	go client.WritePump()
	go client.ReadPump(hub)

	// notify first peer if nobody joins in time
	const WAIT_TIME = 2 * time.Minute
	go func(roomID, clientID string, c *pkg.Client) {
		time.Sleep(WAIT_TIME)

		stats := hub.RoomStats(roomID)
		if len(stats.Clients) < 2 {
			// non-fatal: just inform; your UI can show a toast/option to keep waiting
			c.Send <- []byte(`{"type":"timeout","data":{"reason":"no-peer-joined","afterSec":60}}`)
		}
	}(roomID, clientId, client)
}
