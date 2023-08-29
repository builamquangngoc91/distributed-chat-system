package handlers

import (
	"net/http"
	"strings"
	"time"
	"tylerbui430/user-service/domains"
	"tylerbui430/user-service/models"
	"tylerbui430/user-service/repositories"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserHandlers interface {
	RouteGroup(r *gin.Engine)

	CreateUser(c *gin.Context)
	GetUser(c *gin.Context)

	GetToken(c *gin.Context)
}

type userHandlers struct {
	db *gorm.DB
}

func NewUserHandlers(db *gorm.DB) UserHandlers {
	return &userHandlers{
		db: db,
	}
}

func (u *userHandlers) RouteGroup(rg *gin.Engine) {
	rg.POST("/users/create", u.CreateUser)
	rg.GET("/users/profile", u.GetUserProfile)

	rg.POST("/users/getToken", u.GetToken)
}

func (u *userHandlers) CreateUser(c *gin.Context) {
	ctx := c.Request.Context()

	var req domains.CreateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	_, err := repositories.GetUser(ctx, u.db, &repositories.GetUserArgs{
		Username: req.Username,
	})
	switch err {
	case nil:
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "user already exists",
		})
		return
	case gorm.ErrRecordNotFound:
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		if err := repositories.CreateUser(ctx, u.db, &models.User{
			UserID:       uuid.NewString(),
			Username:     req.Username,
			Name:         req.Name,
			PasswordHash: string(passwordHash),
		}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (u *userHandlers) GetUserProfile(c *gin.Context) {
	ctx := c.Request.Context()

	token := strings.TrimLeft("bearer", c.GetHeader("authorization"))
	
}

func (u *userHandlers) GetToken(c *gin.Context) {
	ctx := c.Request.Context()

	var req domains.GetTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	user, err := repositories.GetUser(ctx, u.db, &repositories.GetUserArgs{
		Username: req.Username,
	})
	switch err {
	case nil:
		if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "username/password is not correct",
			})
			return
		}

		expirationTime := time.Now().Add(30 * 24 * 60 * time.Minute)
		claims := &domains.Claims{
			Username: user.Username,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
			},
		}

		jwtKey := []byte("fdsfdsafdsafdsfdadcxa")
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"type":         "bearer",
			"access_token": tokenString,
			"expires_at":   expirationTime.Format(time.DateTime),
		})
	case gorm.ErrRecordNotFound:
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "username not found",
		})
		return
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
