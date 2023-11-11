package handlers

import (
	"fmt"
	"net/http"

	"chat-service/domains"
	contextHelpers "chat-service/domains/helpers/context"
	"chat-service/enums"
	"chat-service/middlewares"
	"chat-service/models"
	"chat-service/repositories"
	"chat-service/services/authservice"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var _ ChatHandlers = &chatHandlers{}

type ChatHandlers interface {
	RouteGroup(r *gin.Engine)

	CreateChatGroup(c *gin.Context)
}

type ChatHandlersDeps struct {
	DB *gorm.DB
}

type chatHandlers struct {
	db                  *gorm.DB
	authService         authservice.AuthService
	groupRepository     repositories.GroupRepositoryI
	groupUserRepository repositories.GroupUserRepositoryI
}

func NewChatHandlers(deps *ChatHandlersDeps) ChatHandlers {
	if deps == nil {
		return nil
	}

	return &chatHandlers{
		db:                  deps.DB,
		authService:         authservice.NewAuthService(),
		groupRepository:     repositories.NewGroupRepository(),
		groupUserRepository: repositories.NewGroupUserRepository(),
	}
}

func (u *chatHandlers) RouteGroup(rg *gin.Engine) {
	rg.POST("/chats/createChatGroup", middlewares.TokenAuthMiddleware(u.authService), u.CreateChatGroup)
}

func (u *chatHandlers) CreateChatGroup(c *gin.Context) {
	ctx := c.Request.Context()
	userID := contextHelpers.GetUserIDFromCtx(ctx)

	fmt.Println(userID)

	var req domains.CreateChatGroupReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError,
			domains.ErrorResp{
				Message: err.Error(),
			},
		)
		c.Error(err)
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest,
			domains.ErrorResp{
				Message: err.Error(),
			},
		)
		return
	}

	req.UserIDs = append(req.UserIDs, userID)
	listUserResp, err := u.authService.ListUsers(ctx, &authservice.ListUsersReq{
		UserIDs: req.UserIDs,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			domains.ErrorResp{
				Message: err.Error(),
			},
		)
		return
	}

	usersMap := make(map[string]*authservice.User)
	for _, user := range listUserResp.Users {
		usersMap[user.UserID] = user
	}

	// TODO: check user_ids and users

	err = u.db.Transaction(func(tx *gorm.DB) error {
		group := &models.Group{
			GroupID: uuid.NewString(),
			Name:    req.Name,
			Type:    req.Type,
		}
		if err := u.groupRepository.Create(ctx, tx, group); err != nil {
			return err
		}

		for _, userID := range req.UserIDs {
			groupUser := &models.GroupUser{
				GroupUserID: uuid.NewString(),
				GroupID:     group.GroupID,
				UserID:      userID,
				Name:        usersMap[userID].Name,
				Status:      enums.Joined.String(),
			}
			if err := u.groupUserRepository.Create(ctx, tx, groupUser); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			domains.ErrorResp{
				Message: err.Error(),
			})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
