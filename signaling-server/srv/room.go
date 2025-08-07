package srv

import (
	"fmt"

	pkg "signaling-server-webrtc/pkg"
	"signaling-server-webrtc/pkg/types"
	"signaling-server-webrtc/utils"
)

// Logic to handle joining a room
func JoinRoom(hub *pkg.Hub, room *types.Room, client *pkg.Client) types.Room {
	hub.Mu.Lock()
	defer hub.Mu.Unlock()

	// if no room from client then create new room for it.
	if room.RoomId == nil {
		room.RoomId = utils.Ptr(utils.GenerateShortID())
	}

	// create room it it not exists
	if _, ok := hub.Rooms[*room.RoomId]; !ok {
		hub.Rooms[*room.RoomId] = make(map[*pkg.Client]bool)
	}

	// set roomId for client
	client.RoomID = *room.RoomId

	/*
		We will not add the client in room here just return a newly created {roomId, clientId}.
		It will be added in the websocket connection
	*/

	// add client to room. a single room is `hub.Rooms[*room.RoomId]`
	// hub.Rooms[*room.RoomId][client] = true

	utils.LogRoom(*room.RoomId, client.ClientID, "Client joined room")

	return types.Room{
		RoomId:   room.RoomId,
		ClientID: &client.ClientID,
		Status:   utils.Ptr("joined"),
	}
}

/*
This will be deleted. The leave room will be called in websocket connection only.
*/
// Logic to handle leaving a room
func LeaveRoom(hub *pkg.Hub, room types.Room, client *pkg.Client) (types.Room, error) {
	hub.Mu.Lock()
	defer hub.Mu.Unlock()

	res := types.Room{
		RoomId:   room.RoomId,
		ClientID: &client.ClientID,
		Status:   utils.Ptr("left"),
	}

	clients, ok := hub.Rooms[*room.RoomId] // check if room id exist or not in hub. if yes then gather all clients fo that room
	if !ok {
		return types.Room{}, fmt.Errorf("room with ID %s does not exist", *room.RoomId)
	}
	if _, exists := clients[client]; !exists { // check if client exist or not in room
		return types.Room{}, fmt.Errorf("client %s not found in room %s", client.ClientID, *room.RoomId)
	}

	delete(clients, client) // delete client from room

	utils.LogRoom(*room.RoomId, client.ClientID, "Client left room")

	if len(clients) == 0 { // clean empty room
		delete(hub.Rooms, *room.RoomId)
		utils.LogRoom(*room.RoomId, client.ClientID, "Room is empty. Deleted.")
	}

	return res, nil
}

func HubStats(hub *pkg.Hub) types.HubStats {
	stats := types.HubStats{}
	for roomID, clients := range hub.Rooms {
		clientList := []string{}
		for client := range clients {
			clientList = append(clientList, client.ClientID)
		}

		stats.Rooms = append(stats.Rooms, types.RoomStats{
			RoomID:  roomID,
			Clients: clientList,
		})
	}

	return stats
}

func RoomStats(hub *pkg.Hub, roomId string) types.RoomStats {
	stats := types.RoomStats{}
	for rId, clients := range hub.Rooms {
		if rId == roomId {
			stats.RoomID = rId
			for client := range clients {
				stats.Clients = append(stats.Clients, client.ClientID)
			}
			break
		}
	}
	return stats
}
