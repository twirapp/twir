package mcp_oauth

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrClientNotFound            = errors.New("MCP OAuth client not found")
	ErrClientAlreadyExists       = errors.New("MCP OAuth client already exists")
	ErrAuthorizationCodeNotFound = errors.New("MCP OAuth authorization code not found")
	ErrRefreshTokenNotFound      = errors.New("MCP OAuth refresh token not found")
	ErrAccessTokenNotFound       = errors.New("MCP OAuth access token not found")
	ErrRefreshTokenReuse         = errors.New("MCP OAuth refresh token reuse detected")
)

type RefreshTokenReuseError struct {
	FamilyID uuid.UUID
}

func (err *RefreshTokenReuseError) Error() string {
	return ErrRefreshTokenReuse.Error()
}

func (err *RefreshTokenReuseError) Is(target error) bool {
	return target == ErrRefreshTokenReuse
}
