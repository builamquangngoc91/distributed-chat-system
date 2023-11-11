package authservice

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"chat-service/configs"
)

type AuthService interface {
	GetUserProfile(ctx context.Context, args *GetUserProfileReq) (*GetUserProfileResp, error)
}

type authService struct {
}

type (
	GetUserProfileReq struct {
		Token string `json:"token"`
	}

	GetUserProfileResp struct {
		UserID   string `json:"user_id,omitempty"`
		Username string `json:"username,omitempty"`
		Name     string `json:"name,omitempty"`
	}
)

const (
	userProfileURL = "/users/profile"
)

func NewAuthService() AuthService {
	return &authService{}
}

func (s *authService) GetUserProfile(ctx context.Context, args *GetUserProfileReq) (*GetUserProfileResp, error) {
	url := configs.Cfg.AuthService.BaseURL + userProfileURL
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", args.Token))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var userPofile GetUserProfileResp
	if err := json.NewDecoder(res.Body).Decode(&userPofile); err != nil {
		return nil, err
	}

	return &userPofile, err
}
