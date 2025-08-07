package handlers

import (
	"encoding/json"
	"net/http"
	"signaling-server-webrtc/pkg/types"
	"signaling-server-webrtc/srv"
	"signaling-server-webrtc/utils"
	"time"
)

func HandlerHealthCheck(serviceName string) http.HandlerFunc {
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

func HandlerJoinRoom(hub *types.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		joinRoom, err := utils.DecodeRoomRequest(r)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "Invalid payload")
			return
		}

		// creating random client id and adding it to the room
		client := &types.Client{
			ClientID: utils.GenerateShortID(), // or UUID
		}

		room := srv.JoinRoom(hub, &joinRoom, client)
		// client will be added to a new room or to room.RoomId

		client.RoomID = *room.RoomId

		utils.WriteJSON(w, http.StatusOK, room)
	}
}

func HandleLeaveRoom(hub *types.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		room, err := utils.DecodeRoomRequest(r)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "Invalid payload")
			return
		}

		err = room.ValidateLeaveRoom()
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		client := hub.GetClientFromRoom(*room.RoomId, *room.ClientID)
		if client == nil {
			utils.WriteError(w, http.StatusNotFound, "Client not found in room")
			return
		}

		leftRoom, err := srv.LeaveRoom(hub, room, client)
		if err != nil {
			http.Error(w, "Some Error Occured", http.StatusInternalServerError)
			return
		}

		utils.WriteJSON(w, http.StatusOK, leftRoom)
	}
}

func HandleRoomStats(hub *types.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomId := r.URL.Query().Get("roomId")
		if roomId == "" {
			utils.WriteJSON(w, http.StatusOK, srv.HubStats(hub))
			return
		}
		utils.WriteJSON(w, http.StatusOK, srv.RoomStats(hub, roomId))
	}
}
