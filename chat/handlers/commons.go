package handlers

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var (
	UserConnectionsMap = make(map[string]map[string]*websocket.Conn) // map[user_id][conn_id]
	ServerID           = uuid.NewString()
)
