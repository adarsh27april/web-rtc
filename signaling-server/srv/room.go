package srv

import (
	"fmt"

	"signaling-server-webrtc/pkg"
	"signaling-server-webrtc/pkg/types"
	"signaling-server-webrtc/utils"
)

// only client A can create a room
func CreateRoom(hub *pkg.Hub) (types.Room, error) {
	roomId := utils.GenerateShortID()
	clientId := utils.GenerateShortID()

	hub.Mu.Lock()
	defer hub.Mu.Unlock()

	if hub.Rooms[roomId] == nil {
		hub.Rooms[roomId] = make(map[string]*pkg.Client)
	}

	// adding placeholder client until WS Connects
	hub.Rooms[roomId][clientId] = nil
	// actual client object will be formed when WS connection is made to connect

	utils.LogRoom(roomId, clientId, "room created placeholder client attached")

	return types.Room{
		RoomId:   &roomId,
		ClientId: &clientId,
		Status:   utils.Ptr("created"),
	}, nil
}

// client B,C,... will join the room created by client A
func JoinRoom(hub *pkg.Hub, roomId string) (types.Room, error) {
	hub.Mu.Lock()
	defer hub.Mu.Unlock()

	if _, exists := hub.Rooms[roomId]; !exists {
		return types.Room{}, fmt.Errorf("invalid room id! room doesn't exist")
	}

	clientId := utils.GenerateShortID()

	// if room exist add a placeholder value for client in a room.
	// will be updated in WS connection
	hub.Rooms[roomId][clientId] = nil

	utils.LogRoom(roomId, clientId, "Client joined room (placeholder)")

	return types.Room{
		RoomId:   &roomId,
		ClientId: &clientId,
		Status:   utils.Ptr("pending"),
	}, nil
}

/*
This will be deleted. The leave room will be called in websocket connection only.
*/
// Logic to handle leaving a room
// func LeaveRoom(hub *pkg.Hub, room types.Room, client *pkg.Client) (types.Room, error) {
// 	hub.Mu.Lock()
// 	defer hub.Mu.Unlock()

// 	res := types.Room{
// 		RoomId:   room.RoomId,
// 		ClientId: &client.ClientId,
// 		Status:   utils.Ptr("left"),
// 	}

// 	clients, ok := hub.Rooms[*room.RoomId] // check if room id exist or not in hub. if yes then gather all clients fo that room
// 	if !ok {
// 		return types.Room{}, fmt.Errorf("room with ID %s does not exist", *room.RoomId)
// 	}
// 	if _, exists := clients[client]; !exists { // check if client exist or not in room
// 		return types.Room{}, fmt.Errorf("client %s not found in room %s", client.ClientId, *room.RoomId)
// 	}

// 	delete(clients, client) // delete client from room

// 	utils.LogRoom(*room.RoomId, client.ClientId, "Client left room")

// 	if len(clients) == 0 { // clean empty room
// 		delete(hub.Rooms, *room.RoomId)
// 		utils.LogRoom(*room.RoomId, client.ClientId, "Room is empty. Deleted.")
// 	}

// 	return res, nil
// }
