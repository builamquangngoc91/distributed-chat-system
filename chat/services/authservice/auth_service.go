package authservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"chat-service/configs"
)

type AuthService interface {
	GetUserProfile(ctx context.Context, args *GetUserProfileReq) (*GetUserProfileResp, error)
	ListUsers(ctx context.Context, args *ListUsersReq) (*ListUsersResp, error)
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

	ListUsersReq struct {
		UserIDs []string `json:"user_ids"`
	}

	User struct {
		UserID   string `json:"user_id"`
		Username string `json:"username"`
		Name     string `json:"name"`
	}

	ListUsersResp struct {
		Users []*User `json:"users"`
	}
)

const (
	userProfileURL = "/users/profile"
	checkUsersURL  = "/users/listUsers"
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

	// TODO: add commons parse Error
	if res.StatusCode == http.StatusOK {
		var userPofile GetUserProfileResp
		if err := json.NewDecoder(res.Body).Decode(&userPofile); err != nil {
			return nil, err
		}

		return &userPofile, err
	}

	return nil, fmt.Errorf("status %s", res.Status)
}

func (s *authService) ListUsers(ctx context.Context, args *ListUsersReq) (*ListUsersResp, error) {
	url := configs.Cfg.AuthService.BaseURL + checkUsersURL

	body, err := json.Marshal(args)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var resp ListUsersResp
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, err
	}

	return &resp, err
}
