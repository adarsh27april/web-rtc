package types

import (
	"sync"
)

/*
Hub.Rooms is like

	"roomId": {
		"client1_ptr": true    // all clients of a room are in map
		"client2_ptr": true    // instead of array for easy/fast access
		"client3_ptr": true
	}
*/
type Hub struct {
	Rooms      map[string]map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan MessageEnvelope
	Mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan MessageEnvelope),
	}
}

// it will traverse the Hub to get client ptr from clientId if it is present in roomId
func (h *Hub) GetClientFromRoom(roomID, clientID string) *Client {
	h.Mu.RLock()
	defer h.Mu.RUnlock()
	for c := range h.Rooms[roomID] {
		if c.ClientID == clientID {
			return c
		}
	}
	return nil
}

// func (h *Hub) Run() {
// 	for {
// 		select {
// 		case c := <-h.Register: // get value from Register channel
// 			h.addClient(c) // add client to the hub
// 		case c := <-h.Unregister: // get value from Unregister channel
// 			h.removeClient(c) // remove client from the hub
// 		case msg := <-h.Broadcast: // get value from Broadcast channel
// 			h.sendToRoom(msg) // send message to the room
// 		}
// 	}
// }

// func (h *Hub) addClient(c *Client) {
// 	h.Mu.Lock() // this will make a write lock on h.Mu value to make a write on the HUB Rooms map
// 	defer h.Mu.Unlock()
// 	if h.Rooms[c.RoomID] == nil {
// 		h.Rooms[c.RoomID] = make(map[*Client]bool)
// 	}
// 	h.Rooms[c.RoomID][c] = true
// 	fmt.Println(c.RoomID, c.ClientID, "✅ Joined room")
// }

// func (h *Hub) removeClient(c *Client) {
// 	h.Mu.Lock()
// 	defer h.Mu.Unlock()
// 	if _, ok := h.Rooms[c.RoomID][c]; ok {
// 		delete(h.Rooms[c.RoomID], c)
// 		close(c.Send)
// 	}
// 	utils.LogRoom(c.RoomID, c.ClientID, "❌ Left room")
// }

// func (h *Hub) sendToRoom(msg MessageEnvelope) {
// 	h.Mu.RLock()
// 	defer h.Mu.RUnlock()
// 	for c := range h.Rooms[msg.RoomID] {
// 		if c != msg.Sender {
// 			c.Send <- msg.Data
// 		}
// 	}
// 	utils.LogRoom(msg.RoomID, msg.Sender.ClientID, "📡 Relaying message to other clients in room")
// }

// func (h *Hub) RegisterClient(c *Client) {
// 	h.Register <- c
// }

// func (h *Hub) UnregisterClient(c *Client) {
// 	h.Unregister <- c
// }

// func (h *Hub) BroadcastToClient(msg MessageEnvelope) {
// 	h.Broadcast <- msg
// }
