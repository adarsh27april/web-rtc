package main

import (
	"log"
	"net/http"

	"signaling-server-webrtc/client"
	"signaling-server-webrtc/hub"
)

func main() {
	// this hub renotes a room where lients will be added and removed by using go routines.
	h := hub.NewHub() // this will create 3 new channels for register, unregister, broadcast
	go h.Run()        // this is going to run concurrently and listen to all the data made available in that channel

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		client.ServeWs(h, w, r)
	})
	// The server listens on port 4040 and handles all HTTP methods (GET, POST, PUT, etc.) at the '/ws' route.
	// Any request to localhost:4040/ws will trigger the callback function to establish a WebSocket connection.

	log.Println("Signaling server started on :4040")
	log.Fatal(http.ListenAndServe(":4040", nil))
}
