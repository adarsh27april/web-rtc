package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	pkg "signaling-server-webrtc/pkg"
	"signaling-server-webrtc/pkg/handlers"
	"signaling-server-webrtc/srv"
)

const PORT = ":1337"

func main() {
	// this hub renotes a room where lients will be added and removed by using go routines.

	h := pkg.NewHub() // this will create 3 new channels for register, unregister, broadcast
	go h.Run()        // this is going to run concurrently and listen to all the data made available in that channel

	r := mux.NewRouter()

	r.HandleFunc("/api/health", handlers.HandleHealthCheck("Signaling Server")).Methods("GET")

	r.HandleFunc("/api/rooms/create", handlers.HandleCreateRoom(h)).Methods("POST")
	r.HandleFunc("/api/rooms/join", handlers.HandleJoinRoom(h)).Methods("POST")
	// r.HandleFunc("/api/rooms/leave", handlers.HandleLeaveRoom(h)).Methods("POST")
	r.HandleFunc("/api/rooms/stats", handlers.HandleRoomStats(h)).Methods("GET")

	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		srv.ServeWS(h, w, r)
	}).Methods("GET")
	// The '/ws' route listens for WebSocket upgrade requests over HTTP GET.
	// Clients connect to this endpoint to establish a persistent WebSocket connection.

	// Configure CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Your frontend origin
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
		Debug:            true, // Enable for debugging CORS issues
	})
	handler := c.Handler(r)

	log.Println("Signaling server started on ", PORT)
	log.Fatal(http.ListenAndServe(PORT, handler))
}
