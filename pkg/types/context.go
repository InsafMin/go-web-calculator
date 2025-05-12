package types

import "context"

type contextKey string

const UserIDKey = contextKey("user_id")

func WithUserID(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

func GetUserID(ctx context.Context) int {
	val := ctx.Value(UserIDKey)
	if id, ok := val.(int); ok {
		return id
	}
	return -1
}
