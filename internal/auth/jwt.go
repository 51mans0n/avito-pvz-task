package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"strings"
)

var secret = []byte("supersecret") // TODO: брать из config/env

func ExtractRole(token string) (string, error) {
	// dummy‑token
	if strings.HasPrefix(token, "SOME_TOKEN_") {
		role := strings.TrimPrefix(token, "SOME_TOKEN_")
		switch role {
		case "moderator", "employee", "client":
			return role, nil
		default:
			return "", errors.New("unknown role in dummy token")
		}
	}

	// real JWT
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})
	if err != nil || !parsed.Valid {
		return "", errors.New("invalid token")
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("claims cast error")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return "", errors.New("role claim missing")
	}
	return role, nil
}

// helper для dummyLogin
func IssueDummyToken(role string) string { // "moderator"/"employee"/"client"
	return "SOME_TOKEN_" + role
}
