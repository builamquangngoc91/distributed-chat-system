package repositories

import (
	"context"

	"auth-service/models"

	"gorm.io/gorm"
)

var _ UserRepository = &userRepository{}

//go:generate mockery --name UserRepository
type UserRepository interface {
	Create(ctx context.Context, db *gorm.DB, user *models.User) error
	Get(ctx context.Context, db *gorm.DB, args *GetUserArgs) (*models.User, error)
	List(ctx context.Context, db *gorm.DB, args *ListUsersArgs) (users []*models.User, _ error)
}

type userRepository struct {
}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (u *userRepository) Create(ctx context.Context, db *gorm.DB, user *models.User) error {
	return db.WithContext(ctx).Table("users").Create(user).Error
}

type GetUserArgs struct {
	Username string
}

func (u *userRepository) Get(ctx context.Context, db *gorm.DB, args *GetUserArgs) (*models.User, error) {
	query := db.WithContext(ctx).Table("users")
	if args.Username != "" {
		query.Where("username = ?", args.Username)
	}

	var user models.User
	result := query.First(&user)

	return &user, result.Error
}

type ListUsersArgs struct {
	IDs []string
}

func (u *userRepository) List(ctx context.Context, db *gorm.DB, args *ListUsersArgs) (users []*models.User, _ error) {
	query := db.WithContext(ctx).Table("users")
	if len(args.IDs) != 0 {
		query.Where("user_id IN ?", args.IDs)
	}

	result := query.Find(&users)

	return users, result.Error
}
