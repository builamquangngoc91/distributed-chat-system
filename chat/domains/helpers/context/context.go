package context

import "context"

type contextKey string

const (
	UserID contextKey = "user_id"
)

func GetUserIDFromCtx(ctx context.Context) string {
	return ctx.Value(UserID).(string)
}
