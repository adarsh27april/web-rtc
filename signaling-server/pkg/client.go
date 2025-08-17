package pkg

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type Client struct {
	Connection *websocket.Conn // this is a websocket connection
	Send       chan []byte     // this is a channel to send messages to the client
	RoomID     string
	// Hub        hub.Hub
	ClientId string
}

type MessageEnvelope struct {
	Sender *Client
	RoomID string
	Data   []byte
}

func (c *Client) ReadPump(hub *Hub) {
	// this go routine func should run endlessly
	defer func() {
		// this defer func will only be called if the code breaks due to error
		hub.Unregister <- c // sending c to channel Unregister
		c.Connection.Close()
	}()

	for {
		_, message, err := c.Connection.ReadMessage()
		if err != nil {
			break // Client disconnected -> it will Unregister
		}

		hub.Broadcast <- MessageEnvelope{
			Sender: c,
			RoomID: c.RoomID,
			Data:   message,
		}
	}
}

func (c *Client) WritePump() {
	defer c.Connection.Close()

	// here Send is a channel so if it ends then the loop will wait for new value to appear here.
	// if the channle is closed then only the loop ends.
	for message := range c.Send {
		fmt.Println("new message: ", string(message))
		err := c.Connection.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break // Write failed (disconnected or closed)
		}
	}
}
