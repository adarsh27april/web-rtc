package pkg

import (
	"fmt"
	"sort"
	"sync"

	"signaling-server-webrtc/pkg/types"
	"signaling-server-webrtc/utils"
)

/*
Hub is the Central Event Manager. Responsible for:

  - Managing active rooms and the clients inside them
  - Relaying messages between clients in the same room
  - Cleaning up connections when clients disconnect
  - Acting as a message router using Go channels for concurrency safety

Hub.Rooms structure:

	{
		"roomId1": {
			"ClientId1": *Client, // ref. to the connected client struct
			"ClientId2": *Client, // each ClientId maps to its Client instance
			"ClientId3": *Client
		}
		...
	}

Notes:
  - A room is represented as a map of ClientId → *Client.
  - Empty rooms are removed to free up memory.
  - A Client may be pre-registered (with nil connection) via REST before WebSocket connects.
*/
type Hub struct {
	Rooms      map[string]map[string]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan MessageEnvelope
	Mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan MessageEnvelope),
	}
}

// it will traverse the Hub to get client ptr from ClientId if it is present in roomId
func (h *Hub) GetClientFromRoom(roomID, clientId string) *Client {
	h.Mu.RLock()
	defer h.Mu.RUnlock()

	if room, ok := h.Rooms[roomID]; ok {
		if client, ok := room[clientId]; ok {
			return client
		}
	}
	return nil
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.Register: // get value(client) from Register channel
			h.addClient(c) // add client to the hub
			h.assignClientRole(c.RoomID)
		case c := <-h.Unregister: // get value from Unregister channel
			h.removeClient(c) // remove client from the hub
		case msg := <-h.Broadcast: // get value from Broadcast channel
			h.sendToRoom(msg) // send message to the room
		}
	}
}

func (h *Hub) addClient(c *Client) {
	h.Mu.Lock()
	defer h.Mu.Unlock()

	if h.Rooms[c.RoomID] == nil {
		h.Rooms[c.RoomID] = make(map[string]*Client)
	}
	h.Rooms[c.RoomID][c.ClientId] = c
	fmt.Println(c.RoomID, c.ClientId, "✅ Joined room")
}

func (h *Hub) removeClient(c *Client) {
	h.Mu.Lock()
	defer h.Mu.Unlock()

	if _, ok := h.Rooms[c.RoomID][c.ClientId]; ok {
		delete(h.Rooms[c.RoomID], c.ClientId)
		close(c.Send)
	}
	utils.LogRoom(c.RoomID, c.ClientId, "❌ Left room")

	// Clean up room if empty
	if len(h.Rooms[c.RoomID]) == 0 {
		utils.LogRoom(c.RoomID, "Nil", "empty room! Deleting... 🗑️")
		delete(h.Rooms, c.RoomID)
	}
}

func (h *Hub) sendToRoom(msg MessageEnvelope) {
	h.Mu.RLock()
	defer h.Mu.RUnlock()

	for _, c := range h.Rooms[msg.RoomID] { // _ is ClientId
		if c != msg.Sender {
			c.Send <- msg.Data
		}
	}
	utils.LogRoom(msg.RoomID, msg.Sender.ClientId, "📡 Relaying message to other clients in room")
}

func (h *Hub) assignClientRole(roomId string) {
	stats := h.RoomStats(roomId)
	if len(stats.Clients) == 2 {
		cIds := stats.Clients
		sort.Strings(cIds)
		offererClient, answererClient := cIds[0], cIds[1]

		h.Mu.RLock()
		defer h.Mu.RUnlock()
		for clientId, clientPtr := range h.Rooms[roomId] {

			switch clientId {
			case offererClient:
				clientPtr.Send <- []byte(`{"type":"role","data":{"role":"offerer"}}`)
			case answererClient:
				clientPtr.Send <- []byte(`{"type":"role","data":{"role":"answerer"}}`)
			}
		}
	}
}

func (hub *Hub) HubStats() types.HubStats {
	hub.Mu.RLock()
	defer hub.Mu.RUnlock()

	stats := types.HubStats{}
	for roomID, clientsMap := range hub.Rooms {
		roomStats := types.RoomStats{
			RoomID: roomID,
		}
		for clientId := range clientsMap {
			roomStats.Clients = append(roomStats.Clients, clientId)
		}

		stats.Rooms = append(stats.Rooms, roomStats)
	}
	stats.TotalRooms = len(stats.Rooms)

	return stats
}

func (hub *Hub) RoomStats(roomId string) types.RoomStats {
	hub.Mu.RLock()
	defer hub.Mu.RUnlock()

	roomData := hub.Rooms[roomId]
	roomStats := types.RoomStats{
		RoomID: roomId,
	}
	for clientIds := range roomData {
		roomStats.Clients = append(roomStats.Clients, clientIds)
	}

	return roomStats
}
