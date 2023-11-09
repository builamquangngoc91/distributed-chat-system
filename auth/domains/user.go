package domains

import (
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

type (
	CreateUserReq struct {
		Username string `json:"username"`
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	CreateUserResp struct {
		UserID   string `json:"user_id,omitempty"`
		Username string `json:"username,omitempty"`
		Name     string `json:"name,omitempty"`
	}

	GetTokenReq struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	GetTokenResp struct {
		Type        string `json:"type"`
		AccessToken string `json:"access_token"`
		ExpiresAt   string `json:"expires_at"`
	}

	Claims struct {
		Username string `json:"username"`
		jwt.RegisteredClaims
	}
)

func (c *CreateUserReq) Validate() error {
	if c.Username == "" {
		return errors.New("missing username")
	}
	if c.Name == "" {
		return errors.New("missing name")
	}
	if c.Password == "" {
		return errors.New("missing password")
	}

	return nil
}

func (c *GetTokenReq) Validate() error {
	if c.Username == "" {
		return errors.New("missing username")
	}
	if c.Password == "" {
		return errors.New("missing password")
	}

	return nil
}
