package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type (
	UserID  string
	GroupID string

	UserConnectionIdx int64
	ConnectionInfos   struct {
		UserID            UserID
		ExpiresAt         time.Time
		UserConnectionIdx UserConnectionIdx
	}

	MessageType int64

	MessageData struct {
		FromUserID UserID  `json:"from_user_id,omitempty"`
		ToUserID   UserID  `json:"to_user_id,omitempty"`
		ToGroupID  GroupID `json:"to_group_id,omitempty"`
		Text       string  `json:"text,omitempty"`
	}

	Message struct {
		Type MessageType  `json:"type,omitempty"`
		Data *MessageData `json:"data,omitempty"`
	}
)

const (
	PingMessage MessageType = iota + 1
	TextMessage
)

var (
	userConnectionsMap = make(map[UserID]map[UserConnectionIdx]*websocket.Conn)
	connectionInfosMap = make(map[*websocket.Conn]*ConnectionInfos)
	messageCh          = make(chan Message)

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func reader(conn *websocket.Conn) {
	for {
		var msg Message
		if err := conn.ReadJSON(&msg); err != nil {
			log.Println(err)
			return
		}

		switch msg.Type {
		case PingMessage:
			now := time.Now()
			connectionInfos := connectionInfosMap[conn]
			connectionInfos.ExpiresAt = now.Add(5 * time.Second)
			fmt.Printf("user (%s) PING\n", connectionInfos.UserID)
		case TextMessage:
			fmt.Println(msg.Data.Text)
			messageCh <- msg
		}
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("user_id")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	if userConnectionsMap[UserID(userID)] == nil {
		userConnectionsMap[UserID(userID)] = make(map[UserConnectionIdx]*websocket.Conn)
	}
	userConnectionIdx := len(userConnectionsMap[UserID(userID)]) + 1
	connectionInfosMap[conn] = &ConnectionInfos{
		UserID:            UserID(userID),
		ExpiresAt:         time.Now().Add(10 * time.Second),
		UserConnectionIdx: UserConnectionIdx(userConnectionIdx),
	}
	userConnectionsMap[UserID(userID)][UserConnectionIdx(userConnectionIdx)] = conn

	go reader(conn)
	fmt.Printf("user (%s) connected\n", userID)
}

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)
}

func startConnectionsManagement() {
	go func() {
		for {
			for connection, connectionInfos := range connectionInfosMap {
				if connectionInfos.ExpiresAt.Before(time.Now()) {
					connection.Close()
					delete(connectionInfosMap, connection)
					delete(userConnectionsMap[connectionInfos.UserID], connectionInfos.UserConnectionIdx)
					fmt.Printf("user (%s) disconnected\n", connectionInfos.UserID)
				}
			}
			time.Sleep(1 * time.Second)
		}
	}()
}

func sendMessageToUsers() {
	go func() {
		for msg := range messageCh {
			toUserID := msg.Data.ToUserID

			connections := userConnectionsMap[toUserID]

			for _, connection := range connections {
				connection.WriteJSON(msg)
			}
		}
	}()
}

func main() {
	fmt.Println("Hello World")
	setupRoutes()
	//startConnectionsManagement()
	sendMessageToUsers()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
