package api

import "context"

// внутренний ключ ― чтобы не пересекаться с ключами
type ctxKey string

const roleKey ctxKey = "role"

// WithRole кладёт роль в контекст
func WithRole(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, roleKey, role)
}

// GetRole достаёт роль из контекста
func GetRole(ctx context.Context) string {
	if v, ok := ctx.Value(roleKey).(string); ok {
		return v
	}
	return ""
}
