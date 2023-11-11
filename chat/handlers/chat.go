package handlers

import (
	"fmt"
	"net/http"

	contextHelpers "chat-service/domains/helpers/context"
	"chat-service/middlewares"
	"chat-service/services/authservice"

	"github.com/gin-gonic/gin"
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
	db          *gorm.DB
	authService authservice.AuthService
}

func NewChatHandlers(deps *ChatHandlersDeps) ChatHandlers {
	if deps == nil {
		return nil
	}

	return &chatHandlers{
		db:          deps.DB,
		authService: authservice.NewAuthService(),
	}
}

func (u *chatHandlers) RouteGroup(rg *gin.Engine) {
	rg.POST("/chats/createChatGroup", middlewares.TokenAuthMiddleware(u.authService), u.CreateChatGroup)
}

func (u *chatHandlers) CreateChatGroup(c *gin.Context) {
	ctx := c.Request.Context()
	userID := contextHelpers.GetUserIDFromCtx(ctx)

	fmt.Println(userID)

	c.JSON(http.StatusOK, gin.H{})
}
