package api

import "context"

type ctxKeyRole struct{}

func WithRole(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, ctxKeyRole{}, role)
}

func GetRole(ctx context.Context) string {
	v := ctx.Value(ctxKeyRole{})
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
