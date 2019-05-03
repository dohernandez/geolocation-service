package request

import "context"

type contextKey string

const ctxKey = contextKey("id")

// ToContext sets request ID to context
func ToContext(ctx context.Context, reqID string) context.Context {
	return context.WithValue(ctx, ctxKey, reqID)
}

// FromContext extracts request ID from context
func FromContext(ctx context.Context) string {
	value := ctx.Value(ctxKey)
	if value == nil {
		return ""
	}

	return value.(string)
}
