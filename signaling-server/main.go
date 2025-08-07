package main

import (
	"log"
	"net/http"

	// "signaling-server-webrtc/client"
	// "signaling-server-webrtc/hub"

	"signaling-server-webrtc/handlers"
	"signaling-server-webrtc/pkg/types"

	"github.com/gorilla/mux"
	// "signaling-server-webrtc/handlers"
	// "signaling-server-webrtc/hub"
)

func main() {
	// this hub renotes a room where lients will be added and removed by using go routines.

	h := types.NewHub() // this will create 3 new channels for register, unregister, broadcast
	// go h.Run()          // this is going to run concurrently and listen to all the data made available in that channel

	r := mux.NewRouter()

	r.HandleFunc("/api/health", handlers.HandlerHealthCheck("Signaling Server")).Methods("GET")
	r.HandleFunc("/api/rooms/join", handlers.HandlerJoinRoom(h)).Methods("POST")
	r.HandleFunc("/api/rooms/leave", handlers.HandleLeaveRoom(h)).Methods("POST")
	r.HandleFunc("/api/rooms/stats", handlers.HandleRoomStats(h)).Methods("GET")

	// http.HandleFunc("/api/health", handlers.HandlerHealthCheck("Signaling Server"))

	// http.HandleFunc("/api/rooms/join", handlers.HandlerJoinRoom(h))

	// http.HandleFunc("/api/rooms/leave", handlers.HandleLeaveRoom(h))

	// http.HandleFunc("api/rooms/stats", handlers.HandleRoomStats(h))

	// http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
	// 	client.ServeWs(h, w, r)
	// })
	// The server listens on port 4040 and handles all HTTP methods (GET, POST, PUT, etc.) at the '/ws' route.
	// Any request to localhost:4040/ws will trigger the callback function to establish a WebSocket connection.

	log.Println("Signaling server started on :1337")
	log.Fatal(http.ListenAndServe(":1337", r))
}
