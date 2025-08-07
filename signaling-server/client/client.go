package client

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"signaling-server-webrtc/utils"
)

type Hub interface {
	Register(*Client)
	Unregister(*Client)
	Broadcast(MessageEnvelope)
}

type Client struct {
	Connection *websocket.Conn // this is a websocket connection
	Send       chan []byte     // this is a channel to send messages to the client
	RoomID     string
	Hub        Hub
	ClientID   string
}

type MessageEnvelope struct {
	Sender *Client
	RoomID string
	Data   []byte
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func ServeWs(h Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Println("Read message failed:", err)
		return
	}

	var join utils.JoinRoomMessage
	json.Unmarshal(msg, &join)

	c := &Client{
		Connection: conn,
		Send:       make(chan []byte, 256),
		RoomID:     join.RoomID,
		Hub:        h,
		ClientID:   uuid.New().String(),
	}

	h.Register(c)

	go c.readPump()
	go c.writePump()
}

func (c *Client) readPump() {
	defer func() {
		utils.LogRoom(c.RoomID, c.ClientID, "🔌 Disconnected")
		c.Hub.Unregister(c)
		c.Connection.Close()
	}()

	for {
		_, message, err := c.Connection.ReadMessage()
		if err != nil {
			break
		}
		utils.LogRoom(c.RoomID, c.ClientID, "📥 Received message from client")

		c.Hub.Broadcast(MessageEnvelope{
			Sender: c,
			RoomID: c.RoomID,
			Data:   message,
		})
		utils.LogRoom(c.RoomID, c.ClientID, "Broadcasting message to other clients")
	}
}

func (c *Client) writePump() {
	for msg := range c.Send {
		utils.LogRoom(c.RoomID, c.ClientID, "📤 Sending message to client")
		c.Connection.WriteMessage(websocket.TextMessage, msg)
	}
}
