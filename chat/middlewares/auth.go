package middlewares

import (
	"context"
	"net/http"
	"strings"

	"chat-service/domains"
	contextHelpers "chat-service/domains/helpers/context"
	"chat-service/services/authservice"

	"github.com/gin-gonic/gin"
)

func TokenAuthMiddleware(authService authservice.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationValue := c.Request.Header.Get("Authorization")
		token := strings.TrimPrefix(authorizationValue, "Bearer ")
		userProfile, err := authService.GetUserProfile(c.Request.Context(), &authservice.GetUserProfileReq{
			Token: token,
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, &domains.ErrorResp{
				Message: err.Error(),
			})
			return
		}

		ctx := context.WithValue(c.Request.Context(), contextHelpers.UserID, userProfile.UserID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
