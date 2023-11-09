package repositories

import (
	"auth-service/models"
	"context"

	"gorm.io/gorm"
)

func CreateUser(ctx context.Context, db *gorm.DB, user *models.User) error {
	return db.WithContext(ctx).Table("users").Create(user).Error
}

type GetUserArgs struct {
	Username string
}

func GetUser(ctx context.Context, db *gorm.DB, args *GetUserArgs) (*models.User, error) {
	query := db.WithContext(ctx).Table("users")
	if args.Username != "" {
		query.Where("username = ?", args.Username)
	}

	var user models.User
	result := query.First(&user)

	return &user, result.Error
}
