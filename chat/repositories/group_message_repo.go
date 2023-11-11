package repositories

import (
	"chat-service/models"
	"context"

	"gorm.io/gorm"
)

var _ GroupMessageRepositoryI = GroupMessageRepository{}

type (
	GroupMessageRepository struct{}

	GroupMessageRepositoryI interface {
		Create(context.Context, *gorm.DB, *models.GroupMessage) error
	}
)

func NewGroupMessageRepository() GroupMessageRepositoryI {
	return &GroupMessageRepository{}
}

func (GroupMessageRepository) Create(ctx context.Context, db *gorm.DB, group *models.GroupMessage) error {
	return db.
		WithContext(ctx).
		Create(group).
		Error
}
