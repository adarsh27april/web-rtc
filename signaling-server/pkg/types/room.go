package types

import "fmt"

type Room struct {
	RoomId   *string `json:"roomId,omitempty"`
	ClientId *string `json:"clientId,omitempty"`
	Status   *string `json:"status,omitempty" validate:"oneof=joined left created"`
}

// it ensure room and client are non empty
func (r *Room) ValidateLeaveRoom() error {
	if r.RoomId == nil || *r.RoomId == "" {
		return fmt.Errorf("room_id is required")
	}
	if r.ClientId == nil || *r.ClientId == "" {
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
