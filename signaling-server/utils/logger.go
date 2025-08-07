package utils

import (
	"fmt"
	"log"
)

func LogRoom(roomID, clientID, message string, args ...any) {
	logPrefix := fmt.Sprintf("[Room:%s] [Client:%s] ", roomID, clientID)
	log.Printf(logPrefix+message, args...)
}
