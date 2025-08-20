package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"signaling-server-webrtc/pkg"
	"signaling-server-webrtc/srv"
	"signaling-server-webrtc/utils"
)

func HandleHealthCheck(serviceName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]string{
			"serviceName": serviceName,
			"status":      "ok",
			"serverTime":  time.Now().Format("2006-01-02 15:04:05"),
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(resp)
	}
}

func HandleCreateRoom(hub *pkg.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		room, _ := srv.CreateRoom(hub)

		utils.WriteJSON(w, http.StatusOK, room)
	}
}

func HandleJoinRoom(hub *pkg.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomId := r.URL.Query().Get("roomId")

		if roomId == "" {
			utils.WriteError(w, http.StatusBadRequest, "invalid room id!")
		}
		room, err := srv.JoinRoom(hub, roomId)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err.Error())
		}

		utils.WriteJSON(w, http.StatusOK, room)
	}
}

// func HandleLeaveRoom(hub *pkg.Hub) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		room, err := utils.DecodeRoomRequest(r)
// 		if err != nil {
// 			utils.WriteError(w, http.StatusBadRequest, "Invalid payload")
// 			return
// 		}

// 		err = room.ValidateLeaveRoom()
// 		if err != nil {
// 			http.Error(w, "Invalid request payload", http.StatusBadRequest)
// 			return
// 		}

// 		client := hub.GetClientFromRoom(*room.RoomId, *room.ClientId)
// 		if client == nil {
// 			utils.WriteError(w, http.StatusNotFound, "Client not found in room")
// 			return
// 		}

// 		leftRoom, err := srv.LeaveRoom(hub, room, client)
// 		if err != nil {
// 			http.Error(w, "Some Error Occured", http.StatusInternalServerError)
// 			return
// 		}

// 		utils.WriteJSON(w, http.StatusOK, leftRoom)
// 	}
// }

func HandleRoomStats(hub *pkg.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomId := r.URL.Query().Get("roomId")
		if roomId == "" {
			utils.WriteJSON(w, http.StatusOK, hub.HubStats())
			return
		}
		utils.WriteJSON(w, http.StatusOK, hub.RoomStats(roomId))
	}
}
