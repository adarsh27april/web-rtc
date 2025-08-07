package srv

import (
	"log"
	"net/http"
	pkg "signaling-server-webrtc/pkg"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for dev
	},
}

func ServeWS(hub *pkg.Hub, w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("roomId")
	clientID := r.URL.Query().Get("clientId")

	if roomID == "" || clientID == "" {
		http.Error(w, "Missing roomId or clientId", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	client := &pkg.Client{
		Connection: conn,
		ClientID:   clientID,
		RoomID:     roomID,
		Send:       make(chan []byte, 256),
	}

	hub.Register <- client // register the client

	// these are Per-client goroutines
	go client.WritePump()
	go client.ReadPump(hub)
}
