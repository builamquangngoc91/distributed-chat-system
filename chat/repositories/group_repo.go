package repositories

import (
	"chat-service/enums"
	"chat-service/models"
	"context"
	"errors"

	"gorm.io/gorm"
)

var _ GroupRepositoryI = GroupRepository{}

type (
	GroupRepository struct{}

	GetGroupArgs struct {
		ID string
	}
	GetGroupsArgs struct {
		IDs []string
	}

	GroupRepositoryI interface {
		GetGroup(context.Context, *gorm.DB, *GetGroupArgs) (*models.Group, error)
		GetGroups(context.Context, *gorm.DB, *GetGroupsArgs) (models.Groups, error)
		Create(context.Context, *gorm.DB, *models.Group) error
		Update(context.Context, *gorm.DB, *models.Group) error
		Delete(context.Context, *gorm.DB, *models.Group) error
	}
)

func NewGroupRepository() GroupRepositoryI {
	return &GroupRepository{}
}

func (GroupRepository) Create(ctx context.Context, db *gorm.DB, group *models.Group) error {
	return db.
		WithContext(ctx).
		Create(group).
		Error
}

// Get implements GroupRepositoryI.
func (GroupRepository) GetGroup(ctx context.Context, db *gorm.DB, args *GetGroupArgs) (group *models.Group, err error) {
	db = db.
		WithContext(ctx)

	if args.ID != "" {
		db.Where("id = ?", args.ID)
	}
	err = db.First(group).Error

	return
}

// GetGroups implements GroupRepositoryI.
func (GroupRepository) GetGroups(ctx context.Context, db *gorm.DB, args *GetGroupsArgs) (groups models.Groups, err error) {
	db = db.
		WithContext(ctx)

	if len(args.IDs) > 0 {
		db.Where("id IN ?", args.IDs)
	}
	err = db.Find(groups).Error

	return
}

// Update implements GroupRepositoryI.
func (GroupRepository) Update(ctx context.Context, db *gorm.DB, group *models.Group) (err error) {
	db = db.
		WithContext(ctx).
		Where("group_id = ?", group.GroupID).
		Updates(group)
	if err = db.Error; err != nil {
		return err
	}

	if db.RowsAffected == 0 {
		return errors.New(enums.NotRowsAffected)
	}

	return nil
}

// Delete implements GroupRepositoryI.
func (GroupRepository) Delete(ctx context.Context, db *gorm.DB, group *models.Group) (err error) {
	db = db.
		WithContext(ctx).
		Where("group_id = ?", group.GroupID).
		Update("deleted_at", "now()")
	if err = db.Error; err != nil {
		return err
	}

	if db.RowsAffected == 0 {
		return errors.New(enums.NotRowsAffected)
	}

	return nil
}
