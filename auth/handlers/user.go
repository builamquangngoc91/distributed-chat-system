package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"auth-service/configs"
	"auth-service/domains"
	"auth-service/models"
	"auth-service/repositories"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	jwtClaimIDTmpl      = "jwt-claim-id-%s"
	authorizationHeader = "Authorization"
	bearerType          = "Bearer"
	expiration          = 30 * 24 * 60 * 60 * time.Second
)

var (
	_ UserHandlers = &userHandlers{}
)

type UserHandlers interface {
	RouteGroup(r *gin.Engine)

	CreateUser(c *gin.Context)
	GetUserProfile(c *gin.Context)
	GetToken(c *gin.Context)
	RevokeToken(c *gin.Context)

	ListUsers(c *gin.Context)
}

type UserHandlersDeps struct {
	DB          *gorm.DB
	RedisClient *redis.Client
	UserRepo    repositories.UserRepository
}

type userHandlers struct {
	db          *gorm.DB
	redisClient *redis.Client
	userRepo    repositories.UserRepository
}

func NewUserHandlers(deps *UserHandlersDeps) UserHandlers {
	if deps == nil {
		return nil
	}

	return &userHandlers{
		db:          deps.DB,
		redisClient: deps.RedisClient,
		userRepo:    deps.UserRepo,
	}
}

func (u *userHandlers) RouteGroup(rg *gin.Engine) {
	rg.POST("/users/create", u.CreateUser)
	rg.GET("/users/profile", u.GetUserProfile)
	rg.POST("/users/getToken", u.GetToken)
	rg.DELETE("/users/revokeToken", u.RevokeToken)

	rg.POST("/users/listUsers", u.ListUsers)
}

func (u *userHandlers) CreateUser(c *gin.Context) {
	ctx := c.Request.Context()

	var req domains.CreateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest,
			domains.ErrorResponse{
				Message: err.Error(),
			},
		)
		return
	}

	_, err := u.userRepo.Get(ctx, u.db, &repositories.GetUserArgs{
		Username: req.Username,
	})
	switch err {
	case nil:
		c.JSON(http.StatusBadRequest,
			domains.ErrorResponse{
				Message: "user already exists",
			},
		)
		return
	case gorm.ErrRecordNotFound:
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				domains.ErrorResponse{
					Message: err.Error(),
				},
			)
			return
		}

		if err := u.userRepo.Create(ctx, u.db, &models.User{
			UserID:       uuid.NewString(),
			Username:     req.Username,
			Name:         req.Name,
			PasswordHash: string(passwordHash),
		}); err != nil {
			c.JSON(http.StatusInternalServerError,
				domains.ErrorResponse{
					Message: err.Error(),
				},
			)
			return
		}
	default:
		c.JSON(http.StatusInternalServerError,
			domains.ErrorResponse{
				Message: err.Error(),
			},
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (u *userHandlers) GetUserProfile(c *gin.Context) {
	ctx := c.Request.Context()

	var claims domains.Claims
	tokenStr := strings.TrimPrefix(c.GetHeader(authorizationHeader), fmt.Sprintf("%s ", bearerType))
	_, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(configs.Cfg.AuthService.JwtKey), nil
	})
	switch err {
	case nil:
		// no-op
	default:
		c.JSON(http.StatusInternalServerError,
			domains.ErrorResponse{
				Message: err.Error(),
			},
		)
		return
	}

	jwtClaimID := fmt.Sprintf(jwtClaimIDTmpl, claims.ID)
	_, err = u.redisClient.Get(ctx, jwtClaimID).Result()
	switch err {
	case redis.Nil:
		// no-op
	case nil:
		c.JSON(http.StatusUnauthorized,
			domains.ErrorResponse{
				Message: "token was revoked",
			},
		)
		return
	default:
		c.JSON(http.StatusInternalServerError,
			domains.ErrorResponse{
				Message: err.Error(),
			},
		)
		return
	}

	user, err := u.userRepo.Get(ctx, u.db, &repositories.GetUserArgs{
		Username: claims.Username,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			domains.ErrorResponse{
				Message: err.Error(),
			},
		)
		return
	}

	c.JSON(http.StatusOK,
		domains.CreateUserResp{
			UserID:   user.UserID,
			Username: user.Username,
			Name:     user.Name,
		},
	)
}

func (u *userHandlers) GetToken(c *gin.Context) {
	ctx := c.Request.Context()

	var req domains.GetTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError,
			domains.ErrorResponse{
				Message: err.Error(),
			},
		)
		c.Error(err)
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest,
			domains.ErrorResponse{
				Message: err.Error(),
			},
		)
		return
	}

	user, err := u.userRepo.Get(ctx, u.db, &repositories.GetUserArgs{
		Username: req.Username,
	})
	switch err {
	case nil:
		if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
			c.JSON(http.StatusBadRequest,
				domains.ErrorResponse{
					Message: "password is incorrect",
				},
			)
			return
		}

		expirationTime := time.Now().Add(expiration)
		claims := &domains.Claims{
			Username: user.Username,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
				ID:        uuid.NewString(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(configs.Cfg.AuthService.JwtKey))
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				domains.ErrorResponse{
					Message: err.Error(),
				},
			)
			return
		}

		c.JSON(http.StatusOK,
			domains.GetTokenResp{
				Type:        strings.ToLower(bearerType),
				AccessToken: tokenString,
				ExpiresAt:   expirationTime.Format(time.DateTime),
			},
		)
		return
	case gorm.ErrRecordNotFound:
		c.JSON(http.StatusBadRequest,
			domains.ErrorResponse{
				Message: "username not found",
			},
		)
		return
	default:
		c.JSON(http.StatusInternalServerError,
			domains.ErrorResponse{
				Message: err.Error(),
			},
		)
		return
	}
}

func (u *userHandlers) RevokeToken(c *gin.Context) {
	ctx := c.Request.Context()

	var claims domains.Claims
	tokenStr := strings.TrimPrefix(c.GetHeader(authorizationHeader), fmt.Sprintf("%s ", bearerType))
	_, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(configs.Cfg.AuthService.JwtKey), nil
	})
	switch err {
	case nil:
		// no-op
	default:
		c.JSON(http.StatusInternalServerError,
			domains.ErrorResponse{
				Message: err.Error(),
			},
		)
		return
	}

	jwtClaimID := fmt.Sprintf(jwtClaimIDTmpl, claims.ID)
	if err := u.redisClient.Set(ctx, jwtClaimID, time.Now(), 0).Err(); err != nil {
		c.JSON(http.StatusInternalServerError,
			domains.ErrorResponse{
				Message: err.Error(),
			},
		)
		return
	}

	c.Status(http.StatusOK)
}

func (u *userHandlers) ListUsers(c *gin.Context) {
	ctx := c.Request.Context()

	var req domains.ListUsersReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError,
			domains.ErrorResponse{
				Message: err.Error(),
			},
		)
		c.Error(err)
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest,
			domains.ErrorResponse{
				Message: err.Error(),
			},
		)
		return
	}

	users, err := u.userRepo.List(ctx, u.db, &repositories.ListUsersArgs{
		IDs: req.UserIDs,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			domains.ErrorResponse{
				Message: err.Error(),
			},
		)
	}

	resp := domains.ListUsersResp{}
	for _, user := range users {
		resp.Users = append(resp.Users, &domains.User{
			UserID:   user.UserID,
			Username: user.Name,
			Name:     user.Name,
		})
	}

	c.JSON(http.StatusOK, &resp)
}
