package repositories

import (
	"chat-service/models"
	"context"

	"gorm.io/gorm"
)

var _ GroupUserRepositoryI = &groupUserRepository{}

type (
	groupUserRepository struct{}

	GroupUserRepositoryI interface {
		Create(context.Context, *gorm.DB, *models.GroupUser) error
		GetGroupUser(ctx context.Context, db *gorm.DB, args *GetGroupUserArgs) (groupUser *models.GroupUser, err error)
		ListGroupUsers(ctx context.Context, db *gorm.DB, args *ListGroupUsersArgs) (groupUsers models.GroupUsers, err error)
	}

	GetGroupUserArgs struct {
		GroupID string
		UserID  string
	}

	ListGroupUsersArgs struct {
		GroupID string
	}
)

func NewGroupUserRepository() GroupUserRepositoryI {
	return &groupUserRepository{}
}

func (r *groupUserRepository) Create(ctx context.Context, db *gorm.DB, group *models.GroupUser) error {
	return db.
		WithContext(ctx).
		Create(group).
		Error
}

func (r *groupUserRepository) GetGroupUser(ctx context.Context, db *gorm.DB, args *GetGroupUserArgs) (*models.GroupUser, error) {
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

func (r *groupUserRepository) ListGroupUsers(ctx context.Context, db *gorm.DB, args *ListGroupUsersArgs) (groupUsers models.GroupUsers, err error) {
	db = db.WithContext(ctx)

	if args.GroupID != "" {
		db = db.Where("group_id = ?", args.GroupID)
	}

	if err := db.Find(&groupUsers).Error; err != nil {
		return nil, err
	}

	return groupUsers, nil
}
