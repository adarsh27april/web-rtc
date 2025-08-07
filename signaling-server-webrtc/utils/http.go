package utils

import (
	"encoding/json"
	"net/http"
	"signaling-server-webrtc/pkg/types"
)

func WriteError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func WriteJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

func DecodeRoomRequest(r *http.Request) (types.Room, error) {
	var room types.Room
	err := json.NewDecoder(r.Body).Decode(&room)
	return room, err
}
