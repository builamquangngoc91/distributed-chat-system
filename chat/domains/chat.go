package domains

type (
	CreateChatGroupReq struct {
		Type    string   `json:"type"`
		UserIDs []string `json:"user_ids"`
		Name    string   `json:"name"`
	}

	CreateChatGroupResp struct {
		Success bool `json:"success"`
	}
)
