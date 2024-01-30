package models

import (
	"time"

	"gorm.io/gorm"
)

type ChatMessage struct {
	gorm.Model

	Sender string
	Recipient string
	Content string 
	CreatedAt time.Time
}

