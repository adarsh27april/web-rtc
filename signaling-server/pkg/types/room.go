package types

import "fmt"

type Room struct {
	RoomId   *string `json:"roomId,omitempty"`
	ClientID *string `json:"clientId,omitempty"`
	Status   *string `json:"status,omitempty" validate:"oneof=joined left"`
}

// it ensure room and client are non empty
func (r *Room) ValidateLeaveRoom() error {
	if r.RoomId == nil || *r.RoomId == "" {
		return fmt.Errorf("room_id is required")
	}
	if r.ClientID == nil || *r.ClientID == "" {
		return fmt.Errorf("client_id is required")
	}
	return nil
}

type RoomStats struct {
	RoomID  string   `json:"roomId"`
	Clients []string `json:"clients"`
}

type HubStats struct {
	TotalRooms int         `json:"totalRooms"`
	Rooms      []RoomStats `json:"rooms"`
}
