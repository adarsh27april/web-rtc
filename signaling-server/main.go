package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"signaling-server-webrtc/pkg"
	"signaling-server-webrtc/pkg/handlers"
	"signaling-server-webrtc/srv"
	"signaling-server-webrtc/utils"
)

func main() {
	// this hub denotes a room where clients will be added and removed by using go routines.
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
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:3000", "http://localhost:4173"}, // Your frontend origin
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "Accept"},
		AllowCredentials: true,
		Debug:            true, // Enable for debugging CORS issues
	})
	handler := c.Handler(r)

	// handling server start and shutdown
	var server *http.Server
	{
		server = &http.Server{
			Addr:    ":" + utils.GetEnv("PORT"),
			Handler: handler,
		}

		go func() {
			log.Printf("Signaling server started, PORT%v\n", server.Addr)
			err := server.ListenAndServe()

			if err != nil && err != http.ErrServerClosed {
				log.Fatalf("listen: %s\n", err)
			}
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		<-quit
		log.Println("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Fatal("Server forced to shutdown:", err)
		}
	}
}
