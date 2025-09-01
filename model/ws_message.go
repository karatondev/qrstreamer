package model

import (
	"time"
)

type WSMessage struct {
	MsgStatus  bool      `json:"msg_status"`
	Type       string    `json:"type"`
	WhatsappId string    `json:"whatsapp_id"`
	Data       string    `json:"data"`
	Timestamp  time.Time `json:"timestamp"`
}
