package repositories

import (
	"chat-service/models"
	"context"

	"gorm.io/gorm"
)

var _ GroupUserRepositoryI = GroupUserRepository{}

type (
	GroupUserRepository struct{}

	GroupUserRepositoryI interface {
		Create(context.Context, *gorm.DB, *models.GroupUser) error
	}
)

func NewGroupUserRepository() GroupUserRepositoryI {
	return &GroupUserRepository{}
}

func (GroupUserRepository) Create(ctx context.Context, db *gorm.DB, group *models.GroupUser) error {
	return db.
		WithContext(ctx).
		Create(group).
		Error
}
