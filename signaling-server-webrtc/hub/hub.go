package hub

import (
	"signaling-server-webrtc/client"
	"sync"
)

// this hub renotes a room where lients will be added and removed by using go routines and the.
type Hub struct {
	rooms      map[string]map[*client.Client]bool
	register   chan *client.Client
	unregister chan *client.Client
	broadcast  chan client.MessageEnvelope
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[string]map[*client.Client]bool),
		register:   make(chan *client.Client),
		unregister: make(chan *client.Client),
		broadcast:  make(chan client.MessageEnvelope),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register: // get value from register channel
			h.addClient(c) // add client to the hub
		case c := <-h.unregister: // get value from unregister channel
			h.removeClient(c) // remove client from the hub
		case msg := <-h.broadcast: // get value from broadcast channel
			h.sendToRoom(msg) // send message to the room
		}
	}
}

func (h *Hub) addClient(c *client.Client) {
	h.mu.Lock() // this will make a write lock on h.mu value to make a write on the HUB struct
	defer h.mu.Unlock()
	if h.rooms[c.RoomID] == nil {
		h.rooms[c.RoomID] = make(map[*client.Client]bool)
	}
	h.rooms[c.RoomID][c] = true
}

func (h *Hub) removeClient(c *client.Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.rooms[c.RoomID][c]; ok {
		delete(h.rooms[c.RoomID], c)
		close(c.Send)
	}
}

func (h *Hub) sendToRoom(msg client.MessageEnvelope) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for c := range h.rooms[msg.RoomID] {
		if c != msg.Sender {
			c.Send <- msg.Data
		}
	}
}

func (h *Hub) Register(c *client.Client) {
	h.register <- c
}

func (h *Hub) Unregister(c *client.Client) {
	h.unregister <- c
}

func (h *Hub) Broadcast(msg client.MessageEnvelope) {
	h.broadcast <- msg
}
