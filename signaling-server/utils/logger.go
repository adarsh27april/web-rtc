package utils

import (
	"fmt"
	"log"
)

func LogRoom(roomID, ClientId, message string, args ...any) {
	logPrefix := fmt.Sprintf("[Room:%s] [Client:%s] ", roomID, ClientId)
	log.Printf(logPrefix+message, args...)
}
