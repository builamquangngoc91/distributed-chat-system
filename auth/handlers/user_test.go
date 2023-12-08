package handlers

import (
	"auth-service/configs"
	"auth-service/domains"
	"auth-service/models"
	"auth-service/repositories/mocks"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redismock/v9"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_CreateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	oldJwt := configs.Cfg.AuthService.JwtKey
	defer func() {
		configs.Cfg.AuthService.JwtKey = oldJwt
	}()

	jwtKey := "jwt-key"
	configs.Cfg.AuthService.JwtKey = jwtKey

	claimID := "claim-id"
	claims := &domains.Claims{
		Username: "username",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			ID:        claimID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(configs.Cfg.AuthService.JwtKey))
	assert.NoError(t, err)

	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()

		router := gin.Default()
		userRepoMock := new(mocks.UserRepository)
		redisDBMock, redisMock := redismock.NewClientMock()

		userHandlers := NewUserHandlers(&UserHandlersDeps{
			UserRepo:    userRepoMock,
			RedisClient: redisDBMock,
		})
		userHandlers.RouteGroup(router)

		redisMock.
			ExpectGet(fmt.Sprintf(jwtClaimIDTmpl, claimID)).
			RedisNil()
		userRepoMock.
			On("Get", mock.Anything, mock.Anything, mock.Anything).
			Return(&models.User{
				UserID:   "user-id",
				Username: "username",
				Name:     "name",
			}, nil).
			Once()

		request, err := http.NewRequest(http.MethodGet, "/users/profile", nil)
		request.Header.Set(authorizationHeader, tokenString)
		assert.NoError(t, err)

		router.ServeHTTP(w, request)
		fmt.Println(w.Code)
	})
}
