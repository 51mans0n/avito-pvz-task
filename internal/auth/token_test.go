package auth_test

import (
	"testing"

	"github.com/51mans0n/avito-pvz-task/internal/auth"
	"github.com/stretchr/testify/require"
)

func TestExtractRole_OK(t *testing.T) {
	role, err := auth.ExtractRole("SOME_TOKEN_moderator")
	require.NoError(t, err)
	require.Equal(t, "moderator", role)
}

func TestExtractRole_UnknownRole(t *testing.T) {
	_, err := auth.ExtractRole("SOME_TOKEN_hacker")
	require.Error(t, err)
}

func TestExtractRole_BadFormat(t *testing.T) {
	_, err := auth.ExtractRole("Bearer xxx")
	require.Error(t, err)
}
