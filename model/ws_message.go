package model

import (
	"time"
)

type WSMessage struct {
	Type      string    `json:"type"`
	DeviceId  string    `json:"device_id"`
	Data      string    `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}
