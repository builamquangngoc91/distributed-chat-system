package domains

import "time"

type (
	CreateMessageRequest struct {
		GroupID string `json:"group_id"`
		Content string `json:"content"`
	}

	CreateMessageResponse struct {
		ID        string    `json:"id"`
		GroupID   string    `json:"group_id"`
		Content   string    `json:"content"`
		SentBy    string    `json:"sent_by"`
		CreatedAt time.Time `json:"created_at"`
	}

	MessageForMSQ struct {
		ID        string    `json:"id"`
		GroupID   string    `json:"group_id"`
		Content   string    `json:"content"`
		SentBy    string    `json:"sent_by"`
		UserIDs   []string  `json:"user_ids"`
		CreatedAt time.Time `json:"created_at"`
	}
)
