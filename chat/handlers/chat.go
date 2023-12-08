package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"chat-service/domains"
	contextHelpers "chat-service/domains/helpers/context"
	"chat-service/enums"
	"chat-service/middlewares"
	"chat-service/models"
	"chat-service/repositories"
	"chat-service/services/authservice"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	_                  ChatHandlers = &chatHandlers{}
	userConnectionsMap              = make(map[string][]*websocket.Conn)
)

type ChatHandlers interface {
	RouteGroup(r *gin.Engine)

	CreateChatGroupHandler(c *gin.Context)
	WSHandler(c *gin.Context)
}

type ChatHandlersDeps struct {
	DB            *gorm.DB
	RedisClient   *redis.Client
	KafkaProducer sarama.SyncProducer
}

type chatHandlers struct {
	db                     *gorm.DB
	redisClient            *redis.Client
	authService            authservice.AuthService
	groupRepository        repositories.GroupRepositoryI
	groupUserRepository    repositories.GroupUserRepositoryI
	groupMessageRepository repositories.GroupMessageRepositoryI
	kafkaProducer          sarama.SyncProducer
}

func NewChatHandlers(deps *ChatHandlersDeps) ChatHandlers {
	if deps == nil {
		return nil
	}

	return &chatHandlers{
		db:                     deps.DB,
		redisClient:            deps.RedisClient,
		authService:            authservice.NewAuthService(),
		groupRepository:        repositories.NewGroupRepository(),
		groupUserRepository:    repositories.NewGroupUserRepository(),
		groupMessageRepository: repositories.NewGroupMessageRepository(),
		kafkaProducer:          deps.KafkaProducer,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (u *chatHandlers) RouteGroup(rg *gin.Engine) {
	rg.POST("/chats/createChatGroup", middlewares.TokenAuthMiddleware(u.authService), u.CreateChatGroupHandler)
	rg.GET("/chats/ws", middlewares.TokenAuthMiddleware(u.authService), u.WSHandler)
}

func (u *chatHandlers) CreateChatGroupHandler(c *gin.Context) {
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

func (u *chatHandlers) WSHandler(c *gin.Context) {
	ctx := c.Request.Context()
	userID := contextHelpers.GetUserIDFromCtx(ctx)
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, domains.ErrorResp{
			Message: err.Error(),
		})
		return
	}

	userConnectionsMap[userID] = append(userConnectionsMap[userID], conn)
	u.listenChatMessages(ctx, conn)
}

func (u *chatHandlers) listenChatMessages(ctx context.Context, conn *websocket.Conn) {
	userID := contextHelpers.GetUserIDFromCtx(ctx)
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		var msg domains.Message
		json.Unmarshal(p, &msg)

		_, err = u.groupUserRepository.GetGroupUser(ctx, u.db, &repositories.GetGroupUserArgs{
			GroupID: msg.GroupID,
			UserID:  userID,
		})
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				fmt.Printf("can find groupUser(groupID: %s, userID: %s)", msg.GroupID, userID)
				continue
			}
			fmt.Printf("getGroupUser(groupID: %s, userID: %s) error: %s", msg.GroupID, userID, err.Error())
			continue
		}

		if err := u.groupMessageRepository.Create(ctx, u.db, &models.GroupMessage{
			MessageID: uuid.NewString(),
			GroupID:   msg.GroupID,
			UserID:    userID,
			Content:   msg.Content,
		}); err != nil {
			fmt.Printf("user %s save message content(%s) error : %s", userID, msg.Content, err.Error())
		}

		_, _, err = u.kafkaProducer.SendMessage(&sarama.ProducerMessage{
			Topic: "chat-message",
			Key:   sarama.StringEncoder(msg.GroupID),
			Value: sarama.StringEncoder(p),
		})
		if err != nil {
			fmt.Printf("kafka error: %s\n", err.Error())
		}

	}

}
