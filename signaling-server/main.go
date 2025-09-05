package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"signaling-server-webrtc/pkg"
	"signaling-server-webrtc/pkg/handlers"
	"signaling-server-webrtc/srv"
	"signaling-server-webrtc/utils"
)

func init() {
	// Seed the math/rand package with the current time
	rand.Seed(time.Now().UnixNano())
}

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
	var allowedOrigins []string
	originsStr := os.Getenv("CORS_ALLOWED_ORIGINS")

	if originsStr == "" {
		log.Fatalf("FATAL: 'CORS_ALLOWED_ORIGINS' environment variable not set")
	} else {
		allowedOrigins = strings.Split(originsStr, ",")
	}
	c := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "Accept"},
		AllowCredentials: true,
		Debug:            true, // Enable for debugging CORS issues
	})
	handler := c.Handler(r)

	// handling server start and shutdown
	var server *http.Server
	{
		port := utils.GetEnv("PORT")
		if port == "" {
			port = "1337" // Default for local
		}
		server = &http.Server{
			Addr:    ":" + port,
			Handler: handler,
		}

		go func() {
			log.Printf("Signaling server started, PORT: %v\n", server.Addr)

			var err error
			switch utils.GetEnv("ENV") {
			case "local":
				err = server.ListenAndServeTLS("cert.pem", "key.pem")

			case "prod":
				err = server.ListenAndServe()

			default:
				log.Fatalf("FATAL: 'ENV' environment variable not set or invalid. Must be 'local' or 'prod'.")
			}

			if err != nil && err != http.ErrServerClosed {
				log.Fatalf("Server failed to start or unexpectedly closed: %s\n", err)
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
