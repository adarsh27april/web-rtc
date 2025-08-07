package types

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	Connection *websocket.Conn // this is a websocket connection
	Send       chan []byte     // this is a channel to send messages to the client
	RoomID     string
	// Hub        hub.Hub
	ClientID string
}

type MessageEnvelope struct {
	Sender *Client
	RoomID string
	Data   []byte
}
