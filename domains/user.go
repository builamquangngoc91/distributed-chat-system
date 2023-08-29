package domains

import (
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

type (
	CreateUserReq struct {
		Username string
		Name     string
		Password string
	}

	GetTokenReq struct {
		Username string
		Password string
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
