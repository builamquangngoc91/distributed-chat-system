package domains

import (
	"chat-service/enums"
	"fmt"
)

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

func (c *CreateChatGroupReq) Validate() error {
	if len(c.UserIDs) == 0 {
		return fmt.Errorf("missing user_ids")
	}
	if c.Type != enums.Group.String() && c.Type != enums.Personal.String() {
		return fmt.Errorf("invalid type")
	}
	if c.Type == enums.Personal.String() {
		if len(c.UserIDs) != 1 {
			return fmt.Errorf("user_ids must contains one value")
		}
	}

	return nil
}
