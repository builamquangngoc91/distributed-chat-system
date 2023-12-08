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
		GetGroupUser(ctx context.Context, db *gorm.DB, args *GetGroupUserArgs) (group *models.GroupUser, err error)
	}

	GetGroupUserArgs struct {
		GroupID string
		UserID  string
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

func (GroupUserRepository) GetGroupUser(ctx context.Context, db *gorm.DB, args *GetGroupUserArgs) (*models.GroupUser, error) {
	db = db.
		WithContext(ctx)

	if args.GroupID != "" {
		db = db.Where("group_id = ?", args.GroupID)
	}
	if args.UserID != "" {
		db = db.Where("user_id = ?", args.UserID)
	}

	var groupUser models.GroupUser
	if err := db.First(&groupUser).Error; err != nil {
		return nil, err
	}

	return &groupUser, nil
}
