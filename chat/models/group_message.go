package models

import (
	"chat-service/enums"
	"time"
)

type Message struct {
	GroupID   string
	Name      string
	Type      enums.GroupType
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (Message) TableName() string {
	return "messages"
}
