package telemetry

import "context"

const userIDKey contextKey = "user_id"

// WithUserID stores a user ID in the context.
func WithUserID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, userIDKey, id)
}

// UserIDFrom extracts the user ID from context.
func UserIDFrom(ctx context.Context) string {
	id, _ := ctx.Value(userIDKey).(string)
	return id
}
