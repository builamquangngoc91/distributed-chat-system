package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"chat-service/domains"

	"github.com/redis/go-redis/v9"
)

var (
	_ ChatHandlers = &chatHandlers{}
)

type ChatListener interface {
	StartChatListeners()
}

type ChatListenerDeps struct {
	RedisClient *redis.Client
}

type chatListener struct {
	redisClient *redis.Client
}

func NewChatListener(deps *ChatListenerDeps) ChatListener {
	if deps == nil {
		return nil
	}

	return &chatListener{
		redisClient: deps.RedisClient,
	}
}

func (u *chatListener) StartChatListeners() {
	ctx := context.Background()
	rdb := u.redisClient.PSubscribe(ctx, ChatGroupChannelPrefix)

	go func() {
		for msg := range rdb.Channel() {
			fmt.Println(msg.Channel, msg.Payload)

			var message domains.MessageForMSQ
			if err := json.Unmarshal([]byte(msg.Payload), &message); err != nil {
				fmt.Printf("unmarshal message error: %s", err.Error())
				continue
			}

			msg := domains.CreateMessageResponse{
				ID:        message.ID,
				GroupID:   message.GroupID,
				SentBy:    message.SentBy,
				Content:   message.Content,
				CreatedAt: message.CreatedAt,
			}

			for _, userID := range message.UserIDs {
				for _, conn := range UserConnectionsMap[userID] {
					conn.WriteJSON(msg)
				}
			}
		}
	}()
}
