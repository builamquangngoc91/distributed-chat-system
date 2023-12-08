package repositories

import (
	"auth-service/models"
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Test_NewUserRepository(t *testing.T) {
	repo := NewUserRepository()
	assert.NotNil(t, repo)
}

func Test_UserRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when creating gormDB", err)
	}

	user := &models.User{
		UserID:       "userID",
		Username:     "username",
		PasswordHash: "passwordHash",
		Name:         "name",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO \"users\" \\(\"user_id\",\"username\",\"password_hash\",\"name\",\"created_at\",\"updated_at\"\\) VALUES \\(\\$1,\\$2,\\$3,\\$4,\\$5,\\$6\\)").
			WithArgs(&user.UserID, &user.Username, &user.PasswordHash, &user.Name, &user.CreatedAt, &user.UpdatedAt).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		userRepo := NewUserRepository()

		err = userRepo.Create(context.Background(), gormDB, user)
		assert.NoError(t, err)
	})
}

func Test_UserRepository_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when creating gormDB", err)
	}

	user := &models.User{
		UserID:       "userID",
		Username:     "username",
		PasswordHash: "passwordHash",
		Name:         "name",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"user_id", "username", "password_hash", "name", "created_at", "updated_at"}).
			AddRow(user.UserID, user.Username, user.PasswordHash, user.Name, user.CreatedAt, user.UpdatedAt)
		mock.ExpectQuery("SELECT \\* FROM \"users\" WHERE username = \\$1 ORDER BY \"users\".\"user_id\" LIMIT 1").
			WithArgs(&user.Username).
			WillReturnRows(rows)

		userRepo := NewUserRepository()

		userDB, err := userRepo.Get(context.Background(), gormDB, &GetUserArgs{
			Username: user.Username,
		})
		assert.NoError(t, err)
		assert.NotNil(t, userDB)
	})
}

func Test_UserRepository_List(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when creating gormDB", err)
	}

	user := &models.User{
		UserID:       "userID",
		Username:     "username",
		PasswordHash: "passwordHash",
		Name:         "name",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"user_id", "username", "password_hash", "name", "created_at", "updated_at"}).
			AddRow(user.UserID, user.Username, user.PasswordHash, user.Name, user.CreatedAt, user.UpdatedAt)
		mock.ExpectQuery("SELECT \\* FROM \"users\" WHERE user_id IN \\(\\$1\\)").
			WithArgs(&user.UserID).
			WillReturnRows(rows)

		userRepo := NewUserRepository()

		userDB, err := userRepo.List(context.Background(), gormDB, &ListUsersArgs{
			IDs: []string{user.UserID},
		})
		assert.NoError(t, err)
		assert.NotNil(t, userDB)
	})
}
