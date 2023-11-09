package handlers

import (
	"chat-service/middlewares"
	"fmt"
	"net/http"

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
	db *gorm.DB
}

func NewChatHandlers(deps *ChatHandlersDeps) ChatHandlers {
	if deps == nil {
		return nil
	}

	return &chatHandlers{
		db: deps.DB,
	}
}

func (u *chatHandlers) RouteGroup(rg *gin.Engine) {
	rg.POST("/chats/createChatGroup", middlewares.TokenAuthMiddleware(), u.CreateChatGroup)
}

func (u *chatHandlers) CreateChatGroup(c *gin.Context) {
	fmt.Printf("user_id : %s", c.Request.Context().Value("user_id"))
	c.JSON(http.StatusOK, gin.H{})
}
