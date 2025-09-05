package srv

import (
	"fmt"
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
	room, roomExists := hub.Rooms[roomID]
	_, clientExistsInRoom := room[clientId]
	hub.Mu.RUnlock()

	if !roomExists || !clientExistsInRoom {
		http.Error(w, "Unauthorized: Invalid room or client ID", http.StatusUnauthorized)
		return
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
	const WAIT_TIME = 1 * time.Minute
	go func(roomID, clientId string) {
		time.Sleep(WAIT_TIME)

		clientPtr := hub.GetClientFromRoom(roomID, clientId)
		if clientPtr == nil {
			log.Printf("[Room:%s] [Client:%s] Client no longer exists, skipping timeout", roomID, clientId)
			return
		}

		stats := hub.RoomStats(roomID)
		if len(stats.Clients) < 2 {
			timeOutMsg := []byte(fmt.Sprintf(`{
				"type":"timeout", 
				"message": "no peer joined in %d seconds"
			}`, int(WAIT_TIME.Seconds())))

			select {
			case clientPtr.Send <- timeOutMsg:
				log.Printf("[Room:%s] [Client:%s] No peer joined in time, notifying client", roomID, clientId)
			default:
				log.Printf("[Room:%s] [Client:%s] Cannot send timeout - channel unavailable", roomID, clientId)
			}
		}
	}(roomID, clientId)
}
