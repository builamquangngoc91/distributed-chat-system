package models

import (
	"chat-service/enums"
	"time"
)

type Group struct {
	GroupID   string
	Name      string
	Type      enums.GroupType
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (Group) TableName() string {
	return "groups"
}

type Groups []*Group
