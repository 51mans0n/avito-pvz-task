package auth

import (
	"errors"
	"strings"
)

func ExtractRoleFromToken(token string) (string, error) {
	if strings.HasPrefix(token, "SOME_TOKEN_") {
		role := strings.TrimPrefix(token, "SOME_TOKEN_")
		if role == "moderator" || role == "employee" || role == "client" {
			return role, nil
		}
		return "", errors.New("unknown role in dummy token")
	}
	return "", errors.New("invalid token format")
}
