package models

import (
	"time"
)

type GroupMessage struct {
	MessageID string
	GroupID   string
	UserID    string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (GroupMessage) TableName() string {
	return "group_messages"
}
