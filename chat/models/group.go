package models

import (
	"time"
)

type Group struct {
	GroupID   string
	Name      string
	Type      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (Group) TableName() string {
	return "groups"
}

type Groups []*Group
