package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type GetProfileResp struct {
	UserID   string `json:"user_id,omitempty"`
	Username string `json:"username,omitempty"`
	Name     string `json:"name,omitempty"`
}

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationValue := c.Request.Header.Get("Authorization")
		token := strings.TrimPrefix(authorizationValue, "Bearer ")

		url := fmt.Sprintf("%s/users/profile", os.Getenv("AUTH_URL"))
		req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodGet, url, nil)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
		req.Header.Add("authorization", fmt.Sprintf("Bearer %s", token))

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
		defer res.Body.Close()

		var profile GetProfileResp
		if err := json.NewDecoder(res.Body).Decode(&profile); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		ctx := context.WithValue(c.Request.Context(), "user_id", profile.UserID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
