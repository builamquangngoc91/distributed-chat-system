package models

import (
	"chat-service/enums"
	"time"
)

type GroupUser struct {
	GroupUserID string
	UserID      string
	GroupID     string
	Name        string
	Status      enums.GroupUserStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

func (GroupUser) TableName() string {
	return "groups_users"
}
